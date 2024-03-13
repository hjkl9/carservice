package data

import (
	"carservice/internal/datatypes/userorder"
	"context"

	"github.com/jmoiron/sqlx"
)

type UserOrderRepo interface {
	GetIfOrderExistsById(ctx context.Context, memberId interface{}, orderId uint) (bool, error)
	GetFacadesOrderById(ctx context.Context, memberId interface{}, orderId uint) (*userorder.Facades_UserOrder, error)
	GetOrderById(ctx context.Context, memberId interface{}, orderId uint) (*userorder.Native_UserOrder, error)
	GetIfTheListExists(ctx context.Context, memberId uint) (bool, error)
	GetOrderList(ctx context.Context, memberId uint) (*[]*userorder.UserOrderListItem, error)
	SoftDeleteOrderById(ctx context.Context, memberId, orderId uint) error
	Update(ctx context.Context, id uint, payload *userorder.UpdatePayload) error
}

type userOrder struct {
	db *sqlx.DB
}

func newUserOrder(db *sqlx.DB) *userOrder {
	return &userOrder{db}
}

func (uo *userOrder) GetIfOrderExistsById(ctx context.Context, memberId interface{}, orderId uint) (bool, error) {
	var hasOrder uint8
	query := "SELECT (COUNT(1) = 1) AS `hasOrder` FROM `user_orders` WHERE `member_id` = ? AND `id` = ? AND `deleted_at` IS NULL LIMIT 1"
	stmt, err := uo.db.PreparexContext(ctx, query)
	if err != nil {
		return false, err
	}
	if err = stmt.GetContext(ctx, &hasOrder, memberId, orderId); err != nil {
		return false, err
	}
	return hasOrder == 1, nil
}

// func (uo *userOrder) GetOrderById(ctx context.Context, memberId interface{}, orderId uint) (*userorder.UserOrder, error) {
// 	var order userorder.UserOrder
// 	query := "SELECT `uo`.`id` AS `id`, `uo`.`order_number` AS `orderNumber`, `coi`.`name` AS `carOwnerName`, `coi`.`multilevel_address` AS `carOwnerMultiLvAddr`, `coi`.`full_address` AS `carOwnerFullAddress` , `ps`.`title` AS `partnerStore`, `ps`.`full_address` AS `partnerStoreAddress`, `uo`.`comment` AS `comment`, `cb`.`brand_name` AS `carBrandName`, `cbs`.`series_name` AS `carSeriesName` , `uo`.`order_status` AS `orderStatus`, `uo`.`created_at` AS `createdAt`, `uo`.`updated_at` AS `updatedAt` FROM `user_orders` `uo` LEFT JOIN `partner_stores` `ps` ON `ps`.`id` = `uo`.`partner_store_id` JOIN `car_owner_infos` `coi` ON `coi`.`id` = `uo`.`car_owner_info_id` JOIN `car_brands` `cb` ON `cb`.`brand_id` = `uo`.`car_brand_id` JOIN `car_brand_series` `cbs` ON `cbs`.`series_id` = `uo`.`car_brand_series_id` WHERE `uo`.`id` = ? AND `uo`.`member_id` = ? LIMIT 1"
// 	stmt, err := uo.db.PreparexContext(ctx, query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if err = stmt.GetContext(ctx, &order, orderId, memberId); err != nil {
// 		return nil, err
// 	}
// 	return &order, nil
// }

func (uo *userOrder) GetFacadesOrderById(ctx context.Context, memberId interface{}, orderId uint) (*userorder.Facades_UserOrder, error) {
	var order userorder.Facades_UserOrder
	query := "SELECT `uo`.`id` AS `id`, `uo`.`order_number` AS `orderNumber`, `coi`.`name` AS `carOwnerName`, `coi`.`multilevel_address` AS `carOwnerMultiLvAddr`, `coi`.`full_address` AS `carOwnerFullAddress` , `ps`.`title` AS `partnerStore`, `ps`.`full_address` AS `partnerStoreAddress`, `uo`.`comment` AS `comment`, `cb`.`brand_name` AS `carBrandName`, `cbs`.`series_name` AS `carSeriesName` , `uo`.`order_status` AS `orderStatus`, `uo`.`created_at` AS `createdAt`, `uo`.`updated_at` AS `updatedAt` FROM `user_orders` `uo` LEFT JOIN `partner_stores` `ps` ON `ps`.`id` = `uo`.`partner_store_id` JOIN `car_owner_infos` `coi` ON `coi`.`id` = `uo`.`car_owner_info_id` JOIN `car_brands` `cb` ON `cb`.`brand_id` = `uo`.`car_brand_id` JOIN `car_brand_series` `cbs` ON `cbs`.`series_id` = `uo`.`car_brand_series_id` WHERE `uo`.`id` = ? AND `uo`.`member_id` = ? LIMIT 1"
	stmt, err := uo.db.PreparexContext(ctx, query)
	if err != nil {
		return nil, err
	}
	if err = stmt.GetContext(ctx, &order, orderId, memberId); err != nil {
		return nil, err
	}
	return &order, nil
}

func (uo *userOrder) GetOrderById(ctx context.Context, memberId interface{}, orderId uint) (*userorder.Native_UserOrder, error) {
	var order userorder.Native_UserOrder
	query := "SELECT `id`, `car_brand_id`, `car_brand_series_id`, `member_id`, `car_info_id`, `order_number`, `est_amount`, `act_amount`, `expired_at`, `payment_method`, `paid_at`, `order_status`, `comment`, `created_at`, `updated_at`, `partner_store_id`, `deleted_at`, `installer_phone_number`, `installer_name` FROM `user_orders` WHERE `id` = ? AND `member_id` = ? LIMIT 1;"
	stmt, err := uo.db.PreparexContext(ctx, query)
	if err != nil {
		return nil, err
	}
	if err = stmt.GetContext(ctx, &order, orderId, memberId); err != nil {
		return nil, err
	}
	return &order, nil
}

func (uo *userOrder) GetIfTheListExists(ctx context.Context, memberId uint) (bool, error) {
	var hasList uint8
	query := "SELECT (COUNT(1) > 0) AS `hasList` FROM `user_orders` WHERE `member_id` = ? AND `deleted_at` IS NULL;"
	stmt, err := uo.db.PreparexContext(ctx, query)
	if err != nil {
		return false, err
	}
	if err = stmt.GetContext(ctx, &hasList, memberId); err != nil {
		return false, err
	}
	return hasList == 1, nil
}

func (uo *userOrder) GetOrderList(ctx context.Context, memberId uint) (*[]*userorder.UserOrderListItem, error) {
	var items []*userorder.UserOrderListItem
	query := "SELECT `uo`.`id`, `uo`.`order_number` AS `orderNumber`, `ps`.`title` AS `partnerStore`, `uo`.`comment` AS `requirements`, `uo`.`order_status` AS `orderStatus`, `uo`.`created_at` AS `createdAt`, `uo`.`updated_at` AS `updatedAt` FROM `%s` AS `uo` LEFT JOIN `%s` AS `ps` ON `uo`.`partner_store_id` = `ps`.`id` WHERE `uo`.`member_id` = ? AND `uo`.`deleted_at` IS NULL"
	stmt, err := uo.db.PreparexContext(ctx, query)
	if err != nil {
		return nil, err
	}
	if err = stmt.SelectContext(ctx, &items, memberId); err != nil {
		return nil, err
	}
	return &items, nil
}

func (uo *userOrder) SoftDeleteOrderById(ctx context.Context, memberId, orderId uint) error {
	query := "UPDATE `user_orders` SET `deleted_at` = NOW() WHERE `member_id` = ? AND `id` = ?"
	stmt, err := uo.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	rs, err := stmt.ExecContext(ctx, memberId, orderId)
	if err != nil {
		return err
	}
	_, err = rs.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}

func (uo *userOrder) Update(ctx context.Context, id uint, payload *userorder.UpdatePayload) error {

	return nil
}
