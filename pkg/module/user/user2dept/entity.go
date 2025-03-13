package user2dept

import (
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/eta/constants"
	"gorm.io/gorm/schema"
)

var (
	DefaultColumns = append(entity.DefaultEdgeColumns, "is_main", "post", "is_manager", "role_ids")
)

type User2Dept struct {
	entity.BaseEdgeEntity
	IsMain    ctype.Boolean      `json:"is_main" gorm:"column:is_main;comment:是否主部门"`
	Post      string             `json:"post,omitempty" gorm:"column:post;comment:在这个部门的岗位"` //一个人的岗位可以有多个，比如在这个部门是软件工程师，在另一个部门是产品
	IsManager ctype.Boolean      `json:"is_manager" gorm:"column:is_manager;comment:是否是管理员"`
	RoleIds   *ctype.StringArray `json:"role_ids,omitempty" gorm:"column:role_ids;comment:角色IDS"` //用户在当前部门下分配的角色
}

func (u User2Dept) TableName(namer schema.Namer) string {
	if namer == nil {
		return constants.TableUser2Dept
	}
	return namer.TableName(constants.TableUser2Dept)
}
