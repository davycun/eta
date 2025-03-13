package user2app

import (
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/eta/constants"
	"gorm.io/gorm/schema"
)

type User2App struct {
	entity.BaseEdgeEntity
	IsDefault *ctype.Boolean `json:"is_default,omitempty" gorm:"column:is_default;comment:是否为默认应用"` //当登录的时候没有指定app，那么就登录这个默认的app
	IsManager *ctype.Boolean `json:"is_manager,omitempty" gorm:"column:is_manager;comment:是否管理员"`
}

func (u User2App) TableName(namer schema.Namer) string {
	if namer == nil {
		return constants.TableUser2App
	}
	return namer.TableName(constants.TableUser2App)
}
