package usecases

type LinkUseCasesInterface interface {
	ShortenLink(link string) (string, error)
	DeleteLink(link string) ()
	GetRealLink(shortLink string) (string, error)
	GetUserLinks(userId string) ([]string, error)
	GetLinkStats(link string) ([]LinkStat, error)
}

type LinkStat struct {
	LinkName   string `json:"linkName"`
	UseCounter int64  `json:"useCounter"`
}

type LinkUseCases struct{}

func (LinkUseCases) ShortenLink(link string) (string, error) {
	panic("not implemented method")
}

func (LinkUseCases) DeleteLink(link string) () {
	panic("not implemented method")
}

func (LinkUseCases) GetRealLink(shortLink string) (string, error) {
	panic("not implemented method")
}

func (LinkUseCases) GetUserLinks(userId string) ([]string, error) {
	panic("not implemented method")
}

func (LinkUseCases) GetLinkStats(link string) ([]LinkStat, error) {
	panic("not implemented method")
}