package captcha

import (
	"github.com/davycun/eta/pkg/common/global"
)

func InitModule() {
	group := global.GetGin().Group("/captcha")
	group.POST("/image_code", GenerateImage) // openapi 获取 access token
	group.POST("/sms_code", SendSmsCode)     // 发送验证码
}
