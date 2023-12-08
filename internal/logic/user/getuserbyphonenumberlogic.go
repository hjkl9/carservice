package user

import (
	"context"
	"fmt"

	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserByPhoneNumberLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserByPhoneNumberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserByPhoneNumberLogic {
	return &GetUserByPhoneNumberLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserByPhoneNumberLogic) GetUserByPhoneNumber(req *types.GetUserByPhoneNumberReq) (resp *types.GetUserByPhoneNumberRep, err error) {
	// todo: add your logic here and delete this line
	fmt.Println("You're certified.")
	return
}
