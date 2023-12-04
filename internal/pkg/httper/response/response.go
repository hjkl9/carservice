package http

import (
	"carservice/internal/pkg/common/errcode"
	"context"
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type Body struct {
	HttpCode int         `json:"httpCode"`
	ErrCode  string      `json:"errCode"`
	Msg      string      `json:"msg"`
	Data     interface{} `json:"data,omitempty"`
}

// todo: wrapping duplicate code
func makeBody() Body {
	// Do something here...
	return Body{}
}

func Response(w http.ResponseWriter, resp interface{}, err error) {
	var body Body
	realErr := err.(*errcode.ErrCode)

	if err != nil {
		body.HttpCode = realErr.Code
		body.Msg = err.Error()
	} else {
		body.HttpCode = 0
		body.Msg = "OK"
		body.Data = resp
	}
	body.ErrCode = "todo."
	fmt.Printf("111111111111111%#v\n", realErr)
	httpx.WriteJson(w, realErr.Code, body)
}

func ResponseWithCtx(ctx context.Context, w http.ResponseWriter, err error) {
	var body Body
	realErr := err.(*errcode.ErrCode)
	if err != nil {
		body.HttpCode = realErr.Code
		body.Msg = err.Error()
	} else {
		body.HttpCode = 0
		body.Msg = "OK"
	}
	body.ErrCode = "todo."
	fmt.Printf("111111111111111%#v\n", realErr)
	httpx.WriteJsonCtx(ctx, w, realErr.Code, body)
}
