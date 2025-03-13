package logger

import (
	"fmt"
	"log"
)

type logger struct {
	*log.Logger
}

// Info 下面的方法完全为了对接nebula_go的logger
func (l logger) Info(msg string) {
	Logger.Output(3, fmt.Sprintf("[%s] %s", "info", msg))
}

func (l logger) Warn(msg string) {
	Logger.Output(3, fmt.Sprintf("[%s] %s", "Warn", msg))
}

func (l logger) Error(msg string) {
	Logger.Output(3, fmt.Sprintf("[%s] %s", "Error", msg))
}

func (l logger) Fatal(msg string) {
	Logger.Output(3, fmt.Sprintf("[%s] %s", "Fatal", msg))
}
