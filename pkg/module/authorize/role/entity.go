package role

import (
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/eta/constants"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

const (
	CategoryDefault  = "默认"
	CategoryOther    = "其他"
	CategoryPosition = "职位"
)

var (
	DefaultColumns = append(entity.DefaultVertexColumns, "namespace", "target", "name", "category")
)

type Role struct {
	entity.BaseEntity
	Namespace string `json:"namespace,omitempty" gorm:"column:namespace;comment:命名空间" binding:"required"` //为了区分不同的定制化项目
	Category  string `json:"category,omitempty" gorm:"column:category;comment:类似于分组的功能"`                  //
	Name      string `json:"name,omitempty" gorm:"column:name;comment:角色名称" binding:"required"`
	Target    string `json:"target,omitempty" gorm:"column:target;comment:这个角色主要用在什么目标上"` //比如主要用在菜单上
}

func (r Role) TableName(namer schema.Namer) string {
	if namer == nil {
		return constants.TableRole
	}
	return namer.TableName(constants.TableRole)
}
func (r Role) DefaultColumns() []string {
	return DefaultColumns
}

func (r Role) AfterMigrator(db *gorm.DB, c *ctx.Context) error {
	//如果没有数据就初始化一份数据
	return caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			return createRole(db, c)
		}).Err
}

func createRole(db *gorm.DB, c *ctx.Context) error {
	// 初始化两种角色

	return caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			if len(defaultRoleList) < 1 {
				return nil
			}
			cfl := clause.OnConflict{
				Columns: []clause.Column{
					{Name: entity.IdDbName},
				},
				UpdateAll: true,
			}
			return dorm.WithContext(db, c).Model(&defaultRoleList).Clauses(cfl).Create(&defaultRoleList).Error
		}).Err
}
