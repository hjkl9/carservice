package sms

import (
	"net/http"

	"github.com/zeromicro/x/errors"

	"carservice/internal/logic/sms"
	smsutil "carservice/internal/pkg/sms"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func SendCaptchaHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SendCaptchaReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// Customize validation.
		if !smsutil.CheckPhoneNumber(req.PhoneNumber) {
			httpx.ErrorCtx(r.Context(), w, errors.New(http.StatusBadRequest, "无效的手机号码"))
			return
		}

		l := sms.NewSendCaptchaLogic(r.Context(), svcCtx)
		resp, err := l.SendCaptcha(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
