package auth

import (
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/eta/constants"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	DefaultColumnsAuth2Role = append(entity.DefaultEdgeColumns, "origin_schema", "target_schema", "origin_table", "target_table")
)

// Auth2Role
// AuthType 说明
// 请参看 auth包下的常量，read、edit、delete、create、usage、admin、UsageWithRead
// read->filter rows use where
// edit->filter rows use where
// delete->filter rows use where
// usage->filter rows use where
// create->verify column content、verify url
// AuthTable 说明
// 是表示当前权限目标运用到哪张表，比如下面的一些场景的示例
// 1）菜单权限管理场景，from_id是菜单ID，to_id是RoleTable表的ID，AuthTable是t_menu
// 2）t_template表只能查询自己创建的数据，假设用户A的ID是1，那么在t_permission表创建一个filter{Column:"creator_id",Operator:"=",Value:1}
// from_id是t_permission的ID，to_id是RoleTable的ID，AuthTable是t_template
// 3）
type Auth2Role struct {
	entity.BaseEdgeEntity
	AuthType Type `json:"auth_type,omitempty" gorm:"column:auth_type;comment:权限的类型" binding:"required"`
	//from_id对应的表名
	FromTable string `json:"from_table,omitempty" gorm:"column:from_table;comment;from_id对应的表名;not null" binding:"required"`
	//to_id对应的表名，可能是t_role或者t_department
	//定义角色的表的名称
	ToTable string `json:"to_table,omitempty" gorm:"column:to_table;comment;from_id对应的表名;not null" binding:"required"`
	//权限运用到的实际业务表的表名
	//比如菜单与角色关联的时候，from_id是t_menu的ID即FromTable是t_menu，to_id是roleTable(可能是t_role或者t_department)的ID
	//比如permission表创建一个了只能查询自己创建数据的filter，需要把这个permission运用到t_people表
	//那么from_id是t_permission的ID（也就是FromTable是t_permission），to_id是RoleTable(可能是t_role或者t_department)的ID，AuthTable是t_people
	AuthTable string `json:"auth_table,omitempty" gorm:"column:auth_table;comment:权限数据的表名" binding:"required"`

	//RoleTable string `json:"role_table,omitempty" gorm:"column:role_table;comment:目标role的表名" binding:"required"` //大部门是t_role,也有可能是其他表充当role的角色，比如department表
}

func (o Auth2Role) TableName(namer schema.Namer) string {
	if namer == nil {
		return constants.TableAuth2Role
	}
	return namer.TableName(constants.TableAuth2Role)
}

func (o Auth2Role) AfterMigrator(db *gorm.DB, c *ctx.Context) error {
	return caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			return dorm.CreateUniqueIndex(db, constants.TableAuth2Role, "from_id", "to_id", "from_table", "to_table")
		}).
		Call(func(cl *caller.Caller) error {
			if len(defaultAuth2Role) < 1 {
				return nil
			}
			return dorm.Table(db, constants.TableAuth2Role).Create(&defaultAuth2Role).Error
		}).Err
}
