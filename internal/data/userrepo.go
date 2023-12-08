package data

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type UserRepo interface {
	GetById() error
	GetByPhoneNumber(string) error
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

func (u *user) GetByPhoneNumber(phoneNumber string) error {
	fmt.Println("Get user by phone number.")
	// getting logic here...
	return nil
}
