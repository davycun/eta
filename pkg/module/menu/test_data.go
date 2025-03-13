package menu

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/mock/faker"
	"github.com/davycun/eta/pkg/core/entity"
)

func NewTestData() []Menu {
	ids := createIds()
	namespace := "eta"
	return []Menu{
		{BaseEntity: entity.BaseEntity{ID: ids[0]},
			Name:      fmt.Sprintf("%s_菜单", faker.Name()),
			Namespace: namespace,
			Api:       []Api{{Uri: "/*", Method: "*"}}, //为了测试的时候避免403
		},
		{BaseEntity: entity.BaseEntity{ID: ids[1]},
			Name:      fmt.Sprintf("%s_菜单", faker.Name()),
			ParentId:  ids[0],
			Namespace: namespace,
		},
		{BaseEntity: entity.BaseEntity{ID: ids[2]},
			Name:      fmt.Sprintf("%s_菜单", faker.Name()),
			ParentId:  ids[1],
			Namespace: namespace,
		},
		{BaseEntity: entity.BaseEntity{ID: ids[3]},
			Name:      fmt.Sprintf("%s_菜单", faker.Name()),
			ParentId:  ids[0],
			Namespace: namespace,
		},
	}
}
func createIds() []string {
	return []string{
		global.GenerateIDStr(),
		global.GenerateIDStr(),
		global.GenerateIDStr(),
		global.GenerateIDStr(),
	}
}
