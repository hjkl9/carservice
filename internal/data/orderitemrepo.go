package data

import "github.com/jmoiron/sqlx"

type OrderItemRepo interface {
}

type orderItem struct {
	db *sqlx.DB
}

func (oi *orderItem) _() {
}
