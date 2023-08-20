package websocket_monitoring

import (
	"github.com/gorilla/websocket"
	"prerequisites-tester/tests/logging"
)

type WebsocketMonitor interface {
	Monitor(*websocket.Conn) bool
	SetLogger(logger *logging.TestLogger)
}
