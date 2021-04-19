package linkrepo

import "database/sql"

type Postgres struct {
	conn *sql.DB
}

func New(conn *sql.DB) *Postgres {
	return &Postgres{conn: conn}
}

func (p *Postgres) CreateShortLink(link string, userId string) (string, error) {
	panic("implement me")
}

func (p *Postgres) GetLinkByKey(key string) (string, error) {
	panic("implement me")
}

func (p *Postgres) MakeRedirect(key string) (string, error) {
	panic("implement me")
}

func (p *Postgres) DeleteLink(key string, userId string) (string, error) {
	panic("implement me")
}

func (p *Postgres) GetUserLinks(userId string) ([]string, error) {
	panic("implement me")
}

func (p *Postgres) GetLinkStat(link string) (uint64, error) {
	panic("implement me")
}

func (p *Postgres) CreateUserLinksStorage(userId string) (string, error) {
	panic("implement me")
}
