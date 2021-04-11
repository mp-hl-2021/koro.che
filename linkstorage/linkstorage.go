package linkstorage

import "errors"

var (
	ErrNotExist = errors.New("lin does not exist")
)

type Interface interface {
	CreateShortLink(link string) (string, error)
	GetLinkByKey(key string) (string, error)
	MakeRedirect(key string) (string, error)
}