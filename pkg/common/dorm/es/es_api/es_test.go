package es_api_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.idatatlas.com/delta/pkg/common/dorm/es"
	"gitlab.idatatlas.com/delta/pkg/delta/constants"
	"gitlab.idatatlas.com/delta/pkg/delta/core/entity"
	"gitlab.idatatlas.com/delta/pkg/diy/citizen/people/people_entity"
	"gitlab.idatatlas.com/delta/pkg/diy/citizen/wide/wide_pep"
)

func TestCast(t *testing.T) {

	p1 := people_entity.People{BaseEntity: entity.BaseEntity{ID: "1"}}
	p2 := &people_entity.People{BaseEntity: entity.BaseEntity{ID: "2"}}
	p3 := wide_pep.PeopleWide{People: people_entity.People{BaseEntity: entity.BaseEntity{ID: "3"}}}

	type Args struct {
		ID *string
	}

	id := "4"
	assert.Equal(t, "1", getId(p1))
	assert.Equal(t, "2", getId(p2))
	assert.Equal(t, "3", getId(p3))
	assert.Equal(t, "4", getId(&Args{ID: &id}))
}

func getId(obj any) string {
	return es.GetStringValue(reflect.ValueOf(obj), entity.IdFieldName)
}

func TestIndexName(t *testing.T) {

	assert.Equal(t, "test_t_people", es.GetIndexName("test", people_entity.People{}))
	assert.Equal(t, "test_t_people", es.GetIndexName("test", constants.TablePeople))
	assert.Equal(t, "test_t_wide_people", es.GetIndexName("test", constants.TableWidePeople))

}
