package sms

import (
	"context"
	"net/http"

	"carservice/internal/pkg/constant"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"
)

type SendCaptchaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendCaptchaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendCaptchaLogic {
	return &SendCaptchaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendCaptchaLogic) SendCaptcha(req *types.SendCaptchaReq) (resp *types.SendCaptchaRep, err error) {
	// todo: add your logic here and delete this line
	key := req.PhoneNumber + constant.SmsCaptchaPrefix
	cmd := l.svcCtx.RDBC.Exists(l.ctx, key)
	n, err := cmd.Result()
	if err != nil {
		return nil, errors.New(http.StatusInternalServerError, "Redis 数据库查询数据时出现错误")
	}
	// cannot be sent repeatedly.
	if n != 1 {
		return nil, errors.New(http.StatusBadRequest, "不能重复发送验证码")
	}
	// todo: send sms logic.
	// sms := smsutil.NewSms(l.svcCtx.Config)
	// sms.Send()
	resp = &types.SendCaptchaRep{}
	return resp, errors.New(http.StatusOK, "It's ok.")
}
