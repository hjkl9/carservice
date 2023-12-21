package carownerinfo

import (
	"context"

	"carservice/internal/pkg/common/errcode"
	"carservice/internal/pkg/jwt"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCarOwnerInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCarOwnerInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCarOwnerInfoLogic {
	return &GetCarOwnerInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// todo: testing to do.
// GetCarOwnerInfo 获取用户车主信息
func (l *GetCarOwnerInfoLogic) GetCarOwnerInfo(req *types.GetCarOwnerInfoReq) (resp *types.GetCarOwnerInfoRep, err error) {
	userId := jwt.GetUserId(l.ctx)
	var count int
	query := "SELECT count(1) AS `count` FROM `car_owner_infos` WHERE `user_id` = ? AND `id` = ? LIMIT 1"
	if err = l.svcCtx.DBC.Get(&count, query, userId, req.Id); err != nil {
		return nil, errcode.InternalServerError.SetMsg("获取数据时出现错误").SetDetails(err.Error())
	}
	if count == 0 {
		return nil, errcode.NotFound.SetMsg("找不到该用户的车主信息")
	}
	query = "SELECT `id`, `name`, `phoneNumber`, `multilevelAddress`, `fullAddress` FROM `car_owner_infos` WHERE `user_id` = ? AND `id` = ? LIMIT 1"
	var info types.GetCarOwnerInfoRep
	if err = l.svcCtx.DBC.Get(&info, query, userId, req.Id); err != nil {
		return nil, errcode.InternalServerError.SetMsg("获取数据时出现错误").SetDetails(err.Error())
	}
	return &info, nil
}
