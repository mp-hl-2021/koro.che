package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"koro.che/usecases"
	"net/http"
	"time"
)

type Api struct {
	AccountUseCases usecases.AccountUseCasesInterface
	LinkUseCases    usecases.LinkUseCasesInterface
}

func NewApi(a usecases.AccountUseCasesInterface, l usecases.LinkUseCasesInterface) *Api {
	return &Api{
		AccountUseCases: a,
		LinkUseCases:    l,
	}
}

func (a *Api) Router() http.Handler {
	router := mux.NewRouter()

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
		if cookie.Expires.Unix() < time.Now().Unix() {
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
	if link, err := a.LinkUseCases.MakeRedirect(vars["key"]); err == nil {
		http.Redirect(writer, request, "https://" + link, http.StatusMovedPermanently)
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
	if link, err := a.LinkUseCases.GetRealLink(vars["key"]); err == nil {
		o := linkModel{Link: link}
		if err := json.NewEncoder(writer).Encode(o); err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)
	} else {
		writer.WriteHeader(http.StatusNotFound)
	}
}

func (a *Api) shortenLink(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	var m linkModel
	if err := json.NewDecoder(request.Body).Decode(&m); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	var shortLink, _ = a.LinkUseCases.ShortenLink(m.Link)
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
	writer.WriteHeader(http.StatusOK)
}

func (a *Api) getUserLinks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	o := struct {
		Links []string `json:"links"`
	}{
		Links: []string{"dota2", "lol_kek"},
	}
	if err := json.NewEncoder(w).Encode(o); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (a *Api) getUserLinkStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var dotaStat = usecases.LinkStat{LinkName: "dota2", UseCounter: 228}
	var lolKekStat = usecases.LinkStat{LinkName: "lol_kek", UseCounter: 1337}
	o := struct {
		Stats []usecases.LinkStat `json:"linkStats"`
	}{
		Stats: []usecases.LinkStat{dotaStat, lolKekStat},
	}
	if err := json.NewEncoder(w).Encode(o); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
