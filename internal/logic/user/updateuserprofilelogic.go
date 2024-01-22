package user

import (
	"context"

	"carservice/internal/pkg/common/errcode"
	"carservice/internal/pkg/jwt"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserProfileLogic {
	return &UpdateUserProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserProfileLogic) UpdateUserProfile(req *types.UpdateUserProfileReq) error {
	req.AvatarUrl = ""
	if len(req.Username) == 0 {
		return errcode.InvalidParametersErr.SetMessage("用户名不能为空")
	}
	userId := jwt.GetUserId(l.ctx)
	query := "UPDATE `members` SET `username` = ?, `avatar_url` = ? WHERE `id` = ?"
	rs, err := l.svcCtx.DBC.ExecContext(l.ctx, query, req.Username, req.AvatarUrl, userId)
	if err != nil {
		return errcode.DatabaseExecuteErr
	}
	n, err := rs.RowsAffected()
	if err != nil {
		return errcode.DatabaseUpdateErr
	}
	if n != 1 {
		return errcode.DatabaseUpdateErr
	}
	return nil
}
