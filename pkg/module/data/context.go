package data

import (
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/module/data/template"
)

const (
	templateContextKey = "templateContextKey" // 表配置
)

func GetContextTemplate(c *ctx.Context) *template.Template {
	var (
		u = &template.Template{}
	)
	value, exists := c.Get(templateContextKey)
	if exists {
		u = value.(*template.Template)
	}
	return u
}
func SetContextTemplate(c *ctx.Context, u *template.Template) {
	c.Set(templateContextKey, u)
}
