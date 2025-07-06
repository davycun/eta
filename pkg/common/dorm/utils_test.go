package dorm_test

import (
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJoinColumnsWithPrefixAlias(t *testing.T) {

	alias := dorm.JoinColumnsWithPrefixAlias(dorm.DaMeng, "", "emb_", []string{"id", "name", "id_type", "*"})
	assert.Equal(t, `"id" as "emb_id","name" as "emb_name","id_type" as "emb_id_type"`, alias)

	alias = dorm.JoinColumnsWithPrefixAlias(dorm.DaMeng, "t_people", "emb_", []string{"id", "name", "id_type", "*"})
	assert.Equal(t, `"t_people"."id" as "emb_id","t_people"."name" as "emb_name","t_people"."id_type" as "emb_id_type"`, alias)
}
func TestJoinColumns(t *testing.T) {

	alias := dorm.JoinColumns(dorm.DaMeng, "", []string{"id", "name", "id_type", "*"})
	assert.Equal(t, `"id","name","id_type",*`, alias)

	alias = dorm.JoinColumns(dorm.DaMeng, "t_people", []string{"id", "name", "id_type", "*"})
	assert.Equal(t, `"t_people"."id","t_people"."name","t_people"."id_type","t_people".*`, alias)
}
