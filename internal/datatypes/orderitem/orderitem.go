package orderitem

import "database/sql/driver"

type CreateItem struct {
	UserOrderId      uint `db:"user_order_id"`
	CarReplacementId uint `db:"car_replacement_id"`
}

func (cp CreateItem) Value() (driver.Value, error) {
	return []interface{}{cp.UserOrderId, cp.CarReplacementId}, nil
}
