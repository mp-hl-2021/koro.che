package link

import "errors"

var (
	ErrNotExist = errors.New("link does not exist")
)

type Interface interface {
	CreateShortLink(link string, userId string) (string, error)
	GetLinkByKey(key string) (string, error)
	MakeRedirect(key string) (string, error)
	DeleteLink(key string, userId string) (string, error)
	GetUserLinks(userId string) ([]string, error)
	GetLinkStat(link string) (uint64, error)
	CreateUserLinksStorage(userId string) (string, error)
}