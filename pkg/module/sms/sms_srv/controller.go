package sms_srv

import (
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/core/controller"
	"github.com/davycun/eta/pkg/module/sms"
	"github.com/davycun/eta/pkg/module/sms/sms_sender"
	"github.com/duke-git/lancet/v2/random"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	controller.DefaultController
}

// SendVerifyCode
// Content内容示例
func SendVerifyCode(c *gin.Context) {
	var (
		param = &sms.VerifyCodeParam{}
		ct    = ctx.GetContext(c)
		appDb = ct.GetAppGorm()
	)
	err := controller.BindBody(c, param)
	if err != nil {
		controller.ProcessResult(c, nil, err)
	}

	tg := sms.Target{
		Mobile:        param.Phone,
		SignName:      param.SignName,
		TemplateCode:  param.TemplateCode,
		TemplateParam: param.TemplateParam,
		Content:       ctype.NewStringPrt(param.Content),
		TemplateKey:   param.TemplateKey,
	}
	smsCode := random.RandNumeral(6)
	err = sms_sender.SendVerifyCode(appDb, tg, smsCode)
	controller.ProcessResult(c, nil, err)
}
