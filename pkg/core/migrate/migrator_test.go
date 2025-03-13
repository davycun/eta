package migrate_test

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/core/migrate"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
)

func TestFormat(t *testing.T) {
	assert.Equal(t, "%s_t_wide_people", fmt.Sprintf("%s_%s", "%s", "t_wide_people"))
}

type myStruct struct {
}

func (m myStruct) AfterMigrator(db *gorm.DB, c *ctx.Context) error {
	return nil
}

type myStruct2 struct {
}

func (m *myStruct2) AfterMigrator(db *gorm.DB, c *ctx.Context) error {
	return nil
}

func TestMigrate(t *testing.T) {
	assert.True(t, impl(myStruct{}))
	assert.True(t, impl(&myStruct{}))
	assert.False(t, impl(myStruct2{}))
	assert.True(t, impl(&myStruct2{}))
}

func impl(obj any) bool {
	_, ok := obj.(migrate.MigratorAfter)
	return ok
}
