package oauth2

import (
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/user"
)

var (
	loginFuncMap = map[string]LoginFunc{
		constants.LoginTypeAccount:       LoginByAccount,
		constants.LoginTypeDingService:   LoginByDingCode,
		constants.LoginTypeDingQrcode:    LoginByDingCode,
		constants.LoginTypeZzdService:    LoginByZzDingCode,
		constants.LoginTypeZzdQrcode:     LoginByZzDingQrCode,
		constants.LoginTypeWechatService: LoginByWechatCode,
		constants.LoginTypeWechatQrcode:  LoginByWechatCode,
		constants.LoginTypeWeComService:  LoginByWeComCode,
		constants.LoginTypeWeComQrcode:   LoginByWeComCode,
		constants.LoginTypeSmsService:    LoginBySmsCode,
		constants.LoginTypeAccessToken:   LoginByAccessToken,
	}
)

type LoginFunc func(c *ctx.Context, args any) (user.User, error)

func RegistryLoginFunc(loginType string, fc LoginFunc) {
	if _, ok := loginFuncMap[loginType]; ok {
		logger.Errorf("The login method %s already exists and will be overwritten", loginType)
	}
	loginFuncMap[loginType] = fc
}
