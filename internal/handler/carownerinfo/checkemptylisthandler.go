package carownerinfo

import (
	"net/http"

	"carservice/internal/logic/carownerinfo"
	stdresponse "carservice/internal/pkg/httper/response"
	"carservice/internal/svc"
)

func CheckEmptyListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := carownerinfo.NewCheckEmptyListLogic(r.Context(), svcCtx)
		resp, err := l.CheckEmptyList()
		stdresponse.Response(w, resp, err)
	}
}
