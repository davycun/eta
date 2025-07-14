package iface

import "github.com/davycun/eta/pkg/common/ctx"

const (
	entityConfigContextKey = "entityConfigContextKey" // 表配置
)

func GetContextEntityConfig(c *ctx.Context) *EntityConfig {
	var (
		u = &EntityConfig{}
	)
	value, exists := c.Get(entityConfigContextKey)
	if exists {
		u = value.(*EntityConfig)
	}
	return u
}
func SetContextEntityConfig(c *ctx.Context, u *EntityConfig) {
	c.Set(entityConfigContextKey, u)
}
