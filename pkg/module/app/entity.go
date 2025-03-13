package app

import (
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/eta/constants"
	"gorm.io/gorm/schema"
)

var (
	DefaultAppColumns = append(entity.DefaultVertexColumns, "name", "logo", "slogan", "company")
)

type App struct {
	entity.BaseEntity
	Name      string         `json:"name,omitempty" binding:"required"`
	Logo      string         `json:"logo,omitempty"`
	Slogan    string         `json:"slogan,omitempty"`
	Company   string         `json:"company,omitempty"`
	Valid     ctype.Boolean  `json:"valid,omitempty"  gorm:"column:valid;comment:启用或禁用"`
	IsDefault *ctype.Boolean `json:"is_default,omitempty" gorm:"column:is_default;comment:是否默认的APP"` //平台初始化的时候的APP
	Database  dorm.Database  `json:"database,omitempty" gorm:"column:database;serializer:json;comment:数据库连接信息"`
}

func (a *App) TableName(namer schema.Namer) string {
	if namer == nil {
		return constants.TableApp
	}
	return namer.TableName(constants.TableApp)
}

func (a *App) GetDatabase() dorm.Database {
	return a.Database
}
func (a *App) SetDatabase(db dorm.Database) {
	a.Database = db
}
func (a *App) DefaultColumns() []string {
	return DefaultAppColumns
}
