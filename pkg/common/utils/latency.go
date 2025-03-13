package utils

import (
	"github.com/davycun/eta/pkg/common/logger"
	"time"
)

func Latency(name string, f func()) {
	start := time.Now()
	f()
	latency := time.Now().Sub(start)
	logger.Infof(`%s的执行时间：%s`, name, latency.String())
}
