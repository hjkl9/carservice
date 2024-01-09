package api

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/trace"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type ApiResponse struct {
	HttpCode int         `json:"httpCode"`
	ErrCode  string      `json:"errCode"`
	Msg      string      `json:"msg"`
	Details  []string    `json:"details,omitempty"`
	Data     interface{} `json:"data"`
	TraceId  string      `json:"traceId,omitempty"`
}

func ResponseWithCtx(ctx context.Context, w http.ResponseWriter, data interface{}, err error) {
	ac := ParseApiCoder(err)
	// It's not OK.
	if ac.Code() != OK.Code() {
		body := ApiResponse{
			HttpCode: int(ac.HttpCode()),
			ErrCode:  ac.Code(),
			Msg:      ac.Message(),
			Data:     nil,
			TraceId:  trace.TraceIDFromContext(ctx),
		}
		httpx.WriteJsonCtx(ctx, w, int(ac.HttpCode()), body)
		return
	}
	// It's OK
	body := ApiResponse{
		HttpCode: int(ac.HttpCode()),
		ErrCode:  ac.Code(),
		Msg:      ac.Message(),
		Data:     data,
	}
	httpx.WriteJsonCtx(ctx, w, int(ac.HttpCode()), body)

	// if typ.Coder().Code() != OK.Code() {
	// 	body := ApiResponse{
	// 		HttpCode: int(typ.Coder().HttpCode()),
	// 		ErrCode:  typ.Coder().Code(),
	// 		Msg:      typ.Coder().Message(),
	// 		Data:     nil,
	// 		TraceId:  trace.TraceIDFromContext(ctx),
	// 	}
	// 	httpx.WriteJsonCtx(ctx, w, int(typ.Coder().HttpCode()), body)
	// 	return
	// }
	// // It's OK
	// body := ApiResponse{
	// 	HttpCode: int(typ.Coder().HttpCode()),
	// 	ErrCode:  typ.Coder().Code(),
	// 	Msg:      typ.Coder().Message(),
	// 	Data:     data,
	// }
	// httpx.WriteJsonCtx(ctx, w, int(typ.Coder().HttpCode()), body)
	// return
}
