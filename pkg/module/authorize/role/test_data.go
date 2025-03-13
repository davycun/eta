package role

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/mock/faker"
	"github.com/davycun/eta/pkg/core/entity"
)

func NewTestData() []Role {
	ids := createIds()
	namespace := "delta"
	return []Role{
		{BaseEntity: entity.BaseEntity{ID: ids[0]},
			Name:      fmt.Sprintf("%s_角色1", faker.Name()),
			Namespace: namespace,
		},
		{BaseEntity: entity.BaseEntity{ID: ids[1]},
			Name:      fmt.Sprintf("%s_角色2", faker.Name()),
			Namespace: namespace,
		},
	}
}
func createIds() []string {
	return []string{
		global.GenerateIDStr(),
		global.GenerateIDStr(),
	}
}
