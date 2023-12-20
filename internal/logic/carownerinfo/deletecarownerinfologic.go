package carownerinfo

import (
	"context"

	"carservice/internal/pkg/common/errcode"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCarOwnerInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteCarOwnerInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCarOwnerInfoLogic {
	return &DeleteCarOwnerInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// todo: need to test it.
// DeleteCarOwnerInfo 删除用户车主信息
func (l *DeleteCarOwnerInfoLogic) DeleteCarOwnerInfo(req *types.DeleteCarOwnerInfoReq) error {
	// todo: have not test below getting.
	// user := l.ctx.Value("user").(jwt.UserPayload)
	// userId := user.UserId
	userId := 1
	var count int
	query := "SELECT count(1) AS `count` FROM `car_owner_infos` WHERE `user_id` = ? AND `id` = ? LIMIT 1"
	if err := l.svcCtx.DBC.Get(&count, query, userId, req.Id); err != nil {
		return errcode.InternalServerError.SetMsg("获取数据时出现错误").SetDetails(err.Error())
	}
	if count == 0 {
		return errcode.NotFound.SetMsg("找不到该用户的车主信息")
	}
	query = "DELETE FROM `car_owner_infos` WHERE `user_id` = ? AND `id` = ?"
	_, err := l.svcCtx.DBC.ExecContext(l.ctx, query, userId, req.Id)
	if err != nil {
		return errcode.InternalServerError.SetMsg("删除数据时出现错误").SetDetails(err.Error())
	}
	return nil
}
