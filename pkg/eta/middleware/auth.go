package middleware

import (
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/core/controller"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/app"
	"github.com/davycun/eta/pkg/module/authorize"
	"github.com/davycun/eta/pkg/module/setting"
	"github.com/gin-gonic/gin"
)

func Auth(c *gin.Context) {
	var (
		uri    = c.Request.URL.Path
		method = c.Request.Method
	)
	if setting.IsIgnoreTokenUri(nil, method, uri) {
		initIgnoreApiContext(c)
		return
	}

	token := ctx.GetToken(ctx.GetContext(c))
	if token == "" {
		controller.Fail(c, 401, "非法访问", nil)
		return
	}
	authorize.Authorize(ctx.GetContext(c))
	c.Next()
	ctx.CleanCurrentContext(c)
}

// 解决确保那些不需要权限的接口，导致没有经过Auth的上下文初始化导致的问题
func initIgnoreApiContext(c *gin.Context) {
	var (
		ct    = ctx.GetContext(c)
		appId = ct.GetContextAppId()
	)
	if appId == "" {
		appId = c.GetHeader(constants.HeaderAppId)
	}

	if appId != "" {
		ct.SetContextAppId(appId)
		ap, err := app.LoadAppById(global.GetLocalGorm(), appId)
		if err != nil {
			controller.ProcessResult(c, nil, err)
			return
		}
		app.SetContextApp(ct, &ap)
		appDb, err := global.LoadGorm(ap.GetDatabase())
		if err != nil {
			controller.ProcessResult(c, nil, err)
			return
		}
		if appDb != nil {
			ct.SetAppGorm(appDb)
			ct.SetContextGorm(appDb)
		}
	} else {
		ct.SetContextGorm(global.GetLocalGorm())
	}
}
