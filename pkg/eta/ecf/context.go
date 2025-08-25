package ecf

import (
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/eta/constants"
)

func GetContextEntityConfig(c *ctx.Context) *iface.EntityConfig {
	//var (
	//	u = &iface.EntityConfig{}
	//)
	value, exists := c.Get(constants.EntityConfigContextKey)
	if exists {
		return value.(*iface.EntityConfig)
	}
	return nil
}
func SetContextEntityConfig(c *ctx.Context, u *iface.EntityConfig) {
	c.Set(constants.EntityConfigContextKey, u)
}
