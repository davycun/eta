package middleware

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/gin-gonic/gin"
	"time"
)

var (
	ginLogSkipPaths = []string{"/health", "/health/", "/authorize/check", "/metrics"}
	ginLogConfig    = gin.LoggerConfig{
		Output:    logger.Writer(),
		SkipPaths: ginLogSkipPaths,
		Formatter: ginLogFormatter,
	}
)

func ginLogFormatter(param gin.LogFormatterParams) string {
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}

	if param.Latency > time.Minute {
		param.Latency = param.Latency.Truncate(time.Second)
	}
	return fmt.Sprintf("[GIN] %v |%s %3d %s| %13v | %15s |%s %-7s %s %#v\n%s",
		param.TimeStamp.Format("2006/01/02 - 15:04:05"),
		statusColor, param.StatusCode, resetColor,
		param.Latency,
		param.ClientIP,
		methodColor, param.Method, resetColor,
		param.Path,
		param.ErrorMessage,
	)
}
