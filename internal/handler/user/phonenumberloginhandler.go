package user

import (
	"net/http"

	"carservice/internal/logic/user"
	"carservice/internal/pkg/common/errcode"
	stdresponse "carservice/internal/pkg/httper/response"
	"carservice/internal/pkg/sms"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func PhoneNumberLoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PhoneNumberLoginReq
		if err := httpx.Parse(r, &req); err != nil {
			stdresponse.ResponseWithCtx(r.Context(), w, nil, errcode.New(http.StatusBadRequest, "feature.", err.Error()))
			return
		}

		// Customize validation.
		if !sms.CheckPhoneNumber(req.PhoneNumber) {
			stdresponse.ResponseWithCtx(r.Context(), w, nil, errcode.New(http.StatusBadRequest, "-", "无效的手机号码"))
			return
		}
		if len(req.Captcha) != 6 {
			stdresponse.ResponseWithCtx(r.Context(), w, nil, errcode.New(http.StatusBadRequest, "-", "无效的手机验证码"))
			return
		}

		l := user.NewPhoneNumberLoginLogic(r.Context(), svcCtx)
		resp, err := l.PhoneNumberLogin(&req)
		stdresponse.Response(w, resp, err)
	}
}
