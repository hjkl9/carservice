package userorder

import (
	"context"

	"carservice/internal/svc"
	"carservice/internal/types"

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

func (l *CreateUserOrderLogic) CreateUserOrder(req *types.CreateUserOrderReq) error {
	// todo: add your logic here and delete this line

	return nil
}
