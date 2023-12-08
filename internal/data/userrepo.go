package data

import "github.com/jmoiron/sqlx"

type UserRepo interface {
	GetById() error
	GetByPhoneNumber() error
}

type user struct {
	db *sqlx.DB
}

func newUser(dbc *sqlx.DB) *user {
	return &user{dbc}
}

func (u *user) GetById() error {
	return nil
}

func (u *user) GetByPhoneNumber() error {
	return nil
}
