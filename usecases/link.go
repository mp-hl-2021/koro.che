package usecases

import (
	"koro.che/linkstorage"
)

type LinkUseCasesInterface interface {
	ShortenLink(link string, userId string) (string, error)
	MakeRedirect(key string) (string, error)
	DeleteLink(link string, userId string) (string, error)
	GetRealLink(key string) (string, error)
	GetUserLinks(userId string) ([]string, error)
	GetLinkStats(key string) (LinkStat, error)
}

type LinkStat struct {
	LinkName   string `json:"linkName"`
	UseCounter uint64  `json:"useCounter"`
}

type LinkUseCases struct{
	LinkStorage linkstorage.Interface
}

// const prefix = "koro.che/"
const prefix =  "localhost:8080/"

func (l* LinkUseCases) ShortenLink(link string, userId string) (string, error) {
	var shortLink string
	var err error
	shortLink, err = l.LinkStorage.CreateShortLink(link, userId)
	return prefix + shortLink, err
}

func (l* LinkUseCases) MakeRedirect(key string) (string, error)  {
	var link string
	var err error
	link, err = l.LinkStorage.MakeRedirect(key)
	return link, err
}

func (l* LinkUseCases) DeleteLink(link string, userId string) (string, error) {
	deleteLink, err := l.LinkStorage.DeleteLink(link, userId)
	return deleteLink, err
}

func (l* LinkUseCases) GetRealLink(key string) (string, error) {
	var link string
	var err error
	link, err = l.LinkStorage.GetLinkByKey(key)
	return link, err
}

func (l* LinkUseCases) GetUserLinks(userId string) ([]string, error) {
	var links []string
	var err error
	links, err = l.LinkStorage.GetUserLinks(userId)
	return links, err
}

func (l* LinkUseCases) GetLinkStats(link string) (LinkStat, error) {
	var stat uint64
	var err error
	stat, err = l.LinkStorage.GetLinkStat(link)
	return LinkStat{link, stat}, err
}