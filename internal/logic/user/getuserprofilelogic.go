package user

import (
	"context"
	"net/http"

	"carservice/internal/pkg/common/errcode"
	"carservice/internal/pkg/jwt"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserProfileLogic {
	return &GetUserProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserProfileLogic) GetUserProfile() (resp *types.GetUserProfileRep, err error) {
	userId := jwt.GetUserId(l.ctx)
	// 检查用户是否被锁定
	query := "SELECT `locked` FROM `members` WHERE id = ? LIMIT 1"
	var locked uint8
	if err = l.svcCtx.DBC.Get(&locked, query, userId); err != nil {
		return nil, errcode.NewDatabaseErrorx().GetError(err)
	}
	if locked == 1 {
		return nil, errcode.New(http.StatusLocked, "-", "该用户被锁定")
	}
	// 获取用户
	var userProfile struct {
		Username    string `db:"username"`
		PhoneNumber string `db:"phoneNumber"`
		AvatarUrl   string `db:"avatarUrl"`
	}
	query = "SELECT `username` AS `username`, `phone_number` AS `phoneNumber`, `avatar_url` AS `avatarUrl` FROM `members` WHERE `id` = ? LIMIT 1"
	if err = l.svcCtx.DBC.Get(&userProfile, query, userId); err != nil {
		return nil, errcode.NewDatabaseErrorx().GetError(err)
	}
	userProfile.AvatarUrl = "NO_AVAILABLE_CLOUD_STORAGE"
	return &types.GetUserProfileRep{
		Username:    userProfile.Username,
		PhoneNumber: userProfile.PhoneNumber,
		AvatarUrl:   userProfile.AvatarUrl,
	}, nil
}
