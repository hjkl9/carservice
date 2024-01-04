package userorder

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"carservice/internal/data/tables"
	"carservice/internal/enum/payment"
	"carservice/internal/enum/userorder"
	"carservice/internal/pkg/common/errcode"
	"carservice/internal/pkg/generator/order"
	"carservice/internal/pkg/jwt"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/jmoiron/sqlx"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateUserOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateUserOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserOrderLogic {
	return &CreateUserOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateUserOrderLogic) CreateUserOrder(req *types.CreateUserOrderReq) error {
	// ! 注意 user_id 和 member_id 都是用户的 ID
	userId := jwt.GetUserId(l.ctx)

	// 创建事务
	tx, err := l.svcCtx.DBC.Beginx()
	if err != nil {
		logc.Errorf(l.ctx, "创建订单->创建事务时发生错误, err: %s\n", err.Error())
		return errcode.DatabaseError.SetMsg("创建事务时发生错误").SetDetails(err.Error())
	}

	// 创建车主信息
	query := "INSERT INTO `%s`(`name`, `user_id`, `phone_number`, `multilevel_address`, `full_address`, `longitude`, `latitude`) VALUES(?, ?, ?, ?, ?, ?, ?)"
	result, err := tx.ExecContext(
		l.ctx,
		fmt.Sprintf(query, tables.CarOwnerInfo),
		req.CarOwnerName, userId,
		req.CarOwnerPhoneNumber,
		req.CarOwnerMultiLvAddr,
		req.CarOwnerFullAddress,
		req.CarOwnerLongitude,
		req.CarOwnerLatitude,
	)
	if err != nil {
		// 回滚
		if err1 := tx.Rollback(); err1 != nil {
			return errcode.DatabaseError.SetMsg("回滚失败").SetDetails(err1.Error())
		}
		logc.Errorf(l.ctx, "创建订单->创建车主信息时发生错误, err: %s\n", err.Error())
		return errcode.DatabaseError.SetMsg("创建用户车主信息时发生错误").SetDetails(err.Error())
	}
	// 新创建的车主信息 ID
	newCarOwnerInfoId, _ := result.LastInsertId()

	// 品牌和车系
	// 查询是否存在品牌和车系
	query = "SELECT count(1) AS `count` FROM `%s` AS `cb` JOIN `%s` AS `cbs` ON `cb`.`brand_id` = `cbs`.`brand_id` WHERE `cb`.`brand_id` = %d AND `cbs`.`series_id` = %d"
	var count int
	if err = tx.Get(
		&count,
		fmt.Sprintf(query, tables.CarBrand, tables.CarBrandSeries, req.CarBrandId, req.CarBrandSeriesId),
	); err != nil {
		// 回滚
		if err1 := tx.Rollback(); err1 != nil {
			return errcode.DatabaseError.SetMsg("数据库回滚时发生错误").SetDetails(err1.Error())
		}
		logc.Errorf(l.ctx, "创建订单->查询车数据时发生错误, err: %s\n", err.Error())
		return errcode.DatabaseError.SetMsg("数据库查询时发生错误").SetDetails(err.Error())
	}
	// 如果不存在
	if count == 0 {
		if err1 := tx.Rollback(); err1 != nil {
			return errcode.DatabaseError.SetMsg("数据库回滚时发生错误").SetDetails(err1.Error())
		}
		return errcode.InvalidParamsError.SetMsg("无效的汽车品牌和系列")
	}
	// 安装需求
	// req.Requirements
	// 订单默认状态
	orderStatus := userorder.Pending
	// 订单默认支付方式
	paymentMethod := payment.Unknown
	// 生成订单号
	orderNumber := order.GenerateNumber(time.Now())
	// 创建时间和更新时间
	nowUnix := time.Now()
	var createdAt = nowUnix
	var updatedAt = nowUnix
	// 创建订单
	query = "INSERT INTO `%s`(`member_id`, `car_brand_id`, `car_brand_series_id`, `car_info_id`, `car_owner_info_id`, `partner_store_id`, `order_number`, `order_status`, `comment`, `est_amount`, `act_amount`, `payment_method`, `created_at`, `updated_at`) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	rs, err := tx.ExecContext(
		l.ctx,                                // 上下文
		fmt.Sprintf(query, tables.UserOrder), // 创建语句
		userId,                               // 用户
		req.CarBrandId,                       // 汽车品牌
		req.CarBrandSeriesId,                 // 汽车品牌系列
		0,                                    // ! car_info_id 应该删除该字段
		newCarOwnerInfoId,                    // 用户车主信息
		req.PartnerStoreId,                   // ! partner_store_id 合作门店 ID
		orderNumber,                          // 订单号
		orderStatus,                          // 订单状态
		req.Requirements,                     // 需求
		0.00,                                 // ! unused 默认的预估服务费
		0.00,                                 // ! unused 默认的实际服务费
		paymentMethod,                        // 默认的支付方式
		createdAt,                            // 创建时间
		updatedAt,                            // 更新时间
	)
	if err != nil {
		// 回滚
		if err1 := tx.Rollback(); err1 != nil {
			return errcode.DatabaseError.SetMsg("数据库回滚时发生错误").SetDetails(err1.Error())
		}
		logc.Errorf(l.ctx, "创建订单时发生错误, err: %s\n", err.Error())
		return errcode.DatabaseError.SetMsg("创建订单时发生错误").SetDetails(err.Error())
	}

	rows, _ := rs.RowsAffected()
	if rows == 0 {
		// 回滚
		if err1 := tx.Rollback(); err1 != nil {
			return errcode.DatabaseError.SetMsg("数据库回滚时发生错误").SetDetails(err1.Error())
		}
		logc.Errorf(l.ctx, "创建订单->创建订单失败, err: %s\n", err.Error())
		return errcode.InternalServerError.SetMsg("订单创建失败")
	}

	// 提交事务
	if err = tx.Commit(); err != nil {
		return errcode.DatabaseError.SetMsg("提交事务事务时发生错误").SetDetails(err.Error())
	}

	return nil
}

// Logic related structures. //
type carOwnerInfoCounter struct {
	Count   uint `db:"count"`
	FirstId uint `db:"firstId"`
}

// createUserOrderPayload 创建用户订单数据
// 内存对齐 OK
type createUserOrderPayload struct {
	MemberId         uint    `db:"member_id"`
	CarBrandId       uint    `db:"car_brand_id"`
	CarBrandSeriesId uint    `db:"car_brand_series_id"`
	CarOwnerInfoId   uint    `db:"car_owner_info_id"`
	PartnerStoreId   uint    `db:"partner_store_id"`
	OrderNumber      string  `db:"order_number"`
	Comment          string  `db:"comment"`
	EstAmount        float64 `db:"est_amount"`
	ActAmount        float64 `db:"act_amount"`
	PaymentMethod    uint8   `db:"payment_method"`
	OrderStatus      uint8   `db:"order_status"`
	// CreatedAt        time.Duration `db:"created_at"`
	// UpdatedAt        time.Duration `db:"updated_at"`
}

func (l *CreateUserOrderLogic) CreateUserOrderFeature(req *types.CreateUserOrderReq) error {
	// User
	userId, err := (jwt.GetUserId(l.ctx)).(json.Number).Int64()
	if err != nil {
		return errcode.InternalServerError.Lazy("解析用户 ID 时发生错误", err.Error())
	}

	// check if the PartnerStore already exists.
	hasPartnerStore, err := l.validatePartnerStore(req.PartnerStoreId)
	if err != nil {
		return errcode.NewDatabaseErrorx().GetError(err)
	}
	if !hasPartnerStore {
		return errcode.NotFound.SetMsg("该合作门店不存在")
	}

	// validate CarBrand and CarBrandSeries data.
	hasCar, err := l.validateUserCar(req.CarBrandId, req.CarBrandSeriesId)
	if err != nil {
		return errcode.DatabaseError.Lazy("操作数据库时发生错误", err.Error())
	}
	if !hasCar {
		return errcode.NotFound.SetMsg("该车辆不存在")
	}

	// create database transaction
	tx, err := l.svcCtx.DBC.BeginTxx(l.ctx, &sql.TxOptions{})
	if err != nil {
		return errcode.DatabaseError.Lazy("操作数据库时发生错误", err.Error())
	}

	// update or create the info of UserOwner.
	// create CarOwnerInfo at the same time.
	carOwnerInfoId, err := l.createCarOwnerInfo(tx, uint(userId), req)
	if err != nil {
		if err1 := tx.Rollback(); err1 != nil { // Rollback
			return errcode.DatabaseError.Lazy("数据库回滚时发生错误", err1.Error())
		}
		return errcode.DatabaseError.Lazy("操作数据库时发生错误", err.Error())
	}

	// create the new user order.
	createPayload := &createUserOrderPayload{
		MemberId:         uint(userId),
		CarBrandId:       req.CarBrandId,
		CarBrandSeriesId: req.CarBrandSeriesId,
		CarOwnerInfoId:   *carOwnerInfoId,
		PartnerStoreId:   req.PartnerStoreId,
		OrderNumber:      order.GenerateNumber(time.Now()),
		Comment:          req.Requirements,
		EstAmount:        0.000000,
		ActAmount:        0.000000,
		PaymentMethod:    uint8(payment.DefaultAtCreation),
		OrderStatus:      uint8(userorder.DefaultAtCreation),
	}
	if err = l.createUserOrder(tx, createPayload); err != nil {
		fmt.Printf("Start Rollback.")
		if err1 := tx.Rollback(); err1 != nil { // Rollback
			return errcode.DatabaseError.Lazy("数据库回滚时发生错误", err1.Error())
		}
		return errcode.NewDatabaseErrorx().CreateError(err)
	}
	if err = tx.Commit(); err != nil {
		if err1 := tx.Rollback(); err1 != nil { // Rollback
			return errcode.DatabaseError.Lazy("数据库回滚时发生错误", err1.Error())
		}
		return errcode.DatabaseError.Lazy("数据库提交数据时发生错误", err.Error())
	}
	return nil
}

// createCarOwnerInfo 创建或更新用户车主信息
// todo: 删除订单时同时该 CarOwnerInfo 也被删除
func (l *CreateUserOrderLogic) createCarOwnerInfo(
	tx *sqlx.Tx,
	userId uint,
	req *types.CreateUserOrderReq,
) (*uint, error) {
	/// ! Should be deleted. ///
	// Check if the car owner info was exists.
	// var counter carOwnerInfoCounter
	// query := "SELECT COUNT(1) AS `count`, MIN(`id`) AS `firstId` FROM `car_owner_infos` WHERE `user_id` = ? LIMIT 1"
	// stmtx, err := tx.PreparexContext(l.ctx, query)
	// if err != nil {
	// 	return nil, err
	// }
	// if err = stmtx.GetContext(l.ctx, &counter, userId); err != nil {
	// 	return nil, err
	// }
	// // If doesn't exist then create a new one.
	// if counter.Count == 0 {
	// }
	// // Otherwise update and return id of CarOwnerInfo.
	// query = "UPDATE `car_owner_infos` SET `name` = ?, `phone_number` = ?, `multilevel_address` = ?, `full_address` = ? WHERE `user_id` = ? AND `id` = ?"
	// stmt, err := tx.PrepareContext(l.ctx, query)
	// if err != nil {
	// 	return nil, err
	// }
	// _, err = stmt.ExecContext(l.ctx, req.CarOwnerName, req.CarOwnerPhoneNumber, req.CarOwnerMultiLvAddr, req.CarOwnerFullAddress, userId, counter.FirstId)
	// if err != nil {
	// 	return nil, err
	// }
	/// ! Should be deleted. ///

	query := "INSERT INTO `car_owner_infos`(`user_id`, `name`, `phone_number`, `multilevel_address`, `full_address`, `longitude`, `latitude`) VALUES(?, ?, ?, ?, ?, 0.0000000, 0.000000)"
	stat, err := tx.PrepareContext(l.ctx, query)
	if err != nil {
		return nil, err
	}
	rs, err := stat.ExecContext(l.ctx, userId, req.CarOwnerName, req.CarOwnerPhoneNumber, req.CarOwnerMultiLvAddr, req.CarOwnerFullAddress)
	if err != nil {
		return nil, err
	}
	newId, err := rs.LastInsertId()
	newUintId := uint(newId)
	if err != nil {
		return nil, err
	}
	return &newUintId, nil
}

// validateUserCar 验证车辆信息
func (l *CreateUserOrderLogic) validateUserCar(carBrand, carBrandSeriesId uint) (bool, error) {
	var count uint8
	query := "SELECT COUNT(1) AS `count` FROM `car_brands` `cb` JOIN `car_brand_series` `cbs` ON `cb`.`brand_id` = `cbs`.`brand_id` WHERE `cbs`.`brand_id` = ? AND `cbs`.`series_id` = ? LIMIT 1;"
	stmt, err := l.svcCtx.DBC.PreparexContext(l.ctx, query)
	if err != nil {
		return false, err
	}
	if err = stmt.GetContext(l.ctx, &count, carBrand, carBrandSeriesId); err != nil {
		return false, err
	}
	return count == 1, nil
}

// createUserOrder 创建用户订单
func (l *CreateUserOrderLogic) createUserOrder(tx *sqlx.Tx, payload *createUserOrderPayload) error {
	query := "INSERT INTO `user_orders`(`member_id`, `car_brand_id`, `car_brand_series_id`, `car_info_id`, `car_owner_info_id`, `partner_store_id`, `order_number`, `order_status`, `comment`, `est_amount`, `act_amount`, `payment_method`, `created_at`, `updated_at`) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())"
	stmt, err := tx.PrepareContext(l.ctx, query)
	if err != nil {
		return err
	}
	rs, err := stmt.ExecContext(
		l.ctx,
		payload.MemberId,
		payload.CarBrandId,
		payload.CarBrandSeriesId,
		0,
		payload.CarOwnerInfoId,
		payload.PartnerStoreId,
		payload.OrderNumber,
		payload.OrderStatus,
		payload.Comment,
		payload.EstAmount,
		payload.ActAmount,
		payload.PaymentMethod,
	)
	if err != nil {
		return err
	}
	n, err := rs.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return err
	}

	return nil
}

func (l *CreateUserOrderLogic) validatePartnerStore(partnerStoreId uint) (bool, error) {
	var count uint8
	query := "SELECT COUNT(1) AS `count` FROM `partner_stores` WHERE `id` = ? LIMIT 1"
	stmt, err := l.svcCtx.DBC.PreparexContext(l.ctx, query)
	if err != nil {
		return false, err
	}
	if err = stmt.GetContext(l.ctx, &count, partnerStoreId); err != nil {
		return false, err
	}
	return count == 1, nil
}
