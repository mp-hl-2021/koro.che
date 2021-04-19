package accountrepo

import (
	"database/sql"
	"koro.che/internal/domain/account"
)

type Postgres struct {
	conn *sql.DB
}

func New(conn *sql.DB) *Postgres {
	return &Postgres{conn: conn}
}

func (p *Postgres) CreateAccount(cred account.Credentials) (account.Account, error) {
	panic("implement me")
}

func (p *Postgres) GetAccountById(id string) (account.Account, error) {
	panic("implement me")
}

func (p *Postgres) GetAccountByLogin(login string) (account.Account, error) {
	panic("implement me")
}
