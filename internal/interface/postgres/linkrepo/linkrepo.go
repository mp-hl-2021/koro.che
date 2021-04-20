package linkrepo

import (
	"database/sql"
	link2 "koro.che/internal/domain/link"
	"math/rand"
)

type Postgres struct {
	conn *sql.DB
}

func New(conn *sql.DB) *Postgres {
	return &Postgres{conn: conn}
}

const allowedLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const linkSize = 6

func RandString() string {
	b := make([]byte, linkSize)
	for i := range b {
		b[i] = allowedLetters[rand.Intn(len(allowedLetters))]
	}
	return string(b)
}

type LinkInfo struct {
	Id         int
	CreatorId  string
	RealLink   string
	Key        string
	UseCounter int
}

const queryCreateLink = `
	insert into 
	    links(creator_id, real_link, key) 
	    values ($1, $2, $3)
`

const queryGetRealLinkByKey = `
	select real_link from links 
	where key = $1
`
const queryIncreaseLinkStat = `
	update links
		set use_counter = use_counter + 1
	where key = $1
`

const queryDeleteLink = `
	delete from links
	where key = $1
`

const queryUserLinks = `
	select key from links
	where creator_id = $1
`

const queryLinkStats = `
	select use_counter from links
	where key = $1 
`

func (p *Postgres) CreateShortLink(link string, userId string) (string, error) {
	var key string
	for {
		key = RandString()
		row := p.conn.QueryRow(queryGetRealLinkByKey, key)
		err := row.Scan()
		if err == sql.ErrNoRows {
			p.conn.QueryRow(queryCreateLink, userId, link, key)
		    // todo need wrapping???
			break
		}
	}
	return key, nil
}

func (p *Postgres) GetLinkByKey(key string) (string, error) {
	var realLink string
	row := p.conn.QueryRow(queryGetRealLinkByKey, key)
	err := row.Scan(&realLink)
	if err != nil && err == sql.ErrNoRows {
		return "", link2.ErrNotExist
	}
	return realLink, err
}

func (p *Postgres) MakeRedirect(key string) (string, error) {
	var realLink string
	row := p.conn.QueryRow(queryGetRealLinkByKey, key)
	err := row.Scan(&realLink)

	if err != nil && err == sql.ErrNoRows {
		return "", link2.ErrNotExist
	}

	row = p.conn.QueryRow(queryIncreaseLinkStat, key)

	return realLink, err
}

func (p *Postgres) DeleteLink(key string, userId string) (string, error) {
	var realLink string
	row := p.conn.QueryRow(queryGetRealLinkByKey, key)
	err := row.Scan(&realLink)
	if err != nil && err == sql.ErrNoRows {
		return "", link2.ErrNotExist
	}
	if err != nil {
		return "", err // todo wrapping
	}
	row = p.conn.QueryRow(queryDeleteLink, key)

	return realLink, err
}

func (p *Postgres) GetUserLinks(userId string) ([]string, error) {
	var userLinks = make([]string, 0)
	rows, err := p.conn.Query(queryUserLinks, userId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var link string
		if err := rows.Scan(&link); err != nil {
			return nil, err
		}
		userLinks = append(userLinks, link)
	}
	return userLinks, nil
}

func (p *Postgres) GetLinkStat(key string) (uint64, error) {
	var stat uint64
	row := p.conn.QueryRow(queryLinkStats, key)
	err := row.Scan(&stat) //todo wrapping
	return stat, err
}

func (p *Postgres) CreateUserLinksStorage(userId string) (string, error) {
	return "", nil
}
