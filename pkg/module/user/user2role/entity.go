package user2role

import (
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/eta/constants"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	DefaultColumns = entity.DefaultEdgeColumns
)

type User2Role struct {
	entity.BaseEdgeEntity
}

func (u User2Role) TableName(namer schema.Namer) string {
	if namer == nil {
		return constants.TableUser2Role
	}
	return namer.TableName(constants.TableUser2Role)
}

func (u User2Role) AfterMigrator(db *gorm.DB, c *ctx.Context) error {
	return dorm.CreateUniqueIndex(db, constants.TableUser2Role, entity.FromIdDbName, entity.ToIdDbName)
}
