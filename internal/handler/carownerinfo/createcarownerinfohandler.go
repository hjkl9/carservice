package carownerinfo

import (
	"net/http"

	"carservice/internal/logic/carownerinfo"
	"carservice/internal/pkg/common/errcode"
	stdresponse "carservice/internal/pkg/httper/response"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func CreateCarOwnerInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateCarOwnerInfoReq
		if err := httpx.Parse(r, &req); err != nil {
			stdresponse.ResponseWithCtx(r.Context(), w, errcode.New(http.StatusBadRequest, "feature.", err.Error()))
			return
		}

		l := carownerinfo.NewCreateCarOwnerInfoLogic(r.Context(), svcCtx)
		resp, err := l.CreateCarOwnerInfo(&req)
		stdresponse.Response(w, resp, err)
	}
}
