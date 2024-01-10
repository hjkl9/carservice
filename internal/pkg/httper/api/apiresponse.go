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

func Response(w http.ResponseWriter, data interface{}, err error) {
	// It's not OK.
	if err != nil {
		ac := err.(*apiCode)
		body := ApiResponse{
			HttpCode: int(ac.HttpCode()),
			ErrCode:  ac.Code(),
			Msg:      ac.Message(),
			Data:     nil,
		}
		httpx.WriteJson(w, int(ac.HttpCode()), body)
		return
	}
	// It's OK
	body := ApiResponse{
		HttpCode: int(OK.HttpCode()),
		ErrCode:  OK.Code(),
		Msg:      OK.Message(),
		Data:     data,
	}
	httpx.WriteJson(w, int(OK.HttpCode()), body)
}

func ResponseWithCtx(ctx context.Context, w http.ResponseWriter, data interface{}, err error) {
	// It's not OK.
	if err != nil {
		ac := err.(*apiCode)
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
		HttpCode: int(OK.HttpCode()),
		ErrCode:  OK.Code(),
		Msg:      OK.Message(),
		Data:     data,
	}
	httpx.WriteJsonCtx(ctx, w, int(OK.HttpCode()), body)
}

func ResponseListWithCtx(ctx context.Context, w http.ResponseWriter, data interface{}, err error) {
	// todo
}
