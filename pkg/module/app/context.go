package app

import (
	"github.com/davycun/eta/pkg/common/ctx"
)

const (
	appContextKey = "appContextKey"
)

func GetContextApp(c *ctx.Context) (*App, bool) {
	var (
		u = &App{}
	)
	value, exists := c.Get(appContextKey)

	if exists {
		u = value.(*App)
	}
	return u, exists
}
func SetContextApp(c *ctx.Context, ap *App) {
	c.Set(appContextKey, ap)
}
