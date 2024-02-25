package carbrandseries

import (
	"net/http"

	"carservice/internal/logic/carbrandseries"
	"carservice/internal/pkg/common/errcode"
	"carservice/internal/pkg/httper/api"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func OptionListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetCarBrandSeriesOptionListReq
		if err := httpx.Parse(r, &req); err != nil {
			api.ResponseWithCtx(r.Context(), w, nil, errcode.New(http.StatusBadRequest, "feature.", err.Error()))
			return
		}
		l := carbrandseries.NewOptionListLogic(r.Context(), svcCtx)
		resp, err := l.OptionList(&req)
		api.Response(w, resp, err)
	}
}
