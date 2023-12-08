package user

import (
	"context"

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
	return &types.GetUserByPhoneNumberRep{
		Username:  "You're certified.",
		AvatarUrl: "unknown",
	}, nil
}
