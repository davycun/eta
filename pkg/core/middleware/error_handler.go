package middleware

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/duke-git/lancet/v2/strutil"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ErrorHandler
// 统一错误处理
func ErrorHandler(c *gin.Context) {
	ctx.GetContext(c)
	defer func() {
		if a := recover(); a != nil {
			logger.OutputPanic(a)
			switch cErr := a.(type) {
			case *errs.ClientError:
				abortError(c, http.StatusBadRequest, cErr.Code, cErr.Message, nil)
			case *errs.ServerError:
				abortError(c, http.StatusInternalServerError, cErr.Code, cErr.Message, nil)
			case *errs.AuthError:
				abortError(c, http.StatusForbidden, cErr.Code, cErr.Message, nil)
			case *errs.BaseError:
				abortError(c, http.StatusInternalServerError, cErr.Code, cErr.Message, nil)
			default:
				msg := fmt.Sprintf("%v", a)
				message := strutil.Substring(msg, 0, 50)
				fullMessage := &msg
				ginMode := global.GetConfig().Server.GinMode
				if ginMode == "release" {
					fullMessage = nil
				}
				abortError(c, http.StatusInternalServerError, "unknown", message, fullMessage)
			}
		}
	}()
}
func abortError(c *gin.Context, status int, code, msg string, fullMessage *string) {
	c.AbortWithStatusJSON(status, dto.ControllerResponse{
		Code:        code,
		Message:     msg,
		FullMessage: fullMessage,
		Success:     false,
	})
}
