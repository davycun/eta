package dept_srv

import (
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/expr"
	"github.com/davycun/eta/pkg/common/dorm/filter"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/builder"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/core/service/sqlbd"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/dept"
)

// BuildListSql
// 该数据权限只在后台管理中启用:
// 1. 超级管理员和系统管理员角色可以在后台管理中看到所有的部门和所有的用户（包括虚拟用户）
// 2. 部门管理员可以在后台管理中看到他有部门管理员角色的所在部门及其子部门和部门及子部门的用户。（不包括虚拟用户）
// 3. 一般人如果被分配了后台管理的功能权限，那么在后台管理中只能看到他自己，看不到部门列表。
// disable_perm_filter 控制是否开启权限过滤
func buildListSql(cfg *hook.SrvConfig) (sqlList *sqlbd.SqlList, err error) {

	var (
		dbType   = dorm.GetDbType(cfg.OriginDB)
		scm      = dorm.GetDbSchema(cfg.OriginDB)
		idsAlias = "ids"
		args     = cfg.Param
	)

	if args.WithTree {
		// 部门数据可能会很大，所以必须是带条件并且查询树状结构的情况才允许loadAll
		if len(args.RecursiveFilters) < 1 && len(args.AuthRecursiveFilters) < 1 {
			args.RecursiveFilters = append(args.RecursiveFilters, filter.Filter{Column: "parent_id", Operator: filter.Eq, Value: ""})
		}
		//最多只能查两层
		if (args.TreeDepth < 1 || args.TreeDepth > 2) && !args.IsUp {
			args.TreeDepth = 1
		}
	}

	cte := builder.NewCteSqlBuilder(dbType, scm, constants.TableDept)
	if len(args.Columns) > 0 {
		cte.AddColumn(utils.Merge(args.Columns, "id", "updated_at", "parent_id", "name")...)
	} else {
		cte.AddColumn(dept.DefaultColumns...)
	}

	filterBd := buildListSqlBuilder(cfg, entity.IdDbName)
	if filterBd != nil {
		cte.With(idsAlias, filterBd)
		cte.Join("", idsAlias, entity.IdDbName, constants.TableDept, entity.IdDbName)
	}
	cte.AddOrderBy(args.OrderBy...)
	if !args.WithTree && !args.LoadAll {
		cte.Offset(args.GetOffset()).Limit(args.GetLimit())
	}
	listSql, countSql, err := cte.Build()

	sqlList = sqlbd.NewSqlList().AddSql(sqlbd.ListSql, listSql).AddSql(sqlbd.CountSql, countSql)
	return
}
func buildListSqlBuilder(cfg *hook.SrvConfig, idAlias string) *builder.CteSqlBuilder {

	var (
		scm     = dorm.GetDbSchema(cfg.GetDB())
		dbType  = dorm.GetDbType(cfg.GetDB())
		allSb   = make([]builder.Builder, 0, 4)
		rsAlias = "addr"
		args    = cfg.Param
	)
	if idAlias == "" {
		idAlias = entity.IdDbName
	}
	cte := builder.NewCteSqlBuilder(dbType, "", rsAlias)
	cte.AddExprColumn(expr.NewAliasColumn(entity.FromIdDbName, idAlias))

	if len(args.AuthRecursiveFilters) > 0 {
		authBd := builder.NewRecursiveSqlBuilder(dbType, scm, constants.TableDept).SetCteName("auth_cte")
		authBd.AddRecursiveFilter(args.AuthRecursiveFilters...).AddExprColumn(expr.NewAliasColumn(entity.IdDbName, entity.FromIdDbName))
		allSb = append(allSb, authBd)
	}

	if len(args.RecursiveFilters) > 0 {
		recBd := builder.NewRecursiveSqlBuilder(dbType, scm, constants.TableDept).
			SetUp(args.IsUp).AddRecursiveFilter(args.RecursiveFilters...).SetDepth(args.TreeDepth)
		recBd.AddExprColumn(expr.NewAliasColumn(entity.IdDbName, entity.FromIdDbName))
		allSb = append(allSb, recBd)
	}

	if len(args.Filters) > 0 {
		addBd := builder.NewSqlBuilder(dbType, scm, constants.TableDept).
			AddExprColumn(expr.NewAliasColumn(entity.IdDbName, entity.FromIdDbName)).
			AddFilter(args.Filters...)
		allSb = append(allSb, addBd)
	}

	if len(allSb) < 1 {
		return nil
	}
	first := allSb[0]
	for i := 1; i < len(allSb); i++ {
		switch bd := first.(type) {
		case *builder.SqlBuilder:
			bd.UnionIntersect(allSb[i])
		case *builder.RecursiveSqlBuilder:
			bd.UnionIntersect(allSb[i])
		}
	}
	cte.With(rsAlias, first)
	return cte
}
