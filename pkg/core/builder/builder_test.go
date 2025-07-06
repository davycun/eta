package builder_test

import (
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/builder"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var (
	dbType = dorm.DaMeng
	scm    = "eta_dev_backend"
)

type bdArgs struct {
	bd       builder.Builder
	listSql  string
	countSql string
}

func TestSqlBuilder(t *testing.T) {
	argsList := []bdArgs{
		{
			bd:       builder.NewSqlBuilder(dorm.DaMeng, "eta_dev_backend", "t_people"),
			listSql:  `select  "t_people".* from "eta_dev_backend"."t_people"`,
			countSql: `select count( *) from "eta_dev_backend"."t_people"`,
		},
	}

	for _, v := range argsList {
		listSql, countSql, err := v.bd.Build()
		logger.Info(listSql)
		logger.Info(countSql)
		assert.Nil(t, err)
		assert.Equal(t, v.listSql, strings.TrimSpace(listSql))
		assert.Equal(t, v.countSql, strings.TrimSpace(countSql))
	}

}
func TestValueBuilder(t *testing.T) {
	argsList := make([]bdArgs, 0, 2)

	pepIds := []string{"1", "2", "3"}
	cte := builder.NewCteSqlBuilder(dbType, scm, "t_organization")
	vb := builder.NewValueBuilder(dbType, entity.FromIdDbName, pepIds...)
	cte.With("pep", vb)
	pep2OrgBd := builder.NewSqlBuilder(dbType, scm, "r_pep2org").AddColumn(entity.ToIdDbName).
		Join("", "pep", entity.FromIdDbName, "r_pep2org", entity.FromIdDbName)
	cte.With("pep2org", pep2OrgBd).
		AddTableColumn("t_organization", "*").
		Join("", "pep2org", entity.ToIdDbName, "t_organization", entity.IdDbName)

	argsList = append(argsList, bdArgs{
		bd:       cte,
		listSql:  `with "pep" as (select value as "from_id" from  jsonb_array_elements_text('["1","2","3"]')),"pep2org" as (select  "r_pep2org"."to_id" from "eta_dev_backend"."r_pep2org" join "pep" on "pep"."from_id" = "r_pep2org"."from_id"   ) select  "t_organization".* from "eta_dev_backend"."t_organization" join "pep2org" on "pep2org"."to_id" = "t_organization"."id"`,
		countSql: `with "pep" as (select value as "from_id" from  jsonb_array_elements_text('["1","2","3"]')),"pep2org" as (select  "r_pep2org"."to_id" from "eta_dev_backend"."r_pep2org" join "pep" on "pep"."from_id" = "r_pep2org"."from_id"   ) select count( *) from "eta_dev_backend"."t_organization" join "pep2org" on "pep2org"."to_id" = "t_organization"."id"`,
	})

	for _, v := range argsList {
		listSql, countSql, err := v.bd.Build()
		logger.Info(listSql)
		logger.Info(countSql)
		assert.Nil(t, err)
		assert.Equal(t, v.listSql, strings.TrimSpace(listSql))
		assert.Equal(t, v.countSql, strings.TrimSpace(countSql))
	}

}
