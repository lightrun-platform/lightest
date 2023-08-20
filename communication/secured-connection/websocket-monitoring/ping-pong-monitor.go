package websocket_monitoring

import (
	"fmt"
	"github.com/gorilla/websocket"
	. "prerequisites-tester/tests/logging"
	"sync"
	"time"
)

const (
	pingInterval      = time.Second * 10
	writeWait         = time.Second * 5
	pongWait          = time.Second * 15
	maxMonitoringTime = time.Second * 60
)

type PingPongMonitor struct {
	Connection *websocket.Conn
	*TestLogger

	errChan           chan error
	isMonitoring      bool
	isMonitoringMutex sync.Mutex
	pingTicker        *time.Ticker
	monitoringTimer   <-chan time.Time
	pongTimer         *time.Timer
	pongTimerChan     chan bool
}

func (monitor *PingPongMonitor) Monitor(conn *websocket.Conn) bool {
	monitor.Connection = conn
	monitor.errChan = make(chan error)
	defer close(monitor.errChan)
	monitor.isMonitoring = true
	defer func() {
		monitor.isMonitoringMutex.Lock()
		monitor.isMonitoring = false
		monitor.isMonitoringMutex.Unlock()
	}()

	monitor.PrintStatus(
		fmt.Sprintf(
			"Starting ping ticker. Will send a ping every %d seconds for %d seconds",
			int(pingInterval.Seconds()), maxMonitoringTime,
		),
	)

	// Setting timers
	monitor.setPongHandler()
	monitor.pingTicker = time.NewTicker(pingInterval)
	defer monitor.pingTicker.Stop()
	monitor.monitoringTimer = time.After(maxMonitoringTime)

	go monitor.sendPing()
	monitor.pongTimer = time.AfterFunc(pongWait, func() { monitor.pongTimerChan <- true })

	return monitor.handleEvents()
}

func (monitor *PingPongMonitor) handleEvents() bool {
	for {
		select {
		case <-monitor.pingTicker.C:
			go monitor.sendPing()
		case <-monitor.monitoringTimer:
			{
				monitor.PrintStepSuccess(
					fmt.Sprintf(
						"Monitored the websocket connection for %d seconds. No disconnections occurred.",
						maxMonitoringTime.Seconds(),
					),
				)
				return true
			}
		case <-monitor.pongTimerChan:
			{
				monitor.PrintTestFailed("Did not receive a pong response. Websocket stability in question.")
				return false
			}
		case <-monitor.errChan:
			return false
		}
	}
}

func (monitor *PingPongMonitor) setPongHandler() {
	monitor.Connection.SetPongHandler(
		func(string) error {
			err := monitor.Connection.SetReadDeadline(time.Now().Add(pongWait))
			if err != nil {
				monitor.PrintTestFailed(
					fmt.Sprintf(
						"Setting websocket connection read deadline failed with error: %v", err,
					),
				)
				monitor.errChan <- err
				return err
			}
			monitor.PrintStepSuccess("Pong was received.")
			if !monitor.pongTimer.Stop() {
				<-monitor.pongTimer.C
			}
			return nil
		},
	)
}

func (monitor *PingPongMonitor) sendPing() {
	if !monitor.isMonitoring {
		return
	}
	monitor.PrintStatus("Sending ping...")
	if err := monitor.Connection.WriteControl(
		websocket.PingMessage, []byte{}, time.Now().Add(writeWait),
	); err != nil {
		monitor.PrintTestFailed(
			fmt.Sprintf(
				"Writing ping message failed with error: %v", err,
			),
		)
		if monitor.isMonitoring {
			monitor.isMonitoringMutex.Lock()
			monitor.errChan <- err
			monitor.isMonitoringMutex.Unlock()
		}
		return
	}
	monitor.PrintStepSuccess("Sent ping successfully.")
	monitor.pongTimer.Reset(pongWait)
}

func (monitor *PingPongMonitor) SetLogger(logger *TestLogger) {
	monitor.TestLogger = logger
}
