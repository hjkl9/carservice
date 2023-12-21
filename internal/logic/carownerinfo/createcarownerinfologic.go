package carownerinfo

import (
	"context"

	"carservice/internal/pkg/common/errcode"
	"carservice/internal/pkg/jwt"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCarOwnerInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateCarOwnerInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCarOwnerInfoLogic {
	return &CreateCarOwnerInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// todo: get user id
// CreateCarOwnerInfo 创建用户车主信息
func (l *CreateCarOwnerInfoLogic) CreateCarOwnerInfo(req *types.CreateCarOwnerInfoReq) (resp *types.CreateCarOwnerInfoRep, err error) {
	userId := jwt.GetUserId(l.ctx)
	query := "INSERT INTO `car_owner_infos`(`name`, `user_id`, `phone_number`, `multilevel_address`, `full_address`, `longitude`, `latitude`) VALUES(?, ?, ?, ?, ?, ?, ?)"
	_, err = l.svcCtx.DBC.ExecContext(l.ctx, query, req.Name, userId, req.PhoneNumber, req.MultilevelAddress, req.FullAddress, req.Longitude, req.Latitude)
	if err != nil {
		return nil, errcode.InternalServerError.SetMsg("创建数据时发生错误").SetDetails(err.Error())
	}
	return nil, nil
}
