package main

import (
	"fmt"
	tls_configuration "prerequisites-tester/communication/secured-connection/tls-configuration"
	websocket_monitoring "prerequisites-tester/communication/secured-connection/websocket-monitoring"
	"prerequisites-tester/config"
	"prerequisites-tester/tests"
	. "prerequisites-tester/utils"
)

func main() {
	// Loading the tester configuration
	if err := config.LoadConfig("config.json"); err != nil {
		fmt.Println("Error loading configuration:", err)
		return
	}

	// Initializing a slice of Tester with several instantiations
	testers := []tests.Tester{
		&tests.PollingTest{},
		&tests.WebsocketTest{
			TlsConfigProducer: &tls_configuration.PinnedCertConfigProducer{},
			WebsocketMonitor:  &websocket_monitoring.PingPongMonitor{},
		},
		&tests.WebsocketTest{
			TlsConfigProducer: &tls_configuration.RootCertsConfigProducer{},
			WebsocketMonitor:  &websocket_monitoring.PingPongMonitor{},
		},
	}

	testsPassed := 0
	testsFailed := 0

	PrintMessage(Magenta, "Manager", "START", "Starting test suite...")
	fmt.Println()
	for _, test := range testers {
		test.Initialise()
		if test.Test() {
			testsPassed++
		} else {
			testsFailed++
		}
	}
	PrintMessage(Magenta, "Manager", "END", "Finished test suite...")
	PrintMessage(
		Green, "Manager", "STATUS",
		fmt.Sprintf("Tests passed successfully: %d", testsPassed),
	)
	failureColor := Green
	if testsFailed != 0 {
		failureColor = Red
	}
	PrintMessage(
		failureColor, "Manager", "STATUS",
		fmt.Sprintf("Tests failures: %d", testsFailed),
	)
}
