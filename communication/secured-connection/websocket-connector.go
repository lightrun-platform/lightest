package secured_connection

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"prerequisites-tester/communication/requests"
	tls_configuration "prerequisites-tester/communication/secured-connection/tls-configuration"
	"prerequisites-tester/config"
)

type WebsocketConnector struct {
	TlsConfigProducer tls_configuration.TlsConfigProducer
}

func (conn *WebsocketConnector) CreateWebsocketConnection(accessToken string) (*websocket.Conn, error) {
	dialer := websocket.DefaultDialer
	dialer.TLSClientConfig = conn.TlsConfigProducer.CreateConfig()

	wsURL := requests.ProduceWebsocketConnectionURL(
		config.GetConfig().ServerUrl, accessToken, config.GetConfig().CompanyId,
	)
	c, _, err := dialer.Dial(wsURL, http.Header{})
	if err != nil {
		return nil, WebsocketError{Message: fmt.Sprintf("Failed to establish WSS connection: %v", err)}
	}

	return c, nil
}

type WebsocketError struct {
	Message string
}

func (e WebsocketError) Error() string {
	return e.Message
}
