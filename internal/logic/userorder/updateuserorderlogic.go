package userorder

import (
	"context"
	"encoding/json"

	"carservice/internal/pkg/common/errcode"
	"carservice/internal/pkg/jwt"
	"carservice/internal/svc"
	"carservice/internal/types"

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
	_ = userId

	// validate CarBrand and CarBrandSeries data.
	// hasCar, err := l.validateUserCar(req.CarBrandId, req.CarSeriesId)
	// if err != nil {
	// 	return nil, errcode.DatabaseError.Lazy("操作数据库时发生错误", err.Error())
	// }
	// if !hasCar {
	// 	return nil, errcode.NotFound.SetMsg("该车辆不存在")
	// }

	return nil
}
