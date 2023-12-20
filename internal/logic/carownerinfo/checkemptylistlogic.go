package carownerinfo

import (
	"context"

	"carservice/internal/pkg/common/errcode"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckEmptyListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckEmptyListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckEmptyListLogic {
	return &CheckEmptyListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckEmptyListLogic) CheckEmptyList() (resp *types.CheckEmptyListRep, err error) {
	// todo: have not test below getting.
	// user := l.ctx.Value("user").(jwt.UserPayload)
	// userId := user.UserId
	userId := 1
	var count int
	query := "SELECT COUNT(1) AS `count` FROM `car_owner_infos` WHERE `user_id` = ?"
	if err = l.svcCtx.DBC.Get(&count, query, userId); err != nil {
		return nil, errcode.InternalServerError.SetMsg("查询数据时发生错误").SetDetails(err.Error())
	}
	return &types.CheckEmptyListRep{
		Listable: count > 0,
	}, nil
}
