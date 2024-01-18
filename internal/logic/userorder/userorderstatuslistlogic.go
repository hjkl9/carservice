package userorder

import (
	"context"

	"carservice/internal/enum/userorder"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserOrderStatusListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserOrderStatusListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserOrderStatusListLogic {
	return &UserOrderStatusListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserOrderStatusListLogic) UserOrderStatusList() (resp []types.UserOrderStatusListItem, err error) {
	n, ctl := userorder.ClientTabList()
	for i := 0; i < int(n); i++ {
		resp = append(resp, types.UserOrderStatusListItem{
			Status: ctl[i].Status,
			Label:  ctl[i].Label,
			GoTag:  ctl[i].GoTag,
		})
	}
	return resp, nil
}
