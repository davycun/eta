package menu_srv

import (
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/expr"
	"github.com/davycun/eta/pkg/common/dorm/filter"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/builder"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/core/service/sqlbd"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/menu"
	"gorm.io/gorm"
)

func buildListSql(cfg *hook.SrvConfig) (*sqlbd.SqlList, error) {

	var (
		db       = cfg.GetDB()
		dbType   = dorm.GetDbType(db)
		scm      = dorm.GetDbSchema(db)
		idsAlias = "ids"
		args     = cfg.Param
		sqlList  = sqlbd.NewSqlList(iface.MethodList, false)
	)
	if args.WithTree {
		// 部门数据可能会很大，所以必须是带条件并且查询树状结构的情况才允许loadAll
		if len(args.RecursiveFilters) < 1 && len(args.AuthRecursiveFilters) < 1 {
			args.RecursiveFilters = append(args.RecursiveFilters, filter.Filter{Column: "parent_id", Operator: filter.Eq, Value: ""})
		}
	}

	cte := builder.NewCteSqlBuilder(dbType, scm, constants.TableMenu)
	if len(args.Columns) > 0 {
		cte.AddColumn(utils.Merge(args.Columns, "id", "updated_at", "parent_id", "name")...)
	} else {
		cte.AddColumn(menu.DefaultColumns...)
	}

	filterBd := buildListSqlBuilder(db, args, entity.IdDbName)
	if filterBd != nil {
		cte.With(idsAlias, filterBd)
		cte.Join("", idsAlias, entity.IdDbName, constants.TableMenu, entity.IdDbName)
	}
	cte.AddOrderBy(args.OrderBy...)
	if !args.WithTree && !args.LoadAll {
		cte.Offset(args.GetOffset()).Limit(args.GetLimit())
	}
	listSql, countSql, err := cte.Build()
	sqlList.AddSql(sqlbd.ListSql, listSql).AddSql(sqlbd.CountSql, countSql)
	return sqlList, err
}
func buildListSqlBuilder(db *gorm.DB, args *dto.Param, idAlias string) *builder.CteSqlBuilder {

	var (
		scm     = dorm.GetDbSchema(db)
		dbType  = dorm.GetDbType(db)
		allSb   = make([]builder.Builder, 0, 4)
		rsAlias = "addr"
	)
	if idAlias == "" {
		idAlias = entity.IdDbName
	}
	cte := builder.NewCteSqlBuilder(dbType, "", rsAlias)
	cte.AddExprColumn(expr.NewAliasColumn(entity.FromIdDbName, idAlias))

	if len(args.AuthRecursiveFilters) > 0 {
		authBd := builder.NewRecursiveSqlBuilder(dbType, scm, constants.TableMenu).SetCteName("auth_cte")
		authBd.AddRecursiveFilter(args.AuthRecursiveFilters...).AddExprColumn(expr.NewAliasColumn(entity.IdDbName, entity.FromIdDbName))
		allSb = append(allSb, authBd)
	}

	if len(args.RecursiveFilters) > 0 {
		recBd := builder.NewRecursiveSqlBuilder(dbType, scm, constants.TableMenu).
			SetUp(args.IsUp).AddRecursiveFilter(args.RecursiveFilters...).SetDepth(args.TreeDepth)
		recBd.AddExprColumn(expr.NewAliasColumn(entity.IdDbName, entity.FromIdDbName))
		allSb = append(allSb, recBd)
	}

	if len(args.Filters) > 0 {
		addBd := builder.NewSqlBuilder(dbType, scm, constants.TableMenu).
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
