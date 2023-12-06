package http

import (
	"carservice/internal/pkg/common/errcode"
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type Body struct {
	HttpCode int         `json:"httpCode"`
	ErrCode  string      `json:"errCode"`
	Msg      string      `json:"msg"`
	Details  []string    `json:"details, omitempty"`
	Data     interface{} `json:"data,omitempty"`
}

// todo: wrapping duplicate code
func makeBody() Body {
	// Do something here...
	return Body{}
}

func Response(w http.ResponseWriter, resp interface{}, err error) {
	var body Body
	httpCode := http.StatusOK
	var realErr *errcode.ErrCode
	if err != nil {
		realErr = err.(*errcode.ErrCode)
		httpCode = realErr.Code
	}

	if err != nil {
		body.HttpCode = realErr.Code
		body.Details = realErr.Details
		body.Msg = err.Error()
	} else {
		body.HttpCode = httpCode
		body.Msg = "OK"
		body.Data = resp
	}
	body.ErrCode = "todo."
	httpx.WriteJson(w, httpCode, body)
}

func ResponseWithCtx(ctx context.Context, w http.ResponseWriter, err error) {
	var body Body
	realErr := err.(*errcode.ErrCode)
	if err != nil {
		body.HttpCode = realErr.Code
		body.Details = realErr.Details
		body.Msg = err.Error()
	} else {
		body.HttpCode = realErr.Code
		body.Msg = "OK"
	}
	body.ErrCode = "todo."
	httpx.WriteJsonCtx(ctx, w, realErr.Code, body)
}
