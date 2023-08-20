package utils

import "fmt"

const (
	Red         = "\033[31m"
	Green       = "\033[32m"
	Yellow      = "\033[33m"
	Blue        = "\033[34m"
	Reset       = "\033[0m"
	BrightGreen = "\033[32;1m"
	Orange      = "\033[38;2;255;165;0m"
	Cyan        = "\033[36m"
	Magenta     = "\033[35m"
)

func PrintMessage(color string, name string, status string, message string) {
	fmt.Println(color + fmt.Sprintf("[%v] [%v] %v", name, status, message))
}
