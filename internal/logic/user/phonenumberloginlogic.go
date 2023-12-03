package user

import (
	"context"

	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PhoneNumberLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPhoneNumberLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PhoneNumberLoginLogic {
	return &PhoneNumberLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PhoneNumberLoginLogic) PhoneNumberLogin(req *types.PhoneNumberLoginReq) (resp *types.PhoneNumberLoginRep, err error) {
	// todo: add your logic here and delete this line
	return
}
