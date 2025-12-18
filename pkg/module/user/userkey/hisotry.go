package userkey

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
	Entity UserKey `json:"entity,omitempty" gorm:"embedded;embeddedPrefix:h_"`
	EID    int64   `json:"h_eid,omitempty" gorm:"column:h_eid;comment:实体唯一标识,非全局唯一;" es:"type:long"`
}

func (h History) TableName(namer schema.Namer) string {
	if namer == nil {
		return constants.TableUserKeyHistory
	}
	return namer.TableName(constants.TableUserKeyHistory)
}

func (h History) AfterMigrator(db *gorm.DB, c *ctx.Context) error {
	scm := dorm.GetDbSchema(db)
	return history.CreateTrigger(db, scm, constants.TableUserKey)
}
