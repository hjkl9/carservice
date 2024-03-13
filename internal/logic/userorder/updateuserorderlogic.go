package userorder

import (
	dt_user_order "carservice/internal/datatypes/userorder"
	"context"
	"database/sql"
	"encoding/json"

	"carservice/internal/enum/payment"
	"carservice/internal/enum/userorder"
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

	// only update order.
	updatePayload := &dt_user_order.UpdatePayload{
		MemberId:         uint(userId),
		CarBrandId:       uint(req.CarBrandId),
		CarBrandSeriesId: uint(req.CarSeriesId),
		// OrderNumber:      order.GenerateNumber(time.Now()), // not allow to update.
		Comment:       req.Requirements,
		EstAmount:     0.000000,
		ActAmount:     0.000000,
		PaymentMethod: uint8(payment.DefaultAtCreation),
		OrderStatus:   uint8(userorder.DefaultAtCreation),
		// CarOwnerInfoId:   *carOwnerInfoId, // ! deprecated
		PartnerStoreId: uint(req.PartnerStoreId), // ! deprecated
	}

	_ = updatePayload
	_ = tx

	return nil
}
