package userorder

import (
	"context"
	"database/sql"
	"encoding/json"

	"carservice/internal/pkg/common/errcode"
	"carservice/internal/pkg/jwt"
	"carservice/internal/svc"
	"carservice/internal/types"

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

	// todo: 更新车主信息

	// only update order.
	// updatePayload := &dt_user_order.UpdatePayload{
	// 	// MemberId:         uint(userId), // ! should'n update UserId.
	// 	CarBrandId:       req.CarBrandId,
	// 	CarBrandSeriesId: req.CarSeriesId,
	// 	// OrderNumber:      order.GenerateNumber(time.Now()), // not allow to update.
	// 	Comment: req.Requirements, // ! should deprecated it.
	// 	// EstAmount:     0.000000,         // ! should deprecated it.
	// 	// ActAmount:     0.000000,         // ! should deprecated it.
	// 	// PaymentMethod: uint8(payment.DefaultAtCreation),
	// 	// OrderStatus:   uint8(userorder.DefaultAtCreation),
	// 	// CarOwnerInfoId:   *carOwnerInfoId, // ! deprecated
	// 	// PartnerStoreId: uint(req.PartnerStoreId), // ! deprecated
	// }

	// update user order.
	query := "UPDATE `user_order` SET `car_brand_id` = ?, `car_brand_series_id` = ?, `comment` = ?, updated_at = NOW()"
	if _, err = tx.ExecContext(l.ctx, query, req.CarBrandId, req.CarSeriesId, req.Requirements); err != nil {
		if err1 := tx.Rollback(); err1 != nil {
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
		return errcode.DatabaseUpdateErr
	}

	return nil
}
