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
	Details  []string    `json:"details,omitempty"`
	Data     interface{} `json:"data,omitempty"`
}

// todo: wrapping duplicate code
func makeErrorBody(err *errcode.ErrCode) Body {
	fmt.Printf("%#v\n", err)

	return Body{
		HttpCode: err.Code,
		Msg:      err.Error(),
		Details:  err.Details,
		ErrCode:  "todo.",
	}
}

func makeOkBody(resp *interface{}) Body {
	var body Body
	body.HttpCode = 200
	body.Msg = "ok"
	if resp != nil {
		body.Data = *resp
	}
	body.ErrCode = "todo."
	return body
}

func Response(w http.ResponseWriter, resp interface{}, err error) {
	var body Body
	if err != nil {
		realErr := err.(*errcode.ErrCode)
		body = makeErrorBody(realErr)
	} else {
		body = makeOkBody(&resp)
		// todo: test as below.
		// if resp == nil {
		// 	body = makeOkBody(nil)
		// } else {
		// 	body = makeOkBody(&resp)
		// }
	}
	httpx.WriteJson(w, body.HttpCode, body)
}

func ResponseWithCtx(ctx context.Context, w http.ResponseWriter, err error) {
	var body Body
	if err != nil {
		realErr := err.(*errcode.ErrCode)
		body = makeErrorBody(realErr)
	} else {
		body = makeOkBody(nil)
	}
	body.ErrCode = "todo."
	httpx.WriteJsonCtx(ctx, w, body.HttpCode, body)
}
