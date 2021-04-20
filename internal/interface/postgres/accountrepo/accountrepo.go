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
	aсс := account.Account{Credentials: cred}
	row := p.conn.QueryRow(`
		INSERT INTO accounts(login, password) VALUES ($1, $2)
		RETURNING id`, cred.Login, cred.Password)
	err := row.Scan(&aсс.Id)
	return aсс, err
}

func (p *Postgres) GetAccountById(id string) (account.Account, error) {
	aсс := account.Account{}
	row := p.conn.QueryRow(`
		SELECT id, login, password
		FROM accounts
		WHERE id = $1`, id)
	err := row.Scan(&aсс.Id, &aсс.Login, &aсс.Password)
	return aсс, err
}

func (p *Postgres) GetAccountByLogin(login string) (account.Account, error) {
	aсс := account.Account{}
	row := p.conn.QueryRow(`
		SELECT id, login, password
		FROM accounts
		WHERE login = $1`, login)
	err := row.Scan(&aсс.Id, &aсс.Login, &aсс.Password)
	return aсс, err
}
