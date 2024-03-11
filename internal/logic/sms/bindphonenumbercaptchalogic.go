package sms

import (
	"context"
	"net/http"
	"strconv"

	"carservice/internal/pkg/common/errcode"
	"carservice/internal/pkg/constant"
	"carservice/internal/pkg/generator/captcha"
	smsutil "carservice/internal/pkg/sms"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BindPhoneNumberCaptchaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBindPhoneNumberCaptchaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindPhoneNumberCaptchaLogic {
	return &BindPhoneNumberCaptchaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BindPhoneNumberCaptchaLogic) BindPhoneNumberCaptcha(req *types.SendCaptchaReq) error {
	key := constant.SmsBindPhoneNumberCaptchaPrefix + req.PhoneNumber
	cmd := l.svcCtx.RDBC.Exists(l.ctx, key)
	n, err := cmd.Result()
	if err != nil {
		return errcode.New(http.StatusInternalServerError, "feature.", "Redis 数据库查询数据时出现错误")
	}
	// cannot be sent repeatedly.
	if n == 1 {
		return errcode.New(http.StatusOK, "feature.", "不能重复发送验证码")
	}
	// send sms logic.
	smsutil := smsutil.NewSms(l.svcCtx.Config)
	// ! should specify the template id.
	templateIdSet := []string{"1713784"}
	// generate random CAPTCHA of 6-bit.
	randomCaptcha := captcha.PhoneNumberCaptcha()
	templateSet := []string{randomCaptcha, strconv.Itoa(captchaExpire)}
	phoneNumberSet := []string{req.PhoneNumber}
	err = smsutil.Send(templateIdSet, templateSet, phoneNumberSet)
	if err != nil {
		return errcode.InternalServerError.Lazy("发送短信时发生错误", err.Error())
	}
	return nil
}
