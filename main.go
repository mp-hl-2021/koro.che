package main

import (
	"flag"
	"io/ioutil"
	"koro.che/accountstorage"
	"koro.che/api"
	"koro.che/auth"
	"koro.che/linkstorage"
	"koro.che/usecases"
	"net/http"
	"time"
)

func main() {

	privateKeyPath := flag.String("privateKey", "app.rsa", "file path")
	publicKeyPath := flag.String("publicKey", "app.rsa.pub", "file path")
	flag.Parse()

	privateKeyBytes, err := ioutil.ReadFile(*privateKeyPath)
	publicKeyBytes, err := ioutil.ReadFile(*publicKeyPath)
	a, err := auth.NewToken(privateKeyBytes, publicKeyBytes, 100*time.Minute)
	if err != nil {
		panic(err)
	}
	accountUseCases := usecases.AccountUseCases{
		AccountStorage: accountstorage.NewMemory(),
		Auth:           a,
	}
	linkUseCases := usecases.LinkUseCases{
		LinkStorage: linkstorage.NewMemory(),
	}
	service := api.NewApi(&accountUseCases, &linkUseCases)

	server := http.Server{
		Addr:         "localhost:8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      service.Router(),
	}
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
