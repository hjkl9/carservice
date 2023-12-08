package data

import (
	"github.com/jmoiron/sqlx"
)

type PingRepo interface {
	EchoAsResult(string) (string, error)
}

type ping struct {
	db *sqlx.DB
}

func newPing(db *sqlx.DB) PingRepo {
	return &ping{db}
}

func (p *ping) EchoAsResult(value string) (string, error) {
	query := "SELECT ? AS `result`"
	var result string
	err := p.db.Get(&result, query, value)
	if err != nil {
		return "", err
	}

	return result, nil
}
