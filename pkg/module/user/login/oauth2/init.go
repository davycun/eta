package oauth2

import (
	"github.com/davycun/eta/pkg/common/global"
)

func InitModule() {
	handler := &Controller{}
	group := global.GetGin().Group("/oauth2")
	group.GET("/access_token", handler.AccessToken)   // openapi 获取 access token
	group.POST("/login", handler.LoginByUsername)     // 账号密码登录
	group.POST("/login_by_code", handler.LoginByCode) // 授权码登录
	group.POST("/logout", handler.Logout)             // 登出
	group.GET("/token_renewal", handler.TokenRenewal) // token 续期
}
