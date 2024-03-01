package carreplacement

import (
	"net/http"

	"carservice/internal/logic/carreplacement"
	"carservice/internal/pkg/common/errcode"
	"carservice/internal/pkg/httper/api"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ReplacementListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CarReplacementReq
		if err := httpx.Parse(r, &req); err != nil {
			api.ResponseWithCtx(r.Context(), w, nil, errcode.InvalidParametersErr)
			return
		}
		l := carreplacement.NewReplacementListLogic(r.Context(), svcCtx)
		resp, err := l.ReplacementList(&req)
		api.ResponseWithCtx(r.Context(), w, resp, err)
	}
}
