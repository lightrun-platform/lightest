package websocket_monitoring

import (
	"fmt"
	"github.com/gorilla/websocket"
	. "prerequisites-tester/tests/logging"
	"sync"
	"time"
)

const (
	pingIntervalSeconds  = 10
	pingIntervalDuration = time.Second * pingIntervalSeconds

	readWriteWaitSeconds  = 5
	readWriteWaitDuration = time.Second * readWriteWaitSeconds

	pongWaitSeconds  = 15
	pongWaitDuration = time.Second * pongWaitSeconds

	maxMonitoringSeconds  = 30
	maxMonitoringDuration = time.Second * maxMonitoringSeconds
)

type PingPongMonitor struct {
	Connection *websocket.Conn
	*TestLogger

	isMonitoring             bool
	monitoringMutex          sync.Mutex
	stopPingingTimer         <-chan time.Time
	stopWaitingForPongsTimer <-chan time.Time
	errChan                  chan error
	pingTicker               *time.Ticker
	pongTicker               *time.Ticker
	unansweredPings          int
}

func (monitor *PingPongMonitor) Monitor(conn *websocket.Conn) bool {
	monitor.Connection = conn

	// Setting up the error channel used for signaling the event handler that an error in another goroutine occurred.
	monitor.errChan = make(chan error)
	defer close(monitor.errChan)

	// Setting up the isMonitoring flag that is used to signal other goroutine they should stop. Must be synchronized.
	monitor.isMonitoring = true
	defer func() {
		monitor.monitoringMutex.Lock()
		monitor.isMonitoring = false
		monitor.monitoringMutex.Unlock()
	}()

	// A counter used to track the amount of pings that weren't yet answered with pongs. In a healthy scenario, should
	// be either 0 or 1.
	monitor.unansweredPings = 0

	// Reads messages from the websocket connection in a loop as long as the monitor is active.
	go monitor.readMessages()

	monitor.initialiseTimedEvents()
	defer monitor.pingTicker.Stop()
	defer monitor.pongTicker.Stop()

	// Sending the first ping. Next one will be sent when the ping-ticker will tick.
	go monitor.sendSinglePing()

	// Handling ticker, timer, and channel-based events in a loop.
	return monitor.handleEvents()
}

func (monitor *PingPongMonitor) handleEvents() bool {
	for {
		select {
		case <-monitor.pingTicker.C:
			go monitor.sendSinglePing()
		case <-monitor.pongTicker.C:
			{
				monitor.PrintTestFailed(
					"Did not receive a pong response in too long. Websocket stability in question.",
				)
				return false
			}
		case <-monitor.stopPingingTimer:
			{
				if monitor.handlePingingStopEvent() {
					return true
				}
			}
		case <-monitor.stopWaitingForPongsTimer:
			return monitor.handlePongWaitingStopEvent()
		case <-monitor.errChan:
			return false // Error should have been printed already by method reporting the error.
		}
	}
}

func (monitor *PingPongMonitor) initialiseTimedEvents() {
	monitor.setPongHandler()

	monitor.PrintInfo(
		fmt.Sprintf(
			"Starting ping ticker. Will send a ping every %d seconds for %d seconds",
			pingIntervalSeconds, maxMonitoringSeconds,
		),
	)
	monitor.pingTicker = time.NewTicker(pingIntervalDuration)

	monitor.PrintInfo(
		fmt.Sprintf(
			"Starting pong ticker. Will expect a pong every %d seconds for %d seconds",
			pongWaitSeconds, maxMonitoringSeconds+pongWaitSeconds,
		),
	)
	monitor.pongTicker = time.NewTicker(pongWaitDuration)

	monitor.stopPingingTimer = time.After(maxMonitoringDuration)
	monitor.stopWaitingForPongsTimer = time.After(maxMonitoringDuration + pongWaitDuration)
}

func (monitor *PingPongMonitor) handlePingingStopEvent() bool {
	monitor.PrintInfo("Stopping the pinging loop.")
	monitor.pingTicker.Stop()
	if monitor.unansweredPings == 0 {
		monitor.printTestSuccess(int(maxMonitoringDuration.Seconds()))
		return true
	}
	return false
}

func (monitor *PingPongMonitor) handlePongWaitingStopEvent() bool {
	monitor.PrintInfo("Stopping waiting for pongs.")
	if monitor.unansweredPings == 0 {
		monitor.printTestSuccess(int(maxMonitoringDuration.Seconds()) + int(pongWaitDuration.Seconds()))
		return true
	}
	return false
}

func (monitor *PingPongMonitor) printTestSuccess(monitoringTime int) {
	monitor.PrintStepSuccess(
		fmt.Sprintf(
			"Monitored the websocket connection for %d seconds. All pongs recieved. No disconnections occurred.",
			monitoringTime,
		),
	)
}

func (monitor *PingPongMonitor) setPongHandler() {
	monitor.Connection.SetPongHandler(
		func(string) error {
			monitor.monitoringMutex.Lock()
			err := monitor.Connection.SetReadDeadline(time.Now().Add(pongWaitDuration))
			if err != nil {
				monitor.PrintTestFailed(
					fmt.Sprintf(
						"Setting websocket connection read deadline failed with error: %v", err,
					),
				)
				monitor.errChan <- err
				monitor.monitoringMutex.Unlock()
				return err
			}
			monitor.PrintStepSuccess("Pong was received.")
			monitor.unansweredPings--
			monitor.pongTicker.Reset(pongWaitDuration)
			monitor.monitoringMutex.Unlock()
			return nil
		},
	)
}

func (monitor *PingPongMonitor) readMessages() {
	for monitor.isMonitoring {
		_ = monitor.Connection.SetReadDeadline(time.Now().Add(readWriteWaitDuration))
		messageType, message, _ := monitor.Connection.ReadMessage()
		if messageType == websocket.TextMessage {
			monitor.PrintInfo(fmt.Sprintf("Got websocket message from server: %v", string(message)))
		}
	}
}

func (monitor *PingPongMonitor) sendSinglePing() {
	monitor.monitoringMutex.Lock()
	if !monitor.isMonitoring {
		monitor.monitoringMutex.Unlock()
		return
	}
	monitor.PrintInfo("Sending ping...")
	if err := monitor.Connection.WriteControl(
		websocket.PingMessage, []byte{}, time.Now().Add(readWriteWaitDuration),
	); err != nil {
		monitor.PrintTestFailed(
			fmt.Sprintf(
				"Writing ping message failed with error: %v", err,
			),
		)
		if monitor.isMonitoring {
			monitor.errChan <- err
		}
		monitor.monitoringMutex.Unlock()
		return
	}
	monitor.PrintStepSuccess("Sent ping successfully.")
	monitor.unansweredPings++
	monitor.monitoringMutex.Unlock()
}

func (monitor *PingPongMonitor) SetLogger(logger *TestLogger) {
	monitor.TestLogger = logger
}
