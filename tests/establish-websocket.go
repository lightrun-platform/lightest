package tests

import (
	"fmt"
	"prerequisites-tester/communication/requests"
	. "prerequisites-tester/communication/response-handling"
	secured_connection "prerequisites-tester/communication/secured-connection"
	tls_configuration "prerequisites-tester/communication/secured-connection/tls-configuration"
	websocket_monitoring "prerequisites-tester/communication/secured-connection/websocket-monitoring"
	"prerequisites-tester/config"
)

type WebsocketTest struct {
	TestBase
	TlsConfigProducer tls_configuration.TlsConfigProducer
	WebsocketMonitor  websocket_monitoring.WebsocketMonitor
}

func (test *WebsocketTest) Test() bool {
	test.PrintStart()
	if !test.authenticate() {
		return false
	}

	websocketConnector := secured_connection.WebsocketConnector{TlsConfigProducer: test.TlsConfigProducer}

	test.PrintStatus("Creating websocket connection...")
	websocketConnection, err := websocketConnector.CreateWebsocketConnection(
		test.GetContext().Store["accessToken"].(string),
	)
	if err != nil {
		test.PrintTestFailed(fmt.Sprintf("Websocket connection attemp failed with error: %v", err))
		return false
	}
	defer websocketConnection.Close()

	if !test.WebsocketMonitor.Monitor(websocketConnection) {
		return false
	}

	test.PrintTestPassed()
	return true
}

func (test *WebsocketTest) authenticate() bool {
	test.PrintStatus("Authenticating with the Lightrun server...")
	accessToken := new(string) // Will be populated by the registration handler.
	err := requests.Authenticate(
		config.GetConfig().UserEmail, config.GetConfig().UserPassword,
		DefaultAuthenticationHandler(accessToken),
	)
	if err != nil {
		test.PrintTestFailed(fmt.Sprintf("Authentication failed with error: %v", err))
		return false
	}
	test.GetContext().Store["accessToken"] = *accessToken // Stores the populated string in the test context.
	test.PrintStepSuccess("Authentication succeeded. Access token saved.")
	return true
}

func (test *WebsocketTest) Name() string {
	return "Websocket, " + test.TlsConfigProducer.Name()
}

func (test *WebsocketTest) Initialise() {
	test.TestBase.Initialise(test.Name())
	test.WebsocketMonitor.SetLogger(&test.TestBase.TestLogger)
}
