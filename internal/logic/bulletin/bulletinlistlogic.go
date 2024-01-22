package bulletin

import (
	"context"

	"carservice/internal/pkg/common/errcode"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BulletinListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBulletinListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BulletinListLogic {
	return &BulletinListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BulletinListLogic) BulletinList(req *types.GetBulletinListReq) (resp []types.BulletinListItem, err error) {
	if req.Limit > 20 {
		return []types.BulletinListItem{}, errcode.BulletinLimitTooLarge
	}
	// Mock bulletin list.
	resp = append(resp, types.BulletinListItem{Id: 1, Title: "恭喜用户 xxx 下单成功"})
	resp = append(resp, types.BulletinListItem{Id: 2, Title: "车壹壹小程序处于体验版阶段"})
	resp = append(resp, types.BulletinListItem{Id: 3, Title: "车壹壹小程序正在更新迭代..."})
	return
}
