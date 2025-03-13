package middleware

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/controller"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/service/ecf"
	"github.com/davycun/eta/pkg/module/data"
	"github.com/davycun/eta/pkg/module/data/template"
	"github.com/davycun/eta/pkg/module/setting"
	"github.com/gin-gonic/gin"
	"strings"
)

func LoadTable(c *gin.Context) {
	if strings.HasPrefix(c.Request.RequestURI, "/data") {
		parseTemplate(c)
	} else {
		parseEntity(c)
	}
}

func parseEntity(c *gin.Context) {

	var (
		path  = c.Request.RequestURI
		ct    = ctx.GetContext(c)
		appDb = ct.GetAppGorm()
	)

	ec, ok := ecf.GetEntityConfigByUrl(path)
	if !ok {
		logger.Errorf("not found the EntityConfig which base path is [%s]", path)
		return
	}
	ecTb := ec.GetTable()
	bcTb, b := setting.GetTableConfig(appDb, ecTb.GetTableName())
	if b {
		bcTb.Merge(ecTb)
		entity.SetContextTable(ctx.GetContext(c), &bcTb)
	} else {
		entity.SetContextTable(ctx.GetContext(c), ecTb)
	}
	return
}

func parseTemplate(c *gin.Context) {
	var (
		code = struct {
			Code string `json:"code" uri:"code" binding:"required"`
		}{}
		ct   = ctx.GetContext(c)
		tmpl template.Template
	)

	err := controller.BindUri(c, &code)
	if err != nil {
		controller.ProcessResult(c, nil, err)
		return
	}
	if code.Code == "" {
		controller.ProcessResult(c, nil, errs.NewClientError(fmt.Sprintf("expect url is /data/:code/, but is url[%s]", c.Request.RequestURI)))
		return
	}
	tmpl, err = template.LoadByCode(ct.GetAppGorm(), code.Code)
	if err != nil {
		controller.ProcessResult(c, nil, err)
		return
	}
	entity.SetContextTable(ct, tmpl.GetTable())
	data.SetContextTemplate(ct, &tmpl)
	return
}
