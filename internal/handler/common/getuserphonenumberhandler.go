package common

import (
	"net/http"

	"carservice/internal/logic/common"
	"carservice/internal/pkg/common/errcode"
	stdresponse "carservice/internal/pkg/httper/response"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetUserPhoneNumberHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetUserPhoneNumberReq
		if err := httpx.Parse(r, &req); err != nil {
			stdresponse.ResponseWithCtx(r.Context(), w, nil, errcode.New(http.StatusBadRequest, "feature.", err.Error()))
			return
		}

		l := common.NewGetUserPhoneNumberLogic(r.Context(), svcCtx)
		resp, err := l.GetUserPhoneNumber(&req)
		stdresponse.Response(w, resp, err)
	}
}
