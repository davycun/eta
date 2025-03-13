package controller

import (
	"errors"
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/duke-git/lancet/v2/strutil"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"strconv"
)

func ProcessResult(c *gin.Context, data interface{}, err error) {
	addRespHeader(c)
	if errors.Is(err, errs.NoRecordAffected) {
		c.JSON(200, &dto.ControllerResponse{
			Code:    "404",
			Message: err.Error(),
			Success: true,
			Result:  data,
		})
		return
	}
	if err != nil && !errors.Is(err, errs.NoPermissionNoErr) {
		status, baseError := errs.HttpStatus(err)
		sentry.CaptureException(err)

		switch err.(type) {
		case *errs.ClientError, *errs.ServerError, *errs.AuthError, *errs.BaseError:
			Fail(c, status, baseError.Error(), data)
		default:
			msg := baseError.Error()
			message := strutil.Substring(msg, 0, 50)
			fullMessage := &msg
			ginMode := global.GetConfig().Server.GinMode
			if ginMode == "release" {
				fullMessage = nil
			}
			fail1(c, status, message, fullMessage, data)
		}

	} else {
		Success(c, data)
	}
}
func Success(c *gin.Context, data interface{}) {
	SuccessWithMsg(c, data, "ok")
}
func SuccessWithMsg(c *gin.Context, data interface{}, msg string) {
	c.JSON(200, &dto.ControllerResponse{
		Code:    "200",
		Message: msg,
		Success: true,
		Result:  data,
	})
}

func Fail(c *gin.Context, code int, msg string, data interface{}) {
	fail1(c, code, msg, nil, data)
}
func fail1(c *gin.Context, code int, msg string, fullMessage *string, data interface{}) {
	c.AbortWithStatusJSON(code, dto.ControllerResponse{
		Code:        strconv.Itoa(code),
		Message:     msg,
		FullMessage: fullMessage,
		Success:     false,
		Result:      data,
	})
}
func addRespHeader(c *gin.Context) {
	//ct, b := ctx.GetCurrentContext()
	ct := ctx.GetContext(c)
	if ct != nil {
		c.Header(constants.HeaderRespFetchCacheId, ct.GetString(ctx.FetchCacheContextKey))
		c.Header(constants.HeaderRespOpFromEs, ct.GetString(ctx.OpFromEsContextKey))
	}
}
