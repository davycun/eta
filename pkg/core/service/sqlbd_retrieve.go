package service

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
	"github.com/duke-git/lancet/v2/slice"
	"gorm.io/gorm"
)

func init() {
	sqlbd.AddSqlBuilder(sqlbd.BuildForAllTable, QuerySql, iface.MethodQuery)
	sqlbd.AddSqlBuilder(sqlbd.BuildForAllTable, QuerySql, iface.MethodList)
	sqlbd.AddSqlBuilder(sqlbd.BuildForAllTable, QuerySql, iface.MethodCount)
	sqlbd.AddSqlBuilder(sqlbd.BuildForAllTable, QuerySql, iface.MethodDetail)
	sqlbd.AddSqlBuilder(sqlbd.BuildForAllTable, AggregateSql, iface.MethodAggregate)
	sqlbd.AddSqlBuilder(sqlbd.BuildForAllTable, PartitionSql, iface.MethodPartition)
}

func QuerySql(cfg *hook.SrvConfig) (*sqlbd.SqlList, error) {
	listSql, countSql, err := BuildParamSql(cfg.GetDB(), cfg.Param, cfg.GetEntityConfig())
	return sqlbd.NewSqlList().AddSql(sqlbd.ListSql, listSql).AddSql(sqlbd.CountSql, countSql).SetEsRetriever(QueryFromEs), err
}

// BuildParamSql
// TODO 这里要思考什么情况下允许LoadAll，什么情况下不允许LoadAll？不解决LoadAll问题可能会把应用或数据库搞崩
// 1）通过配置实现？
// 2）如果cfg.param.LoadAll为true，并且非WithTree的时候，那么调用count接口查看大于多少条数据就不允许LoadAll，其他情况不允许LoadAll？
// 3）可以LoadAll，但只能通过websocket或SSE实现，服务端也是通过分页查询返回数据，直到获取完所有数据？
func BuildParamSql(db *gorm.DB, args *dto.Param, ec *iface.EntityConfig) (listSql, countSql string, err error) {

	var (
		//args           = cfg.Param
		dbType   = dorm.GetDbType(db)
		scm      = dorm.GetDbSchema(db)
		idsAlias = "ids"
		mustCols = entity.GetMustColumns(ec.NewEntityPointer())
		cols     = ResolveColumns(args, ec)
	)

	cte := builder.NewCteSqlBuilder(dbType, scm, ec.GetTableName())
	if len(cols) > 0 {
		cte.AddColumn(utils.Merge(cols, mustCols...)...)
	}
	cte.AddExprColumn(args.ExtraColumns...)

	if len(args.AuthRecursiveFilters) < 1 && len(args.RecursiveFilters) < 1 && len(args.Auth2RoleFilters) < 1 {
		cte.AddFilter(args.Filters...)
		cte.AddFilter(filter.KeywordToFilter(entity.RaContentDbName, args.SearchContent)...)
	} else {
		filterBd := buildIdListSqlBuilder(db, args, ec)
		if filterBd != nil {
			cte.With(idsAlias, filterBd)
			cte.Join("", idsAlias, entity.IdDbName, ec.GetTableName(), entity.IdDbName)
		}
	}

	cte.AddOrderBy(args.OrderBy...)
	if !args.WithTree && !args.LoadAll {
		cte.Offset(args.GetOffset()).Limit(args.GetLimit())
	}
	listSql, countSql, err = cte.Build()
	return
}
func buildIdListSqlBuilder(db *gorm.DB, args *dto.Param, ec *iface.EntityConfig, cols ...string) *builder.CteSqlBuilder {

	var (
		dbType  = dorm.GetDbType(db)
		allSb   = BuildParamFilterBuilder(db, args, ec, cols...)
		rsAlias = "rs"
		idAlias = entity.IdDbName
	)
	cte := builder.NewCteSqlBuilder(dbType, "", rsAlias)
	cte.AddColumn(idAlias)

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

func BuildParamFilterBuilder(db *gorm.DB, args *dto.Param, ec *iface.EntityConfig, cols ...string) []builder.Builder {

	var (
		scm    = dorm.GetDbSchema(db)
		dbType = dorm.GetDbType(db)
		allSb  = make([]builder.Builder, 0, 4)
	)
	if len(cols) < 1 {
		cols = []string{entity.IdDbName}
	}
	if len(args.Filters) > 0 || args.SearchContent != "" {
		tmpBd := builder.NewSqlBuilder(dbType, scm, ec.GetTableName()).AddColumn(cols...).AddFilter(args.Filters...)
		tmpBd.AddFilter(filter.KeywordToFilter(entity.RaContentDbName, args.SearchContent)...)
		allSb = append(allSb, tmpBd)
	}

	if len(args.AuthRecursiveFilters) > 0 {
		authBd := builder.NewRecursiveSqlBuilder(dbType, scm, ec.GetTableName()).SetCteName("auth_cte")
		authBd.AddRecursiveFilter(args.AuthRecursiveFilters...).AddColumn(cols...)
		allSb = append(allSb, authBd)
	}

	if len(args.RecursiveFilters) > 0 {
		tmpBd := builder.NewRecursiveSqlBuilder(dbType, scm, ec.GetTableName()).
			SetUp(args.IsUp).AddRecursiveFilter(args.RecursiveFilters...).SetDepth(args.TreeDepth)
		tmpBd.AddColumn(cols...)
		allSb = append(allSb, tmpBd)
	}

	if len(args.AuthFilters) > 0 {
		tmpBd := builder.NewSqlBuilder(dbType, scm, ec.GetTableName()).AddColumn(cols...).AddFilter(args.AuthFilters...)
		allSb = append(allSb, tmpBd)
	}
	if len(args.Auth2RoleFilters) > 0 {
		tmpBd := builder.NewSqlBuilder(dbType, scm, constants.TableAuth2Role).
			AddExprColumn(expr.NewAliasColumn(entity.FromIdDbName, entity.IdDbName)).
			AddFilter(args.AuthFilters...)
		allSb = append(allSb, tmpBd)
	}

	if len(allSb) < 1 {
		return nil
	}
	return allSb
}

func ResolveColumns(args *dto.Param, ec *iface.EntityConfig) []string {

	if len(args.ExcludeColumns) < 1 {
		return args.Columns
	}

	cols := args.Columns
	if len(cols) < 1 {
		cols = entity.GetTableColumns(ec.NewEntityPointer(), args.ExcludeColumns...)
	} else {
		cols = slice.Filter(cols, func(i int, v string) bool {
			return !utils.ContainAny(args.ExcludeColumns, v)
		})
	}
	return cols
}
