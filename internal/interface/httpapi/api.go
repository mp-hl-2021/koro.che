package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	link2 "koro.che/internal/domain/link"
	"koro.che/internal/interface/prom"
	"koro.che/internal/usecases/account"
	"koro.che/internal/usecases/link"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

type Api struct {
	AccountUseCases account.AccountUseCasesInterface
	LinkUseCases    link.LinkUseCasesInterface
	Logger          zerolog.Logger
	Ctx             context.Context
}

func NewApi(ctx context.Context, a account.AccountUseCasesInterface, l link.LinkUseCasesInterface) *Api {
	return &Api{
		AccountUseCases: a,
		LinkUseCases:    l,
		Logger:          log.With().Str("module", "http-server").Logger(),
		Ctx:             ctx,
	}
}

func (a *Api) Router() http.Handler {
	router := mux.NewRouter()
	router.Use(prom.Measurer())
	router.Use(a.logger)
	router.Handle("/metrics", promhttp.Handler())
	router.HandleFunc("/api/register", a.register).Methods(http.MethodPost)
	router.HandleFunc("/api/login", a.login).Methods(http.MethodPut)
	router.HandleFunc("/api/logout", a.authorize(a.logout)).Methods(http.MethodPut)
	router.HandleFunc("/api/shorten", a.shortenLink).Methods(http.MethodPost)
	router.HandleFunc("/api/{key}/real", a.getRealLink).Methods(http.MethodGet)
	router.HandleFunc("/{key}", a.redirectToRealLink).Methods(http.MethodGet)
	router.HandleFunc("/api/manage/{key}", a.authorize(a.deleteLink)).Methods(http.MethodDelete)
	router.HandleFunc("/api/manage/links", a.authorize(a.getUserLinks)).Methods(http.MethodGet)
	router.HandleFunc("/api/manage/stats", a.authorize(a.getUserLinkStats)).Methods(http.MethodGet)

	return router
}

func (a *Api) authorize(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		cookie, err := request.Cookie("token")
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		if cookie.Expires.Unix() < time.Now().Unix() && cookie.Expires.Unix() >= 0 {
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		token := cookie.Value
		id, err := a.AccountUseCases.Authenticate(token)
		if err != nil {
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(request.Context(), "account_id", id)
		handlerFunc(writer, request.WithContext(ctx))
	}
}

type registrationModel struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (a *Api) register(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	var m registrationModel
	if err := json.NewDecoder(request.Body).Decode(&m); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	acc, err := a.AccountUseCases.CreateAccount(m.Login, m.Password)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(err.Error()))
		return
	}
	a.LinkUseCases.CreateUserLinksStorage(acc.Id)

	location := fmt.Sprintf("/accounts/%s", acc.Id)
	writer.Header().Set("Location", location)
	writer.WriteHeader(http.StatusCreated)
}

type loginModel struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (a *Api) login(writer http.ResponseWriter, request *http.Request) {
	var m loginModel
	if err := json.NewDecoder(request.Body).Decode(&m); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := a.AccountUseCases.LoginToAccount(m.Login, m.Password)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
	}
	http.SetCookie(writer, &http.Cookie{Name: "token", Value: token})
	writer.Header().Set("Content-Type", "application/jwt")
	if _, err := writer.Write([]byte(token)); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
}

func (a *Api) logout(writer http.ResponseWriter, request *http.Request) {
	c := http.Cookie{
		Name:   "token",
		MaxAge: -1}
	http.SetCookie(writer, &c)

	if _, err := writer.Write([]byte("Old cookie deleted. Logged out!\n")); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.Redirect(writer, request, "/", http.StatusFound)
}

func (a *Api) redirectToRealLink(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	if lnk, err := a.LinkUseCases.MakeRedirect(vars["key"]); err == nil {
		http.Redirect(writer, request, "https://"+lnk, http.StatusMovedPermanently)
	} else {
		writer.WriteHeader(http.StatusNotFound)
	}
}

type linkModel struct {
	Link string `json:"link"`
}

func (a *Api) getRealLink(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	vars := mux.Vars(request)
	if lnk, err := a.LinkUseCases.GetRealLink(vars["key"]); err == nil {
		o := linkModel{Link: lnk}
		if err := json.NewEncoder(writer).Encode(o); err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)
	} else {
		writer.WriteHeader(http.StatusNotFound)
	}
}

func GetUserId(a *Api, r *http.Request) string {
	cookie, err := r.Cookie("token")
	if err != nil {
		return ""
	}
	if cookie.Expires.Unix() < time.Now().Unix() && cookie.Expires.Unix() >= 0 {
		return ""
	}
	token := cookie.Value
	id, err := a.AccountUseCases.Authenticate(token)
	if err != nil {
		return ""
	}
	return id
}

func (a *Api) shortenLink(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	var m linkModel
	if err := json.NewDecoder(request.Body).Decode(&m); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// get user id if exists
	userId := GetUserId(a, request)

	var shortLink, _ = a.LinkUseCases.ShortenLink(m.Link, userId)
	o := linkModel{Link: shortLink}
	if err := json.NewEncoder(writer).Encode(o); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusCreated)
}

func (a *Api) deleteLink(writer http.ResponseWriter, request *http.Request) {
	var m linkModel
	if err := json.NewDecoder(request.Body).Decode(&m); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// get user id if exists
	userId := GetUserId(a, request)

	if _, err := a.LinkUseCases.DeleteLink(m.Link, userId); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	writer.WriteHeader(http.StatusOK)
}

func urlChecker(ctx context.Context, in <-chan string, out chan<- bool, stop <-chan struct{}) error {
	for {
		select {
		case <-stop:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		case u, ok := <-in:
			if !ok {
				return nil
			}
			flag, err := checkUrl(ctx, u)
			if err != nil {
				return err
			}
			out <- flag
		}
	}
}

func checkUrl(ctx context.Context, urlStr string) (bool, error) {
	c := http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return false, err
	}
	resp, err := c.Do(req)
	if err != nil {
		if e, ok := err.(*url.Error); ok {
			if e.Timeout() || e.Temporary() {
				return false, nil
			} else {
				return false, err
			}
		} else {
			return false, err
		}
	}
	if resp != nil {
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusNotFound {
			return false, nil
		} else {
			return true, nil
		}
	}
	return false, nil
}

const (
	urlsConcurrency = 4
)

func getLinksAvailability(ctx context.Context, urls []string) ([]bool, error) {
	stopChan := make(chan struct{})
	errChan := make(chan string, urlsConcurrency)
	urlsChan := make(chan string)
	availabilitiesChan := make(chan bool, urlsConcurrency)
	go func() {
		for _, u := range urls {
			urlsChan <- u
		}
		close(urlsChan)
	}()
	var wg sync.WaitGroup
	wg.Add(urlsConcurrency)
	for i := 0; i < urlsConcurrency; i++ {
		go func(ctx context.Context, in <-chan string, out chan<- bool, stop <-chan struct{}) {
			err := urlChecker(ctx, in, out, stop)
			if err != nil {
				fmt.Fprintln(os.Stderr, "error url checking")
				errChan <- err.Error()
			}
			wg.Done()
		}(ctx, urlsChan, availabilitiesChan, stopChan)
	}

	availabilities := make([]bool, 0)

	for {
		if len(availabilities) == len(urls) {
			close(stopChan)
			break
		}
		select {
		case e := <-errChan:
			close(stopChan)
			return []bool{}, errors.New(e)
		case <-ctx.Done():
			break
		case b := <-availabilitiesChan:
			availabilities = append(availabilities, b)
		}
	}
	wg.Wait()
	return availabilities, nil
}

type UserLinkResponse struct {
	Link      string `json:"link"`
	Available bool   `json:"available"`
}

func (a *Api) getUserLinks(w http.ResponseWriter, r *http.Request) {
	var links []link2.UserLink
	userId := r.Context().Value("account_id").(string)
	links, err := a.LinkUseCases.GetUserLinks(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	realLinks := make([]string, len(links))
	for i, lnk := range links {
		realLinks[i] = lnk.RealLink
	}
	availabilities, err := getLinksAvailability(a.Ctx, realLinks)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := make([]UserLinkResponse, len(links))
	for i := 0; i < len(links); i++ {
		resp[i] = UserLinkResponse{Link: links[i].ShortenLink, Available: availabilities[i]}
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (a *Api) getUserLinkStats(w http.ResponseWriter, r *http.Request) {
	var m linkModel
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	o, err := a.LinkUseCases.GetLinkStats(m.Link)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := json.NewEncoder(w).Encode(o); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type responseWriterObserver struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (o *responseWriterObserver) WriteHeader(code int) {
	o.ResponseWriter.WriteHeader(code)
	if o.wroteHeader {
		return
	}
	o.wroteHeader = true
	o.status = code
}

func (o *responseWriterObserver) StatusCode() int {
	if !o.wroteHeader {
		return http.StatusOK
	}
	return o.status
}

func (a *Api) logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		o := &responseWriterObserver{ResponseWriter: w}
		next.ServeHTTP(o, r)
		a.Logger.Info().
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Str("protocol", r.Proto).
			Int("status-code", o.StatusCode()).
			Str("remote-addr", r.RemoteAddr).
			Dur("duration", time.Since(start)).
			Msg("")
	})
}
