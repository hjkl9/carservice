package userorder

import (
	"net/http"

	"carservice/internal/logic/userorder"
	"carservice/internal/pkg/httper/api"
	"carservice/internal/svc"
)

func UserOrderStatusListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := userorder.NewUserOrderStatusListLogic(r.Context(), svcCtx)
		resp, err := l.UserOrderStatusList()
		api.ResponseWithCtx(r.Context(), w, resp, err)
	}
}
