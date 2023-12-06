package upload

import (
	"net/http"

	"carservice/internal/logic/common/upload"
	"carservice/internal/pkg/common/errcode"
	stdresponse "carservice/internal/pkg/httper/response"
	"carservice/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func UploadFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			if err == http.ErrMissingFile {
				stdresponse.ResponseWithCtx(r.Context(), w, errcode.New(http.StatusBadRequest, "feature.", err.Error()))
				return
			}
			stdresponse.ResponseWithCtx(r.Context(), w, errcode.New(http.StatusAccepted, "feature.", err.Error()))
			return
		}
		defer file.Close()

		l := upload.NewUploadFileLogic(r.Context(), svcCtx)
		resp, err := l.UploadFile(file, fileHeader)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
