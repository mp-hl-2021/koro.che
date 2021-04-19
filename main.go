package main

import (
	"flag"
	"io/ioutil"
	auth2 "koro.che/internal/auth"
	"koro.che/internal/interface/httpapi"
	"koro.che/internal/interface/memory/accountrepo"
	"koro.che/internal/interface/memory/linkrepo"
	"koro.che/internal/usecases/account"
	"koro.che/internal/usecases/link"
	"net/http"
	"time"
)

func main() {

	privateKeyPath := flag.String("privateKey", "app.rsa", "file path")
	publicKeyPath := flag.String("publicKey", "app.rsa.pub", "file path")
	flag.Parse()

	privateKeyBytes, err := ioutil.ReadFile(*privateKeyPath)
	publicKeyBytes, err := ioutil.ReadFile(*publicKeyPath)
	a, err := auth2.NewToken(privateKeyBytes, publicKeyBytes, 100*time.Minute)
	if err != nil {
		panic(err)
	}
	accountUseCases := account.AccountUseCases{
		AccountStorage: accountrepo.NewMemory(),
		Auth:           a,
	}
	linkUseCases := link.LinkUseCases{
		LinkStorage: linkrepo.NewMemory(),
	}
	service := httpapi.NewApi(&accountUseCases, &linkUseCases)

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
