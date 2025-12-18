package app

import (
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/eta/constants"
	"gorm.io/gorm/schema"
)

type History struct {
	entity.History
	Entity App   `json:"entity,omitempty" gorm:"embedded;embeddedPrefix:h_"`
	EID    int64 `json:"h_eid,omitempty" gorm:"column:h_eid;comment:实体唯一标识,非全局唯一;" es:"type:long"`
}

func (h History) TableName(namer schema.Namer) string {
	if namer == nil {
		return constants.TableAppHistory
	}
	return namer.TableName(constants.TableAppHistory)
}
