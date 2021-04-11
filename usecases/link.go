package usecases

import (
	"koro.che/linkstorage"
)

type LinkUseCasesInterface interface {
	ShortenLink(link string) (string, error)
	MakeRedirect(key string) (string, error)
	DeleteLink(link string) ()
	GetRealLink(key string) (string, error)
	GetUserLinks(userId string) ([]string, error)
	GetLinkStats(key string) ([]LinkStat, error)
}

type LinkStat struct {
	LinkName   string `json:"linkName"`
	UseCounter int64  `json:"useCounter"`
}

type LinkUseCases struct{
	LinkStorage linkstorage.Interface
}

// const prefix = "koro.che/"
const prefix =  "localhost:8080/"

func (l* LinkUseCases) ShortenLink(link string) (string, error) {
	var shortLink string
	var err error
	shortLink, err = l.LinkStorage.CreateShortLink(link)
	return prefix + shortLink, err
}

func (l* LinkUseCases) MakeRedirect(key string) (string, error)  {
	var link string
	var err error
	link, err = l.LinkStorage.MakeRedirect(key)
	return link, err
}

func (l* LinkUseCases) DeleteLink(link string) () {
	panic("not implemented method")
}

func (l* LinkUseCases) GetRealLink(key string) (string, error) {
	var link string
	var err error
	link, err = l.LinkStorage.GetLinkByKey(key)
	return link, err
}

func (l* LinkUseCases) GetUserLinks(userId string) ([]string, error) {
	panic("not implemented method")
}

func (l* LinkUseCases) GetLinkStats(link string) ([]LinkStat, error) {
	panic("not implemented method")
}