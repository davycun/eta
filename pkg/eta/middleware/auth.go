package middleware

import (
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/core/controller"
	"github.com/davycun/eta/pkg/module/app"
	"github.com/davycun/eta/pkg/module/authorize"
	"github.com/gin-gonic/gin"
)

func Auth(c *gin.Context) {
	uri := c.Request.RequestURI
	//TODO 把ignore做成一个管理功能
	if global.IsIgnoreUri(uri) {
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

func initIgnoreApiContext(c *gin.Context) {
	var (
		ct    = ctx.GetContext(c)
		appId = ct.GetContextAppId()
	)
	if appId != "" {
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
		ct.SetAppGorm(appDb)
	}
	ct.SetContextGorm(global.GetLocalGorm())
}
