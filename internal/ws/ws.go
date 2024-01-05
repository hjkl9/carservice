package ws

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logc"
)

type UserID uint

// global
type WebsocketManager struct {
	ctx         context.Context
	connections map[UserID]*websocket.Conn
	mutex       sync.Mutex
	upgrader    websocket.Upgrader
}

func NewWebsocketManager() *WebsocketManager {
	wm := &WebsocketManager{
		ctx:         context.Background(),
		connections: make(map[UserID]*websocket.Conn),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
	logc.Info(wm.ctx, "initial websocket manager")
	return wm
}

func (wm *WebsocketManager) GetConnectionsCount() uint {
	return uint(len(wm.connections))
}

func (wm *WebsocketManager) AddConnection(id UserID, conn *websocket.Conn) {
	wm.mutex.Lock()
	defer wm.mutex.Unlock()
	wm.connections[id] = conn
	logc.Infof(wm.ctx, "Some one has connection, connection count is %d\n", wm.GetConnectionsCount())
}

func (wm *WebsocketManager) RemoveConnection(id UserID) {
	wm.mutex.Lock()
	defer wm.mutex.Unlock()
	delete(wm.connections, id)
}

func (wm *WebsocketManager) SendMsgTo(id UserID, msg []byte) error {
	wm.mutex.Lock()
	defer wm.mutex.Unlock()
	fmt.Printf("%#v\n", wm.connections)
	conn, ok := wm.connections[id]
	if !ok {
		return ErrSendMsg
	}
	if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
		return err
	}
	return nil
}
