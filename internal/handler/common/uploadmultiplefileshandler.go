package common

import (
	"net/http"

	"carservice/internal/logic/common"
	"carservice/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func UploadMultipleFilesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := common.NewUploadMultipleFilesLogic(r.Context(), svcCtx)
		resp, err := l.UploadMultipleFiles()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
