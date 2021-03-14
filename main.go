package main

import (
	"koro.che/api"
	"koro.che/usecases"
	"net/http"
	"time"
)

func main() {
	service := api.NewApi(&usecases.AccountUseCases{}, &usecases.LinkUseCases{})

	server := http.Server{
		Addr:         "localhost:8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      service.Router(),
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
