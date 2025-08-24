package template

import (
	"github.com/davycun/eta/pkg/common/ctx"
)

const (
	templateContextKey = "templateContextKey" // 表配置
)

func GetContextTemplate(c *ctx.Context) *Template {
	var (
		u = &Template{}
	)
	value, exists := c.Get(templateContextKey)
	if exists {
		u = value.(*Template)
	}
	return u
}
func SetContextTemplate(c *ctx.Context, u *Template) {
	c.Set(templateContextKey, u)
}
