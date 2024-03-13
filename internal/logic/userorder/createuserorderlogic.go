package userorder

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"carservice/internal/enum/payment"
	"carservice/internal/enum/userorder"
	"carservice/internal/pkg/common/errcode"
	"carservice/internal/pkg/generator/order"
	"carservice/internal/pkg/jwt"
	"carservice/internal/svc"
	"carservice/internal/types"

	dt_user_order "carservice/internal/datatypes/userorder"

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

// CreateUserOrderFeature 创建用户订单
// FIXME 创建订单 > 创建车主信息 >
func (l *CreateUserOrderLogic) CreateUserOrderFeature(req *types.CreateUserOrderReq) (*types.CreateUserOrderRep, error) {
	// 用户是否同意协议
	if req.AgreeToTerms != 1 {
		return nil, errcode.StatusForbiddenError.Lazy("用户未阅读或同意协议")
	}

	// User
	userId, err := (jwt.GetUserId(l.ctx)).(json.Number).Int64()
	if err != nil {
		return nil, errcode.InternalServerError.Lazy("解析用户 ID 时发生错误", err.Error())
	}

	// validate CarBrand and CarBrandSeries data.
	hasCar, err := l.validateUserCar(req.CarBrandId, req.CarSeriesId)
	if err != nil {
		return nil, errcode.DatabaseError.Lazy("操作数据库时发生错误", err.Error())
	}
	if !hasCar {
		return nil, errcode.NotFound.SetMsg("该车辆不存在")
	}

	// create database transaction
	tx, err := l.svcCtx.DBC.BeginTxx(l.ctx, &sql.TxOptions{})
	if err != nil {
		return nil, errcode.DatabaseError.Lazy("操作数据库时发生错误", err.Error())
	}

	// create the new user order.
	createPayload := &dt_user_order.CreatePayload{
		MemberId:         uint(userId),
		CarBrandId:       uint(req.CarBrandId),
		CarBrandSeriesId: uint(req.CarSeriesId),
		OrderNumber:      order.GenerateNumber(time.Now()),
		Comment:          req.Requirements,
		EstAmount:        0.000000,
		ActAmount:        0.000000,
		PaymentMethod:    uint8(payment.DefaultAtCreation),
		OrderStatus:      uint8(userorder.DefaultAtCreation),
		// CarOwnerInfoId:   *carOwnerInfoId, // ! deprecated
		PartnerStoreId: uint(req.PartnerStoreId), // ! deprecated
	}
	newUserOrderId, err := l.createUserOrder(tx, createPayload)
	if err != nil {
		fmt.Printf("Start Rollback.")
		if err1 := tx.Rollback(); err1 != nil { // Rollback
			return nil, errcode.DatabaseError.Lazy("数据库回滚时发生错误", err1.Error())
		}
		return nil, errcode.NewDatabaseErrorx().CreateError(err)
	}

	// update or create the info of UserOwner.
	// create CarOwnerInfo at the same time.
	_, err = l.createCarOwnerInfo(tx, uint(userId), *newUserOrderId, req)
	if err != nil {
		if err1 := tx.Rollback(); err1 != nil { // Rollback
			return nil, errcode.DatabaseError.Lazy("数据库回滚时发生错误", err1.Error())
		}
		return nil, errcode.DatabaseError.Lazy("操作数据库时发生错误", err.Error())
	}

	// filter out the ids.
	// and compute total amount, then update to user order.
	carReplacementIds := func() (result []uint) {
		for _, r := range req.CarReplacements {
			result = append(result, r.Id)
		}
		return
	}()
	// update or create the replacement items.
	err = l.createOrderItems(tx, *newUserOrderId, carReplacementIds)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		if err1 := tx.Rollback(); err1 != nil { // Rollback
			return nil, errcode.DatabaseError.Lazy("数据库回滚时发生错误", err1.Error())
		}
		return nil, errcode.DatabaseError.Lazy("数据库提交数据时发生错误", err.Error())
	}
	// todo: 发送下单成功短信
	return &types.CreateUserOrderRep{
		NewId: *newUserOrderId,
	}, nil
}

// createCarOwnerInfo 创建或更新用户车主信息
// todo: 删除订单时同时该 CarOwnerInfo 也被删除
func (l *CreateUserOrderLogic) createCarOwnerInfo(
	tx *sqlx.Tx,
	userId,
	userOrderId uint,
	req *types.CreateUserOrderReq,
) (*uint, error) {
	query := "INSERT INTO `car_owner_infos`(`user_id`, `user_order_id`, `name`, `phone_number`, `multilevel_address`, `full_address`, `longitude`, `latitude`) VALUES(?, ?, ?, ?, ?, ?, 0.0000000, 0.000000)"
	stat, err := tx.PrepareContext(l.ctx, query)
	if err != nil {
		return nil, err
	}
	rs, err := stat.ExecContext(l.ctx, userId, userOrderId, req.CarOwnerName, req.CarOwnerPhoneNumber, req.CarOwnerMultiLvAddr, req.CarOwnerFullAddress)
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
func (l *CreateUserOrderLogic) validateUserCar(carBrand, carBrandSeriesId int64) (bool, error) {
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
func (l *CreateUserOrderLogic) createUserOrder(tx *sqlx.Tx, payload *dt_user_order.CreatePayload) (*uint, error) {
	query := "INSERT INTO `user_orders`(`member_id`, `car_brand_id`, `car_brand_series_id`, `car_info_id`, `partner_store_id`, `order_number`, `order_status`, `comment`, `est_amount`, `act_amount`, `payment_method`, `created_at`, `updated_at`) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())"
	stmt, err := tx.PrepareContext(l.ctx, query)
	if err != nil {
		return nil, err
	}
	rs, err := stmt.ExecContext(
		l.ctx,
		payload.MemberId,
		payload.CarBrandId,
		payload.CarBrandSeriesId,
		0,
		// payload.CarOwnerInfoId, // ! deprecated
		payload.PartnerStoreId,
		payload.OrderNumber,
		payload.OrderStatus,
		payload.Comment,
		payload.EstAmount,
		payload.ActAmount,
		payload.PaymentMethod,
	)
	if err != nil {
		return nil, err
	}
	newId, err := rs.LastInsertId()
	if err != nil {
		return nil, err
	}
	newUintId := uint(newId)

	return &newUintId, nil
}

// createOrderItems 创建用户订单的配件项目
func (l *CreateUserOrderLogic) createOrderItems(tx *sqlx.Tx, orderId uint, carReplacementIds []uint) error {
	fmt.Println(carReplacementIds)
	for _, carReplacementId := range carReplacementIds {
		query := "INSERT INTO `order_items`(`user_order_id`, `car_replacement_id`) VALUES(?, ?)"
		_, err := tx.ExecContext(l.ctx, query, orderId, carReplacementId)
		if err != nil {
			if err1 := tx.Rollback(); err1 != nil { // Rollback
				logc.Errorf(l.ctx, "创建用户订单配件项目回滚时发生错误, err: %s\n", err.Error())
				return errcode.DatabaseRollbackErr
			}
			logc.Errorf(l.ctx, "创建用户订单配件项目时发生错误, err: %s\n", err.Error())
			return errcode.DatabaseExecuteErr
		}
	}
	return nil
}

type SmsPayload struct {
	CompanyName  string
	OrderNumber  string
	PartnerStore string
	CreatedAt    string
	// other fields...
}

func (l *CreateUserOrderLogic) sendSms(sp SmsPayload) error {
	// todo 发送短信
	return nil
}
