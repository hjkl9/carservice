package userorder

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"carservice/internal/pkg/common/errcode"
	"carservice/internal/pkg/jwt"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/jmoiron/sqlx"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserOrderLogic {
	return &UpdateUserOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserOrderLogic) UpdateUserOrder(req *types.UpdateUserOrderReq) error {
	userId, err := (jwt.GetUserId(l.ctx).(json.Number)).Int64()
	if err != nil {
		return errcode.InternalServerError.Lazy("UserID 类型转换时发生错误").SetDetails(err.Error())
	}

	hasOrder, err := l.svcCtx.Repo.UserOrder().GetIfOrderExistsById(l.ctx, userId, req.Id)
	if err != nil {
		return errcode.DatabaseGetErr
	}
	if !hasOrder {
		return errcode.InvalidParametersErr.SetMessage("无效的订单")
	}

	// validate CarBrand and CarBrandSeries.
	hasCarSeries, err := l.svcCtx.Repo.
		CarBrandSeriesRepo().
		CheckIfSeriesExists(l.ctx, req.CarBrandId, req.CarSeriesId)
	if err != nil {
		logc.Errorf(l.ctx, "检查车型是否存在时发生错误, err: %s\n", err.Error())
		return errcode.DatabaseGetErr
	}
	if !hasCarSeries {
		return errcode.InvalidParametersErr.SetMessage("无效的车型")
	}

	// create a database transaction.
	tx, err := l.svcCtx.DBC.BeginTxx(l.ctx, &sql.TxOptions{})
	if err != nil {
		logc.Errorf(l.ctx, "创建事务时发生错误, err: %s\n", err.Error())
		return errcode.DatabaseTrasationErr
	}

	// 更新车主信息
	if err = saveCarOwnerInfo(l.ctx, tx, req.Id, req); err != nil {
		logc.Errorf(l.ctx, "更新车主信息时发生错误, err: %s\n", err.Error())
		if err1 := tx.Rollback(); err1 != nil {
			logc.Errorf(l.ctx, "更新车主信息回滚时发生错误, err: %s\n", err1.Error())
			return errcode.DatabaseRollbackErr
		}
		return errcode.DatabaseUpdateErr
	}

	// update user order.
	query := "UPDATE `user_orders` SET `car_brand_id` = ?, `car_brand_series_id` = ?, `comment` = ?, updated_at = NOW()"
	if _, err = tx.ExecContext(l.ctx, query, req.CarBrandId, req.CarSeriesId, req.Requirements); err != nil {
		logc.Errorf(l.ctx, "更新用户订单时发生错误, err: %s\n", err.Error())
		if err1 := tx.Rollback(); err1 != nil {
			logc.Errorf(l.ctx, "更新用户订单回滚时发生错误, err: %s\n", err.Error())
			return errcode.DatabaseRollbackErr
		}
		return errcode.DatabaseUpdateErr
	}

	// update items of user order.
	// filter out the id set, and then update to user order.
	carReplacementIds := func() (result []uint) {
		for _, r := range req.CarReplacements {
			result = append(result, r.Id)
		}
		return
	}()
	// update or create the replacement items.
	if err = saveOrderItems(tx, req.Id, carReplacementIds); err != nil {
		logc.Errorf(l.ctx, "更新用户订单项目时发生错误, err: %s\n", err.Error())
		if err1 := tx.Rollback(); err1 != nil {
			logc.Errorf(l.ctx, "更新用户订单项目回滚时发生错误, err: %s\n", err1.Error())
			return errcode.DatabaseRollbackErr
		}
		return errcode.DatabaseUpdateErr
	}

	if err = tx.Commit(); err != nil {
		logc.Errorf(l.ctx, "更新用户订单提交时发生错误, err: %s\n", err.Error())
		if err1 := tx.Rollback(); err1 != nil {
			logc.Errorf(l.ctx, "更新用户订单提交回滚时发生错误, err: %s\n", err1.Error())
			return errcode.DatabaseRollbackErr
		}
		return errcode.DatabaseCommitErr
	}

	return nil
}

// saveCarOwnerInfo 更新车主信息
func saveCarOwnerInfo(
	ctx context.Context,
	tx *sqlx.Tx,
	orderId uint,
	reqData *types.UpdateUserOrderReq,
) error {
	// 获取订单用户车主信息
	query := "SELECT `name` AS `name`, `phone_number` AS `phoneNumber`, `multilevel_address` AS `multilevelAddress`, `full_address` AS `fullAddress` FROM `car_owner_infos` WHERE `user_order_id` = ? LIMIT 1;"

	var carOwnerInfo struct {
		Name              string `db:"name"`
		PhoneNumber       string `db:"phoneNumber"`
		MultilevelAddress string `db:"multilevelAddress"`
		FullAddress       string `db:"fullAddress"`
	}

	if err := tx.GetContext(ctx, &carOwnerInfo, query, orderId); err != nil {
		return err
	}

	var updateMap map[string]string = make(map[string]string, 4)
	if reqData.CarOwnerName != carOwnerInfo.Name {
		updateMap["name"] = reqData.CarOwnerName
	}
	if reqData.CarOwnerPhoneNumber != carOwnerInfo.PhoneNumber {
		updateMap["phone_number"] = reqData.CarOwnerPhoneNumber
	}
	if reqData.CarOwnerMultiLvAddr != carOwnerInfo.MultilevelAddress {
		updateMap["multilevel_address"] = reqData.CarOwnerMultiLvAddr
	}
	if reqData.CarOwnerFullAddress != carOwnerInfo.FullAddress {
		updateMap["full_address"] = reqData.CarOwnerFullAddress
	}

	// if there is nothing to update
	if len(updateMap) == 0 {
		return nil
	}

	var i = 0
	var n = len(updateMap)
	fieldSetString := ""
	for k, v := range updateMap {
		fieldSetString += fmt.Sprintf("`%s` = \"%s\"", k, v)
		if i != (n - 1) {
			fieldSetString += ","
		}
		i++
	}

	query = fmt.Sprintf("UPDATE `car_owner_infos` SET %s WHERE `user_order_id` = ? LIMIT 1", fieldSetString)
	if _, err := tx.ExecContext(ctx, query, orderId); err != nil {
		return err
	}

	return nil
}
