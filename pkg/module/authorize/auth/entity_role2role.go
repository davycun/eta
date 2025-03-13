package auth

import (
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/eta/constants"
	"gorm.io/gorm/schema"
)

// Role2Role
// 角色和角色的关系，fromId是t_role的ID，toId是t_role或者t_department的ID
type Role2Role struct {
	entity.BaseEdgeEntity
}

func (o Role2Role) TableName(namer schema.Namer) string {
	if namer == nil {
		return constants.TableRole2Role
	}
	return namer.TableName(constants.TableRole2Role)
}
