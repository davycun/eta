package entity

import (
	"github.com/davycun/eta/pkg/common/ctx"
)

const (
	tableContextKey = "tableContextKey" // 表配置
)

func GetContextTable(c *ctx.Context) *Table {
	var (
		u = &Table{}
	)
	value, exists := c.Get(tableContextKey)
	if exists {
		u = value.(*Table)
	}
	return u
}
func SetContextTable(c *ctx.Context, u *Table) {
	c.Set(tableContextKey, u)
}

func GetTableName2(c *ctx.Context) any {
	return GetContextTable(c).GetTableName()
}
func NewEntityPointer(c *ctx.Context) any {
	return GetContextTable(c).NewEntityPointer()
}
func NewEntitySlicePointer(c *ctx.Context) any {
	return GetContextTable(c).NewEntitySlicePointer()
}
func NewRsDataPointer(c *ctx.Context) any {
	return GetContextTable(c).NewRsDataPointer()
}
func NewRsDataSlicePointer(c *ctx.Context) any {
	return GetContextTable(c).NewRsDataSlicePointer()
}
