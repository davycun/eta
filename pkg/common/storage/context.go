package storage

import "github.com/davycun/eta/pkg/common/ctx"

func GetContextStorage(c *ctx.Context) Storage {
	s, exists := c.Get(ctx.StorageContextKey)
	if exists {
		return s.(Storage)
	}
	return DefaultStorage
}
func SetContextStorage(c *ctx.Context, st Storage) {
	c.Set(ctx.StorageContextKey, st)
}
