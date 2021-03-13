package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/register", register).Methods(http.MethodPost)
	router.HandleFunc("/api/login", login).Methods(http.MethodPut)
	router.HandleFunc("/api/logout", logout).Methods(http.MethodPut)
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

type registrationModel struct {
	Login string `json:"login"`
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
	Login string `json:"login"`
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
	writer.Write([]byte(token))
	writer.WriteHeader(http.StatusOK)
}

func logout(writer http.ResponseWriter, request *http.Request) {
	_, err := request.Cookie("_cookie")
	if err != http.ErrNoCookie {
		print(err, "Failed to get cookie")
	}
	http.Redirect(writer, request, "/", http.StatusFound)
}
