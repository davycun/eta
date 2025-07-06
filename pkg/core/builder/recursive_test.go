package builder_test

import (
	"github.com/davycun/eta/pkg/common/dorm/filter"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/builder"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestNewRecursiveSqlBuilder(t *testing.T) {

	var (
		err error
	)
	assert.Nil(t, err)

	type SB struct {
		bd       builder.Builder
		listSql  string
		countSql string
	}

	sb1 := SB{
		bd:       builder.NewRecursiveSqlBuilder(dbType, scm, "t_address").AddRecursiveFilter(filter.Filter{Column: "id", Operator: filter.Eq, Value: "46"}),
		listSql:  `with  "cte"("id","parent_id") as (select  "t_address"."id","t_address"."parent_id" from "eta_dev_backend"."t_address"  where  ( "t_address"."id"  = '46')    union all (select  "t_address"."id","t_address"."parent_id" from "eta_dev_backend"."t_address" join "cte" on "cte"."id" = "t_address"."parent_id"   )),"cte_rs" as (select distinct "cte"."id" from "cte"    ) select  "t_address".* from "delta_dev_backend"."t_address" join "cte_rs" on "cte_rs"."id" = "t_address"."id"`,
		countSql: `with  "cte"("id","parent_id") as (select  "t_address"."id","t_address"."parent_id" from "eta_dev_backend"."t_address"  where  ( "t_address"."id"  = '46')    union all (select  "t_address"."id","t_address"."parent_id" from "eta_dev_backend"."t_address" join "cte" on "cte"."id" = "t_address"."parent_id"   )),"cte_rs" as (select distinct "cte"."id" from "cte"    ) select count( *) from "delta_dev_backend"."t_address" join "cte_rs" on "cte_rs"."id" = "t_address"."id"`,
	}

	b2 := builder.NewRecursiveSqlBuilder(dbType, scm, "t_address").
		AddRecursiveFilter(filter.Filter{Column: "id", Operator: filter.Eq, Value: "46"})
	b2.AddColumn("id", "parent_id", "name", "address")
	sb2 := SB{
		bd:       b2,
		listSql:  `with  "cte"("id","parent_id") as (select  "t_address"."id","t_address"."parent_id" from "eta_dev_backend"."t_address"  where  ( "t_address"."id"  = '46')    union all (select  "t_address"."id","t_address"."parent_id" from "eta_dev_backend"."t_address" join "cte" on "cte"."id" = "t_address"."parent_id"   )),"cte_rs" as (select distinct "cte"."id" from "cte"    ) select  "t_address"."id","t_address"."parent_id","t_address"."name","t_address"."address" from "delta_dev_backend"."t_address" join "cte_rs" on "cte_rs"."id" = "t_address"."id"`,
		countSql: `with  "cte"("id","parent_id") as (select  "t_address"."id","t_address"."parent_id" from "eta_dev_backend"."t_address"  where  ( "t_address"."id"  = '46')    union all (select  "t_address"."id","t_address"."parent_id" from "eta_dev_backend"."t_address" join "cte" on "cte"."id" = "t_address"."parent_id"   )),"cte_rs" as (select distinct "cte"."id" from "cte"    ) select count( *) from "delta_dev_backend"."t_address" join "cte_rs" on "cte_rs"."id" = "t_address"."id"`,
	}

	b3 := builder.NewRecursiveSqlBuilder(dbType, scm, "t_address").
		AddRecursiveFilter(filter.Filter{Column: "id", Operator: filter.Eq, Value: "46"})
	b3.AddColumn("id", "parent_id", "name", "address").AddFilter(filter.Filter{Column: "level", Operator: filter.Eq, Value: 1})
	sb3 := SB{
		bd:       b3,
		listSql:  `with  "cte"("id","parent_id") as (select  "t_address"."id","t_address"."parent_id" from "eta_dev_backend"."t_address"  where  ( "t_address"."id"  = '46')    union all (select  "t_address"."id","t_address"."parent_id" from "delta_dev_backend"."t_address" join "cte" on "cte"."id" = "t_address"."parent_id"   )),"cte_rs" as (select distinct "cte"."id" from "cte"    ) select  "t_address"."id","t_address"."parent_id","t_address"."name","t_address"."address" from "delta_dev_backend"."t_address" join "cte_rs" on "cte_rs"."id" = "t_address"."id" where  ( "t_address"."level"  = 1)`,
		countSql: `with  "cte"("id","parent_id") as (select  "t_address"."id","t_address"."parent_id" from "eta_dev_backend"."t_address"  where  ( "t_address"."id"  = '46')    union all (select  "t_address"."id","t_address"."parent_id" from "delta_dev_backend"."t_address" join "cte" on "cte"."id" = "t_address"."parent_id"   )),"cte_rs" as (select distinct "cte"."id" from "cte"    ) select count( *) from "delta_dev_backend"."t_address" join "cte_rs" on "cte_rs"."id" = "t_address"."id" where  ( "t_address"."level"  = 1)`,
	}

	b4 := builder.NewRecursiveSqlBuilder(dbType, scm, "t_address").
		AddRecursiveFilter(filter.Filter{Column: "id", Operator: filter.Eq, Value: "1"}).SetDepth(2)
	b4.AddColumn("id", "parent_id", "name", "address")
	sb4 := SB{
		bd:       b4,
		listSql:  `with  "cte"("id","parent_id","depth") as (select  "t_address"."id","t_address"."parent_id",1 as "depth"  from "eta_dev_backend"."t_address"  where  ( "t_address"."id"  = '1')    union all (select  "t_address"."id","t_address"."parent_id","cte"."depth" + 1 as "depth"  from "delta_dev_backend"."t_address" join "cte" on "cte"."id" = "t_address"."parent_id" where  ( "cte"."depth"  <= 2)   )),"cte_rs" as (select distinct "cte"."id" from "cte"    ) select  "t_address"."id","t_address"."parent_id","t_address"."name","t_address"."address" from "delta_dev_backend"."t_address" join "cte_rs" on "cte_rs"."id" = "t_address"."id"`,
		countSql: `with  "cte"("id","parent_id","depth") as (select  "t_address"."id","t_address"."parent_id",1 as "depth"  from "eta_dev_backend"."t_address"  where  ( "t_address"."id"  = '1')    union all (select  "t_address"."id","t_address"."parent_id","cte"."depth" + 1 as "depth"  from "delta_dev_backend"."t_address" join "cte" on "cte"."id" = "t_address"."parent_id" where  ( "cte"."depth"  <= 2)   )),"cte_rs" as (select distinct "cte"."id" from "cte"    ) select count( *) from "delta_dev_backend"."t_address" join "cte_rs" on "cte_rs"."id" = "t_address"."id"`,
	}

	bdList := make([]SB, 0, 4)
	bdList = append(bdList, sb1, sb2, sb3, sb4)

	for _, v := range bdList {
		listSql, countSql, err1 := v.bd.Build()
		assert.Nil(t, err1)
		logger.Info(listSql)
		logger.Info(countSql)
		assert.Equal(t, v.listSql, strings.TrimSpace(listSql))
		assert.Equal(t, v.countSql, strings.TrimSpace(countSql))
	}
}
