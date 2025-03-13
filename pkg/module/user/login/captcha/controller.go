package captcha

import (
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/core/controller"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/gin-gonic/gin"
)

func SendSmsCode(c *gin.Context) {
	var (
		param  = &SmsCaptchaParam{}
		srv    = newService(ctx.GetContext(c))
		result = &dto.Result{}
	)

	err := controller.BindBody(c, &param)
	if err != nil {
		controller.ProcessResult(c, nil, err)
		return
	}
	err = srv.SendSmsCode(param, result)
	controller.ProcessResult(c, result, err)
}

func GenerateImage(c *gin.Context) {

	var (
		param  = &ImageCaptchaParam{}
		srv    = newService(ctx.GetContext(c))
		result = &dto.Result{}
	)

	err := controller.BindBody(c, &param)
	if err != nil {
		controller.ProcessResult(c, nil, err)
		return
	}
	err = srv.GenerateImage(param, result)
	controller.ProcessResult(c, result, err)

}
