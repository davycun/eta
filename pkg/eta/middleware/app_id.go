package middleware

import (
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/controller"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/app"
	"github.com/gin-gonic/gin"
)

// 确保那些不需要权限的接口
func contextAppId(c *gin.Context) {

	//需要设置确定appID的基本都是不需要权限的请求
	if !global.IsIgnoreUri(c.Request.RequestURI) {
		return
	}
	rid := c.GetHeader(constants.HeaderAppId)
	if rid == "" {
		return
	}

	ct := ctx.GetContext(c)
	ct.SetContextAppId(rid)
	ap, err := app.LoadAppById(global.GetLocalGorm(), rid)
	if err != nil {
		logger.Errorf("LoadAppById[%s] Err %s", rid, err)
		controller.ProcessResult(c, nil, err)
		return
	}
	appDb, _ := global.LoadGorm(ap.Database)
	if appDb != nil {
		ct.SetAppGorm(appDb)
	}
}
