package dorm_test

import (
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQuotePlaceholder(t *testing.T) {

	rs := []struct {
		Src  string
		Vars []string
		Rs   map[dorm.DbType]string
	}{
		{
			Src:  `select distinct ?.? from ?`,
			Vars: []string{"cte1", "id", "cte1"},
			Rs: map[dorm.DbType]string{
				dorm.DaMeng:     `select distinct "cte1"."id" from "cte1"`,
				dorm.PostgreSQL: `select distinct "cte1"."id" from "cte1"`,
				dorm.Mysql:      "select distinct `cte1`.`id` from `cte1`",
				dorm.Doris:      "select distinct `cte1`.`id` from `cte1`",
			},
		},
		{
			Src:  `select ?.? from ?.?`,
			Vars: []string{"t_label", "*", "eta_dev_backend", "t_label"},
			Rs: map[dorm.DbType]string{
				dorm.DaMeng:     `select "t_label".* from "eta_dev_backend"."t_label"`,
				dorm.PostgreSQL: `select "t_label".* from "eta_dev_backend"."t_label"`,
				dorm.Mysql:      "select `t_label`.* from `eta_dev_backend`.`t_label`",
				dorm.Doris:      "select `t_label`.* from `eta_dev_backend`.`t_label`",
			},
		},
	}

	for _, v := range rs {
		for x, y := range v.Rs {
			tmp := dorm.QuotePlaceholder(x, v.Src, v.Vars...)
			assert.Equal(t, y, tmp)
		}

	}
}

func TestQuote(t *testing.T) {
	assert.Equal(t, `"test"."name"`, dorm.Quote(dorm.DaMeng, "test", "name"))
	assert.Equal(t, `"test"."name"."test"`, dorm.Quote(dorm.DaMeng, "test", "name", "test"))
	assert.Equal(t, `"name"."test"`, dorm.Quote(dorm.DaMeng, "", "name", "test"))
	assert.Equal(t, `"name"`, dorm.Quote(dorm.DaMeng, "", "name"))
	assert.Equal(t, `"t_label".*`, dorm.Quote(dorm.DaMeng, "t_label", "*"))
}
