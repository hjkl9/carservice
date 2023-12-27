package user

import (
	"net/http"

	"carservice/internal/logic/user"
	stdresponse "carservice/internal/pkg/httper/response"
	"carservice/internal/svc"
)

func GetUserProfileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewGetUserProfileLogic(r.Context(), svcCtx)
		resp, err := l.GetUserProfile()
		stdresponse.Response(w, resp, err)
	}
}
