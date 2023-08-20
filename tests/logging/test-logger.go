package logging

import (
	"fmt"
	. "prerequisites-tester/utils"
)

type TestLogger struct {
	TestName string
}

func (logger *TestLogger) PrintStart() {
	logger.printMessage(Cyan, "START", "Starting test...")
}

func (logger *TestLogger) PrintTestFailed(message string) {
	logger.printMessage(Red, "FAILED", message)
	fmt.Println()
}

func (logger *TestLogger) PrintSoftFailure(message string) {
	logger.printMessage(Orange, "STEP FAILED", message)
}

func (logger *TestLogger) PrintStatus(message string) {
	logger.printMessage(Reset, "STATUS", message)
}

func (logger *TestLogger) PrintTestPassed() {
	logger.printMessage(Green, "PASSED", "All steps passed successfully")
	fmt.Println()
}

func (logger *TestLogger) PrintStepSuccess(message string) {
	logger.printMessage(Blue, "STEP SUCCEEDED", message)
}

func (logger *TestLogger) printMessage(color string, status string, message string) {
	PrintMessage(color, logger.TestName, status, message)
}
