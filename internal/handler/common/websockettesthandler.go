package common

import (
	"fmt"
	"net/http"

	"carservice/internal/logic/common"
	"carservice/internal/pkg/common/errcode"
	stdresponse "carservice/internal/pkg/httper/response"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// upgrade to websocket connection.
var upgrader = websocket.Upgrader{}

func WebsocketTestHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.WebsocketTestReq
		if err := httpx.Parse(r, &req); err != nil {
			stdresponse.ResponseWithCtx(r.Context(), w, errcode.New(http.StatusBadRequest, "feature.", err.Error()))
			return
		}

		// create websocket connection and handle error.
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			stdresponse.Response(w, nil, err)
			return
		}
		defer func() {
			err := conn.Close()
			if err != nil {
				stdresponse.Response(w, nil, err)
				return
			}
		}()

		for {
			// read message and handle
			mt, msg, err := conn.ReadMessage()
			if err != nil {
				logc.Errorf(r.Context(), "failed to read message of the current connection, err: %s\n", err.Error())
				continue
			}
			logc.Infof(r.Context(), "get message: %s\n", msg)
			writeMsg := fmt.Sprintf("server message: %s", string(msg))
			err = conn.WriteMessage(mt, []byte(writeMsg))
			if err != nil {
				logc.Errorf(r.Context(), "failed to write message of the current connection, err: %s\n", err.Error())
				continue
			}
		}

		l := common.NewWebsocketTestLogic(r.Context(), svcCtx)
		err = l.WebsocketTest(&req)
		stdresponse.Response(w, nil, err)
	}
}
