package carownerinfo

import (
	"context"

	"carservice/internal/pkg/common/errcode"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCarOwnerInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateCarOwnerInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCarOwnerInfoLogic {
	return &UpdateCarOwnerInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateCarOwnerInfoLogic) UpdateCarOwnerInfo(req *types.UpdateCarOwnerInfoReq) (resp *types.UpdateCarOwnerInfoRep, err error) {
	userId := 1

	var count uint
	query := "SELECT count(1) AS `count` FROM `car_owner_infos` WHERE `id` = ? AND `user_id` = ?"
	l.svcCtx.DBC.Get(&count, query, req.Id, userId)
	if count == 0 {
		return nil, errcode.NotFound.SetMsg("找不到该用户的车主信息")
	}

	query = "UPDATE `car_owner_infos` SET `name`= ?, `phone_number`= ?, `multilevel_address`= ?, `full_address`= ?, `longitude`= ?, `latitude`= ? WHERE `id`= ? AND `user_id`"
	_, err = l.svcCtx.DBC.ExecContext(l.ctx, query, req.Name, req.PhoneNumber, req.MultilevelAddress, req.FullAddress, req.Longitude, req.Latitude, req.Id, userId)
	if err != nil {
		return nil, errcode.InternalServerError.SetMsg("更新数据时发生错误").SetDetails(err.Error())
	}
	return nil, nil
}
