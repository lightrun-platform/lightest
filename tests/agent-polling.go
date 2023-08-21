package tests

import (
	"fmt"
	"prerequisites-tester/communication/requests"
	. "prerequisites-tester/communication/response-handling"
	"prerequisites-tester/config"
	"sync"
	"time"
)

const (
	pollingIntervalSeconds int = 10
	pollingDurationSeconds int = 20
	minimumSuccessfulPolls int = 2
	maximumFailedPolls     int = 1
)

type PollingTest struct {
	TestBase

	mutex           sync.Mutex
	shouldStop      bool
	successfulPolls int
	totalPolls      int
	pollingTimer    *time.Timer
}

func (test *PollingTest) Test() bool {
	test.PrintStart()
	if !test.registerAgent() {
		return false
	}

	return test.performPollingSegment()
}

func (test *PollingTest) performPollingSegment() bool {
	test.PrintInfo(
		fmt.Sprintf(
			"Long polling the agent for %d seconds, with %d seconds interval...",
			pollingDurationSeconds, pollingIntervalSeconds,
		),
	)

	test.shouldStop = false
	endPollingTimer := time.After(time.Duration(pollingDurationSeconds) * time.Second)

	go test.performPolling() // Make one polling iteration before first tick.
	<-endPollingTimer
	test.shouldStop = true
	return test.assessPollingStats()
}

func (test *PollingTest) performPolling() {
	for {
		test.mutex.Lock()
		if test.shouldStop {
			test.mutex.Unlock()
			return
		}
		test.successfulPolls += test.pollServer()
		test.totalPolls++
		if test.totalPolls-test.successfulPolls > maximumFailedPolls {
			test.PrintTestFailed("Too many unsuccessful polling requests.")
			test.mutex.Unlock()
			return
		}
		test.mutex.Unlock()
	}
}

func (test *PollingTest) assessPollingStats() bool {
	test.mutex.Lock()
	test.PrintInfo("Finished polling segment.")
	if test.successfulPolls >= minimumSuccessfulPolls {
		test.PrintStepSuccess("Long polling finished successfully")
		test.PrintTestPassed()
		test.mutex.Unlock()
		return true
	} else {
		test.PrintTestFailed("Not enough successful long polling requests.")
		test.mutex.Unlock()
		return false
	}
}

func (test *PollingTest) registerAgent() bool {
	regData := requests.RegistrationData{
		CompanyId:   config.GetConfig().CompanyId,
		ServerUrl:   config.GetConfig().ServerUrl,
		AgentConfig: config.GetConfig().AgentConfig,
	}

	test.PrintInfo("Registering agent...")
	agentId := new(string) // Will be populated by the registration handler.
	err := requests.RegisterAgent(regData, DefaultRegistrationHandler(agentId))
	if err != nil {
		test.PrintTestFailed(fmt.Sprintf("Registering the agent failed with error: %v", err))
		return false
	}
	test.GetContext().Store["agentId"] = *agentId // Stores the populated string in the test context.
	test.PrintStepSuccess("Agent registered.")
	return true
}

func (test *PollingTest) pollServer() int {
	test.PrintInfo("Making a long polling request...")
	err := requests.GetBreakpoints(test.GetContext().Store["agentId"].(string), &StatusHandler{&HandlingContext{}})
	if err != nil {
		test.PrintSoftFailure(fmt.Sprintf("Request failed with error: %v", err))
		return 0
	}
	test.PrintStepSuccess("Request succeeded. Response parsed successfully.")
	return 1
}

func (test *PollingTest) Name() string {
	return "Agent Long Polling"
}

func (test *PollingTest) Initialise() {
	test.TestBase.Initialise(test.Name())
}
