package user

import (
	"context"
	"fmt"
	"time"

	"carservice/internal/data/tables"
	"carservice/internal/pkg/common/errcode"
	"carservice/internal/pkg/jwt"
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
	mockUserId := 21
	var count int
	query := "SELECT count(1) AS `count` FROM `%s` WHERE `id` = ? LIMIT 1"
	err = l.svcCtx.DBC.Get(&count, fmt.Sprintf(query, tables.User), mockUserId)
	if err != nil {
		return nil, errcode.DatabaseError.SetMsg("查询数据库时发生错误").SetDetails(err.Error())
	}
	token, err := jwt.GetJwtToken(l.svcCtx.Config.JwtConf.AccessSecret, time.Now().Unix(), 36000, uint(mockUserId))
	if err != nil {
		return nil, errcode.DatabaseError.SetMsg("生成 AccessToken 时发生错误").SetDetails(err.Error())
	}
	return &types.MockLoginReq{
		Token: token,
	}, nil
}
