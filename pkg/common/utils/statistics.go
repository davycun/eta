package utils

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/logger"
	"time"
)

func Statistics(startMilliseconds int64, format string, args ...any) {
	str := fmt.Sprintf(format, args...)
	dur := time.Now().UnixMilli() - startMilliseconds
	logger.Infof(str+"| 耗时:%d.%d秒 |", dur/1000, dur%1000)
}
