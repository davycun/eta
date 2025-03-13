package menu_srv

import (
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/controller"
	"github.com/davycun/eta/pkg/module/menu"
	"github.com/gin-gonic/gin"
	"strings"
)

var (
	baseUri = []string{"/role/*", "/permission/*", "/auth2role/*", "/cache/*", "/oauth2/*", "/authorize/*"}
)

func init() {
	baseUri = append(baseUri, "/user/set_current_dept", "/user/current", "/user/update", "/user/id_name", "/user/modify_password", "/user/reset_password")
	baseUri = append(baseUri, "/app/migrate")
	//baseUri = append(baseUri, "/menu/my_menu")
	baseUri = append(baseUri, "/menu/*")
	baseUri = append(baseUri, "/optlog/*")
	baseUri = append(baseUri, "/setting/*")
	baseUri = append(baseUri, "/storage/*")
	//baseUri = append(baseUri, "/data/*")
	baseUri = append(baseUri, "/template/*") //这个要确定是否都放开
	baseUri = append(baseUri, "/tasks/*")
	baseUri = append(baseUri, "/api/*")
	baseUri = append(baseUri, "/crypto/*")
	baseUri = append(baseUri, "/ws/*")
	///////// citizen
	baseUri = append(baseUri, "/neurond/*") //市民从v2同步到v3的接口
	baseUri = append(baseUri, "/citizen/address/*")
	baseUri = append(baseUri, "/citizen/addr2label/*")
	baseUri = append(baseUri, "/citizen/address_history/*")
	baseUri = append(baseUri, "/citizen/address_history/*")
	baseUri = append(baseUri, "/citizen/bd2label/*")
	baseUri = append(baseUri, "/citizen/bd2addr/*")
	///////// tourist_forecast
	baseUri = append(baseUri, "/tourist_forecast/*")
}

func ApiCallAuth(c *gin.Context) {
	var (
		ct  = ctx.GetContext(c)
		uri = c.Request.RequestURI
	)

	if ct.GetContextIsManager() {
		return
	}

	if global.IsAdminUri(uri) {
		controller.Fail(c, 403, "非超级管理员禁止操作", nil)
		return
	}

	if global.IsIgnoreUri(uri) {
		return
	}

	if strings.Contains(uri, "?") {
		uri, _, _ = strings.Cut(uri, "?")
	}

	if utils.IsMatchedUri(uri, baseUri...) {
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
