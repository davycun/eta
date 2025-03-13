package dorm

import (
	"github.com/davycun/eta/pkg/common/ctx"
	"gorm.io/gorm"
)

const (
	DbContextKey = "dbContextKey"
)

func WithContext(db *gorm.DB, c *ctx.Context) *gorm.DB {
	return db.Set(DbContextKey, c)
}

func GetDbContext(db *gorm.DB) *ctx.Context {
	if c, ok := db.Get(DbContextKey); ok {
		return c.(*ctx.Context)
	}
	return nil
}
