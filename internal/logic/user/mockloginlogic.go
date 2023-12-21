package user

import (
	"context"

	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MockLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMockLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MockLoginLogic {
	return &MockLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MockLoginLogic) MockLogin() (resp *types.MockLoginReq, err error) {
	// todo: add your logic here and delete this line

	return
}
