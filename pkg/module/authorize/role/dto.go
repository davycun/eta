package role

import "github.com/davycun/eta/pkg/core/entity"

// RelationRole 基本角色关联信息
type RelationRole struct {
	entity.BaseEdgeEntity
	Role Role `json:"role,omitempty" gorm:"embedded;embeddedPrefix:emb_"`
}
