package user

import (
	"net/http"

	"github.com/zeromicro/x/errors"

	"carservice/internal/logic/user"
	"carservice/internal/pkg/sms"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func PhoneNumberLoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PhoneNumberLoginReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// Customize validation.
		if !sms.CheckPhoneNumber(req.PhoneNumber) {
			httpx.ErrorCtx(r.Context(), w, errors.New(http.StatusBadRequest, "无效的手机号码"))
			return
		}
		if len(req.Captcha) != 6 {
			httpx.ErrorCtx(r.Context(), w, errors.New(http.StatusBadRequest, "无效的验证码"))
			return
		}

		l := user.NewPhoneNumberLoginLogic(r.Context(), svcCtx)
		resp, err := l.PhoneNumberLogin(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
