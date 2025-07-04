package menu_srv

import (
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/controller"
	"github.com/davycun/eta/pkg/module/menu"
	"github.com/davycun/eta/pkg/module/setting"
	"github.com/gin-gonic/gin"
	"strings"
)

func ApiCallAuth(c *gin.Context) {
	var (
		ct     = ctx.GetContext(c)
		uri    = c.Request.URL.Path //c.Request.RequestURI 这个是 path?query
		method = c.Request.Method
	)

	if ct.GetContextIsManager() || setting.IsIgnoreTokenUri(nil, method, uri) || setting.IsIgnoreAuthUri(nil, method, uri) {
		return
	}

	if setting.IsAdminUri(nil, method, uri) && !ct.GetContextIsManager() {
		controller.Fail(c, 403, "非超级管理员禁止操作", nil)
		return
	}

	mn, err := menu.LoadMenuByUserId(ct)
	if err != nil {
		controller.Fail(c, 500, err.Error(), nil)
		return
	}
	//如果没有设置就都放过
	if len(mn) < 1 {
		return
	}

	for _, val := range mn {
		for _, v := range val.Api {
			mtd := strings.ToLower(v.Method)
			if utils.IsMatchedUri(uri, v.Uri) && (mtd == "" || strings.Contains(mtd, "*") || strings.Contains(mtd, strings.ToLower(c.Request.Method))) {
				return
			}
		}
	}
	controller.Fail(c, 403, "无权访问", nil)
}
