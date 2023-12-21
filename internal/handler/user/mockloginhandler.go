package user

import (
	"net/http"

	"carservice/internal/logic/user"
	stdresponse "carservice/internal/pkg/httper/response"
	"carservice/internal/svc"
)

func MockLoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewMockLoginLogic(r.Context(), svcCtx)
		resp, err := l.MockLogin()
		stdresponse.Response(w, resp, err)
	}
}
