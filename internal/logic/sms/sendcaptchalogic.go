package sms

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"carservice/internal/pkg/common/errcode"
	"carservice/internal/pkg/constant"
	smsutil "carservice/internal/pkg/sms"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	captchaExpire = 10
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

// todo: using RUST to make RPC server for SendCaptcha.
func (l *SendCaptchaLogic) SendCaptcha(req *types.SendCaptchaReq) (resp *types.SendCaptchaRep, err error) {
	key := constant.SmsCaptchaPrefix + req.PhoneNumber
	cmd := l.svcCtx.RDBC.Exists(l.ctx, key)
	n, err := cmd.Result()
	if err != nil {
		return nil, errcode.New(http.StatusInternalServerError, "feature.", "Redis 数据库查询数据时出现错误")
	}
	// cannot be sent repeatedly.
	if n == 1 {
		return nil, errcode.New(http.StatusInternalServerError, "feature.", "不能重复发送验证码")
	}
	// send sms logic.
	smsutil := smsutil.NewSms(l.svcCtx.Config)
	templateIdSet := []string{"1713784"}
	templateSet := []string{"123456", strconv.Itoa(captchaExpire)}
	phoneNumberSet := []string{req.PhoneNumber}
	err = smsutil.Send(templateIdSet, templateSet, phoneNumberSet)
	if err != nil {
		return nil, errcode.New(
			http.StatusInternalServerError,
			"feature.",
			"发送短信时出现错误",
		).SetDetails(err.Error())
	}
	// set CAPTCHA in the rdb when sending was successfully.
	setCmd := l.svcCtx.RDBC.Set(l.ctx, key, req.PhoneNumber, captchaExpire*time.Minute)
	if setCmd.Err() != nil {
		return nil, errcode.New(http.StatusInternalServerError, "feature.", "Redis 数据库创建数据时出现错误")
	}

	resp = &types.SendCaptchaRep{}
	return resp, nil
}
