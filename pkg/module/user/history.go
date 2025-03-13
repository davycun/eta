package user

import (
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/eta/constants"
	"gorm.io/gorm/schema"
)

type History struct {
	entity.History
	Entity User `json:"entity,omitempty" gorm:"embedded;embeddedPrefix:h_"`
}

func (h History) TableName(namer schema.Namer) string {
	if namer == nil {
		return constants.TableUserHistory
	}
	return namer.TableName(constants.TableUserHistory)
}
