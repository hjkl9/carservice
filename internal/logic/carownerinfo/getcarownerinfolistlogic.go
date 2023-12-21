package carownerinfo

import (
	"context"

	"carservice/internal/pkg/common/errcode"
	"carservice/internal/pkg/jwt"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCarOwnerInfoListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCarOwnerInfoListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCarOwnerInfoListLogic {
	return &GetCarOwnerInfoListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// todo: do testing.
// GetCarOwnerInfoList 获取用户车主信息列表
func (l *GetCarOwnerInfoListLogic) GetCarOwnerInfoList(req *types.GetCarOwnerInfoListReq) (resp []types.CarOwnerInfoListItem, err error) {
	userId := jwt.GetUserId(l.ctx)
	query := "SELECT `id`, `name`, `phoneNumber`, `multilevelAddress`, `fullAddress` FROM `car_owner_infos` WHERE `user_id` = ?"
	var list []types.CarOwnerInfoListItem
	if err = l.svcCtx.DBC.Select(&list, query, userId); err != nil {
		return nil, errcode.InternalServerError.SetMsg("查询数据时出现错误").SetDetails(err.Error())
	}
	return list, nil
}
