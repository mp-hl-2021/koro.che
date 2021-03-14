package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	//7
	router.HandleFunc("/api/manage/{link}", deleteLink).Methods(http.MethodDelete)
	//8
	router.HandleFunc("/api/manage/links", getUserLinks).Methods(http.MethodGet)
	//9
	router.HandleFunc("/api/manage/stats", getUserLinkStats).Methods(http.MethodGet)

	router.HandleFunc("/api/shorten", shortenLink).Methods(http.MethodPost)
	router.HandleFunc("/api/{shorten_link}/real", getRealLink).Methods(http.MethodGet)
	router.HandleFunc("/{shorten_link}", redirectToRealLink).Methods(http.MethodGet)
	router.HandleFunc("/api/register", register).Methods(http.MethodPost)
	router.HandleFunc("/api/login", login).Methods(http.MethodPut)
	router.HandleFunc("/api/logout", logout).Methods(http.MethodPut)
	server := http.Server{
		Addr:         "localhost:8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      router,
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func redirectToRealLink(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	if vars["shorten_link"] == "kek" {
		http.Redirect(writer, request, "https://tsarn.website/sp", http.StatusMovedPermanently)
	} else {
		writer.WriteHeader(http.StatusNotFound)
	}
}

type linkModel struct {
	Link string `json:"link"`
}

func getRealLink(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	vars := mux.Vars(request)
	if vars["shorten_link"] == "kek" {
		o := linkModel{Link: "https://tsarn.website/sp"}
		if err := json.NewEncoder(writer).Encode(o); err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)
	} else {
		writer.WriteHeader(http.StatusNotFound)
	}
}

func shortenLink(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	var m linkModel
	if err := json.NewDecoder(request.Body).Decode(&m); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	o := linkModel{Link: "https://koro.che/kek"}
	if err := json.NewEncoder(writer).Encode(o); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusCreated)
}

type registrationModel struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func register(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	var m registrationModel
	if err := json.NewDecoder(request.Body).Decode(&m); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	writer.WriteHeader(http.StatusCreated)
}

type loginModel struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func login(writer http.ResponseWriter, request *http.Request) {
	var m loginModel
	if err := json.NewDecoder(request.Body).Decode(&m); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	token := "some_token"

	writer.Header().Set("Content-Type", "application/jwt")
	if _, err := writer.Write([]byte(token)); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
}

func logout(writer http.ResponseWriter, request *http.Request) {
	_, err := request.Cookie("_cookie")
	if err != http.ErrNoCookie {
		print(err, "Failed to get cookie")
	}
	http.Redirect(writer, request, "/", http.StatusFound)
}

func deleteLink(w http.ResponseWriter, r *http.Request) {
	//some delete logic
}

func getUserLinks(w http.ResponseWriter, r *http.Request) {
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

func getUserLinkStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	type linkStat struct {
		LinkName   string `json:"linkName"`
		UseCounter int64  `json:"useCounter"`
	}
	var dotaStat = linkStat{LinkName: "dota2", UseCounter: 228}
	var lolKekStat = linkStat{LinkName: "lol_kek", UseCounter: 1337}
	o := struct {
		Stats []linkStat `json:"linkStats"`
	}{
		Stats: []linkStat{dotaStat, lolKekStat},
	}
	if err := json.NewEncoder(w).Encode(o); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
