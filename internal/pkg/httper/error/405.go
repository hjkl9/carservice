package error

import (
	"carservice/internal/pkg/common/errcode"
	stdresponse "carservice/internal/pkg/httper/response"
	"net/http"
)

func NotAllowedHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stdresponse.Response(w, nil, errcode.New(http.StatusMethodNotAllowed, "-", "错误的请求方式 (Http Method)"))
	}
}
