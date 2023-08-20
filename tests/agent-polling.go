package tests

import (
	"fmt"
	"prerequisites-tester/communication/requests"
	. "prerequisites-tester/communication/response-handling"
	"prerequisites-tester/config"
	"time"
)

const (
	pollingIntervalSeconds int = 10
	pollingDurationSeconds int = 30
	minimumSuccessfulPolls int = 2
	maximumFailedPolls     int = 1
)

type PollingTest struct {
	TestBase

	successfulPolls int
	totalPolls      int
	pollingTicker   *time.Ticker
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
	test.PrintStatus(
		fmt.Sprintf(
			"Long polling the agent for %d seconds, with %d seconds interval...",
			pollingDurationSeconds, pollingIntervalSeconds,
		),
	)

	pollingTimer := time.After(time.Duration(pollingDurationSeconds) * time.Second)
	pollingTicker := time.NewTicker(time.Duration(pollingIntervalSeconds) * time.Second)
	defer pollingTicker.Stop()

	test.handlePollingIteration() // Make one polling iteration before first tick.
	for {
		select {
		case <-pollingTicker.C:
			if !test.handlePollingIteration() {
				return false
			}
		case <-pollingTimer:
			return test.assessPollingStats()
		}
	}
}

func (test *PollingTest) handlePollingIteration() bool {
	test.successfulPolls += test.pollServer()
	test.totalPolls++
	if test.totalPolls-test.successfulPolls > maximumFailedPolls {
		test.PrintTestFailed("Too many unsuccessful polling requests.")
		return false
	}
	return true
}

func (test *PollingTest) assessPollingStats() bool {
	test.PrintStatus("Finished polling segment.")
	if test.successfulPolls >= minimumSuccessfulPolls {
		test.PrintStepSuccess("Long polling finished successfully")
		test.PrintTestPassed()
		return true
	} else {
		test.PrintTestFailed("Not enough successful long polling requests.")
		return false
	}
}

func (test *PollingTest) registerAgent() bool {
	regData := requests.RegistrationData{
		CompanyId:   config.GetConfig().CompanyId,
		ServerUrl:   config.GetConfig().ServerUrl,
		AgentConfig: config.GetConfig().AgentConfig,
	}

	test.PrintStatus("Registering agent...")
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
	test.PrintStatus("Making a long polling request...")
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
