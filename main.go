package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/shorten", shortenLink).Methods(http.MethodPost)
	router.HandleFunc("/api/{shorten_link}/real", getRealLink).Methods(http.MethodGet)
	router.HandleFunc("/{shorten_link}", redirectToRealLink).Methods(http.MethodGet)
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
	{
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
