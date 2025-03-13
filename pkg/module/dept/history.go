package dept

import (
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/history"
	"github.com/davycun/eta/pkg/eta/constants"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type History struct {
	entity.History
	Entity Department `json:"entity,omitempty" gorm:"embedded;embeddedPrefix:h_"`
}

func (h History) TableName(namer schema.Namer) string {
	if namer == nil {
		return constants.TableDeptHistory
	}
	return namer.TableName(constants.TableDeptHistory)
}

func (h History) AfterMigrator(db *gorm.DB, c *ctx.Context) error {
	scm := dorm.GetDbSchema(db)
	return history.CreateTrigger(db, scm, constants.TableDept)
}
