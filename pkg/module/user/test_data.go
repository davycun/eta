package user

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/mock/faker"
	"github.com/davycun/eta/pkg/core/entity"
	"math/rand"
)

func NewTestData() User {
	ids := createIds()
	name := faker.Name()
	account := fmt.Sprintf("%s_%d", name, rand.Int())
	return User{BaseEntity: entity.BaseEntity{ID: ids[0]},
		Name:     name,
		Account:  ctype.NewStringPrt(account),
		Password: "admin@123",
	}
}
func createIds() []string {
	return []string{
		global.GenerateIDStr(),
	}
}
