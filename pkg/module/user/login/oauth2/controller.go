package oauth2

import (
	"github.com/davycun/eta/pkg/common/crypt"
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/core/controller"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/user"
	"github.com/gin-gonic/gin"
)

type Controller struct {
}

// LoginByUsername 账号密码登录
func (handler *Controller) LoginByUsername(c *gin.Context) {
	var (
		err    error
		param  = &LoginByUsernameParam{}
		result = &LoginResult{}
		args   = &LoginParam{Param: param, LoginType: constants.LoginTypeAccount}
	)
	err = controller.BindBody(c, param)
	if err != nil {
		controller.ProcessResult(c, nil, err)
		return
	}
	ct := ctx.GetContext(c)
	if param.AppId != "" {
		ct.SetContextAppId(param.AppId)
	}
	err = NewService(ct, global.GetLocalGorm()).Login(args, result)
	controller.ProcessResult(c, result, err)
}

// LoginByCode 通过编码登录
func (handler *Controller) LoginByCode(c *gin.Context) {
	var (
		err    error
		param  = &LoginByCodeParam{}
		result = &LoginResult{}
		args   = &LoginParam{Param: param, LoginType: param.LoginType}
	)
	err = controller.BindBody(c, param)
	if err != nil {
		controller.ProcessResult(c, nil, err)
		return
	}
	ct := ctx.GetContext(c)
	if param.AppId != "" {
		ct.SetContextAppId(param.AppId)
	}
	err = NewService(ct, global.GetLocalGorm()).Login(args, result)

	controller.ProcessResult(c, result, err)
}

// AccessToken openapi 获取 access token
// 通过access_key 获取 token
func (handler *Controller) AccessToken(c *gin.Context) {
	var (
		err    error
		param  = &AccessKeyParam{}
		result = &LoginResult{}
		args   = &LoginParam{Param: param, LoginType: constants.LoginTypeAccessToken}
	)
	err = c.BindQuery(&param)
	if param.Algo == "" {
		param.Algo = crypt.AlgoSignHmacSha256
	}
	err = NewService(ctx.GetContext(c), global.GetLocalGorm()).Login(args, result)
	controller.ProcessResult(c, result, err)
}

// Logout 登出
func (handler *Controller) Logout(c *gin.Context) {
	token := user.GetToken(ctx.GetContext(c))
	err := user.DelUserToken(token)
	controller.ProcessResult(c, nil, err)
}

func (handler *Controller) TokenRenewal(c *gin.Context) {
	var (
		err error
		ct  = ctx.GetContext(c)
	)

	token := user.GetToken(ct)
	err = user.RenewalToken(token, user.GetTokenExpireIn())
	controller.ProcessResult(c, nil, err)
}
