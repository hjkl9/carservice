package carownerinfo

import (
	"context"

	"carservice/internal/pkg/common/errcode"
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

func (l *CreateCarOwnerInfoLogic) CreateCarOwnerInfo(req *types.CreateCarOwnerInfoReq) (resp *types.CreateCarOwnerInfoRep, err error) {
	// todo: have not test below getting.
	// user := l.ctx.Value("user").(jwt.UserPayload)
	// userId := user.UserId
	userId := 1
	query := "INSERT INTO `car_owner_infos`(`name`, `user_id`, `phone_number`, `multilevel_address`, `full_address`) VALUES(?, ?, ?, ?, ?)"
	_, err = l.svcCtx.DBC.ExecContext(l.ctx, query, req.Name, userId, req.PhoneNumber, req.MultilevelAddress, req.FullAddress)
	if err != nil {
		return nil, errcode.InternalServerError.SetMsg("创建数据时发生错误").SetDetails(err.Error())
	}
	return nil, nil
}
