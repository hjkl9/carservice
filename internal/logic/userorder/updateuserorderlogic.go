package userorder

import (
	"context"
	"encoding/json"

	"carservice/internal/pkg/common/errcode"
	"carservice/internal/pkg/jwt"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserOrderLogic {
	return &UpdateUserOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserOrderLogic) UpdateUserOrder(req *types.UpdateUserOrderReq) error {
	userId, err := (jwt.GetUserId(l.ctx).(json.Number)).Int64()
	if err != nil {
		return errcode.InternalServerError.Lazy("UserID 类型转换时发生错误").SetDetails(err.Error())
	}

	// validate parameters //
	// check phone number.
	if len(req.CarOwnerPhoneNumber) != 11 {
		return errcode.InvalidPhoneNumberError
	}

	// check if the data already exists. //
	// check CarOwnerInfo
	var counter carOwnerInfoCounter
	query := "SELECT COUNT(1) AS `count`, MIN(`id`) AS `firstId` FROM `car_owner_infos` WHERE `user_id` = ? LIMIT 1"
	stmtx, err := l.svcCtx.DBC.PreparexContext(l.ctx, query)
	if err != nil {
		return err
	}
	if err = stmtx.GetContext(l.ctx, &counter, userId); err != nil {
		return err
	}

	// what something need to change?

	// update data.

	return nil
}
