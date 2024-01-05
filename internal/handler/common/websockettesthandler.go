package common

import (
	"net/http"

	"carservice/internal/pkg/common/errcode"
	stdresponse "carservice/internal/pkg/httper/response"
	"carservice/internal/svc"
	"carservice/internal/types"
	"carservice/internal/ws"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// upgrade to websocket connection.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// global
var wm = ws.NewWebsocketManager()

func WebsocketTestHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user.
		var req types.WebsocketTestReq
		if err := httpx.Parse(r, &req); err != nil {
			stdresponse.ResponseWithCtx(r.Context(), w, errcode.New(http.StatusBadRequest, "feature.", err.Error()))
			return
		}
		// create websocket connection and handle error.
		conn, err := upgrader.Upgrade(w, r, nil)

		remoteAddr := conn.RemoteAddr()
		logc.Infof(r.Context(), "用户 [%s], ip: [%s] 连接成功\n", "1", remoteAddr.String())
		if err != nil {
			stdresponse.Response(w, nil, err)
			return
		}
		defer func() {
			logc.Info(r.Context(), "someone has disconnect")
			err := conn.Close()
			if err != nil {
				stdresponse.Response(w, nil, err)
				return
			}
			// delete connection.
			wm.RemoveConnection(1)
		}()

		// add connection.
		wm.AddConnection(1, conn)

		// l := common.NewWebsocketTestLogic(r.Context(), svcCtx)
		// err = l.WebsocketTest(&req)
		// stdresponse.Response(w, nil, err)
	}
}
