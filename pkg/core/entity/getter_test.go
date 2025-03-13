package entity_test

import (
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MyEntity struct {
	entity.BaseEdgeEntity
	MyName ctype.String `json:"my_name"`
}
type HistoryEntity struct {
	entity.History
	Entity MyEntity `json:"entity" gorm:"embedded;embeddedPrefix:h_"`
}

func TestGetter(t *testing.T) {

	e := MyEntity{
		BaseEdgeEntity: entity.BaseEdgeEntity{
			BaseEntity: entity.BaseEntity{
				ID: "abcd123",
			},
		},
		MyName: ctype.NewString("test", true),
	}
	en := HistoryEntity{
		Entity: e,
	}

	entity.GetString(e, "my_name")
	assert.Equal(t, "test", entity.GetString(e, "my_name"))
	assert.Equal(t, "test", entity.GetString(e, "MyName"))
	assert.Equal(t, "abcd123", entity.GetString(e, "ID"))
	assert.Equal(t, "abcd123", entity.GetString(e, "id"))
	assert.Nil(t, entity.Set(&e, "id", "123"))
	assert.Equal(t, "123", entity.GetString(e, "id"))

	assert.Nil(t, entity.Set(&en, "FromId", "abc"))
	assert.Equal(t, "abc", entity.GetString(&en, "FromId"))
}
