package setting

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/eta/constants"
	"golang.org/x/exp/maps"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

// Setting 配置在Local和app都有表的原因
// 1. 有的请求是无需token的，这个时候的上下文就不会存在gorm.DB，但是当前这个请求又需要一些配置信息，那么如果配置表只存在于appDB中就无法获取，所以LocalDB中得有配置表。
// 2. 如果配置表只存在于LocalDB中，那么如果appDB是一个独立的库，而且这个库是归属另一个团队管理的，因为我们不能把LocalDB权限释放给这个独立团队，那么他们的配置怎么办呢，只能存在appDB中。
// 3. 综上，所以LocalDB和AppDB中都存在Setting，并且LocalDB中的Setting是所有AppDB中的大合集，各自AppDB只能看到自己的配置
type Setting struct {
	entity.BaseEntity
	Namespace string     `json:"namespace,omitempty" gorm:"column:namespace;comment:命名空间" binding:"required"`
	Category  string     `json:"category,omitempty" gorm:"column:category;comment:类别" binding:"required" `
	Name      string     `json:"name,omitempty" gorm:"column:name;comment:名称" binding:"required" `
	Content   ctype.Json `json:"content,omitempty" gorm:"column:content;serializer:json;comment:具体配置内容"`
}

func (s Setting) GetKey() string {
	return fmt.Sprintf("%s_%s", s.Category, s.Name)
}

func (s Setting) TableName(namer schema.Namer) string {
	if namer == nil {
		return constants.TableSetting
	}
	return namer.TableName(constants.TableSetting)
}
func (s Setting) AfterMigrator(db *gorm.DB, c *ctx.Context) error {
	return caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			return dorm.CreateUniqueIndex(db, constants.TableSetting, "category", "name")
		}).
		// 如果没有数据就初始化一份数据
		Call(func(cl *caller.Caller) error {
			if len(defaultSettingMap) > 0 {
				setList := maps.Values(defaultSettingMap)
				cfl := clause.OnConflict{
					Columns: []clause.Column{
						{Name: "category"},
						{Name: "name"},
					},
					DoNothing: true,
				}
				return dorm.WithContext(db, c).Model(&setList).Clauses(cfl).Create(&setList).Error
			}
			return nil
		}).Err
}
