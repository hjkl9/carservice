package data

import (
	"carservice/internal/datatypes/orderitem"
	"context"

	"github.com/jmoiron/sqlx"
)

type OrderItemRepo interface {
	BatchCreate(context.Context, uint, []orderitem.CreateItem) error
	DeleteByOrderId(context.Context, uint) error
}

type orderItem struct {
	db *sqlx.DB
}

func newOrderItem(db *sqlx.DB) *orderItem {
	return &orderItem{db}
}

func (oi *orderItem) BatchCreate(ctx context.Context, orderId uint, replacements []orderitem.CreateItem) error {

	query := "INSERT INTO `order_items`(`user_order_id`, `car_replacement_id`) VALUES(:user_order_id, :car_replacement_id)"

	_, err := oi.db.NamedExecContext(ctx, query, replacements)
	if err != nil {
		return err
	}

	return nil
}

func (oi *orderItem) DeleteByOrderId(ctx context.Context, orderId uint) error {
	query := "DELETE FROM `order_items` WHERE `user_order_id` = ?"

	_, err := oi.db.ExecContext(ctx, query, orderId)
	if err != nil {
		return err
	}

	return nil
}
