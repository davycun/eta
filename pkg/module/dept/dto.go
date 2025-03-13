package dept

import (
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/core/entity"
)

var (
	DefaultRelationDeptColumns = []string{"is_main", "post", "is_manager", "role_ids"}
)

type RelationDept struct {
	entity.BaseEdgeEntity
	Dept Department `json:"dept,omitempty" gorm:"embedded;embeddedPrefix:dp_"`

	//user2dept
	IsMain    bool               `json:"is_main,omitempty" gorm:"column:is_main;comment:是否主部门"`
	Post      string             `json:"post,omitempty" gorm:"column:post;comment:在这个部门的岗位"`
	IsManager bool               `json:"is_manager,omitempty" gorm:"column:is_manager;comment:是否是管理员"`
	RoleIds   *ctype.StringArray `json:"role_ids,omitempty" gorm:"column:role_ids;comment:角色IDS"`
}
