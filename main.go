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
	router.HandleFunc("/api/manage/link", deleteLink).Methods(http.MethodDelete)
	//8
	router.HandleFunc("/api/manage/links", getUserLinks).Methods(http.MethodGet)
	//9
	router.HandleFunc("/api/manage/stats", getUserLinkStats).Methods(http.MethodGet)

	server := http.Server{
		Addr:         "localhost:8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,

		Handler: router,
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
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