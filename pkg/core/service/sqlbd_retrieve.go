package service

import (
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/expr"
	"github.com/davycun/eta/pkg/common/dorm/filter"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/builder"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/core/ra"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/core/service/sqlbd"
	"github.com/davycun/eta/pkg/eta/constants"
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
	listSql, countSql, err := buildListSql(cfg)
	return sqlbd.NewSqlList().AddSql(sqlbd.ListSql, listSql).AddSql(sqlbd.CountSql, countSql).SetEsFilter(buildListFilter), err
}

// TODO 这里要思考什么情况下允许LoadAll，什么情况下不允许LoadAll？不解决LoadAll问题可能会把应用或数据库搞崩
// 1）通过配置实现？
// 2）如果cfg.Param.LoadAll为true，并且非WithTree的时候，那么调用count接口查看大于多少条数据就不允许LoadAll，其他情况不允许LoadAll？
// 3）可以LoadAll，但只能通过websocket或SSE实现，服务端也是通过分页查询返回数据，直到获取完所有数据？
func buildListSql(cfg *hook.SrvConfig) (listSql, countSql string, err error) {

	var (
		dbType         = dorm.GetDbType(cfg.OriginDB)
		scm            = dorm.GetDbSchema(cfg.OriginDB)
		idsAlias       = "ids"
		defaultColumns = entity.GetDefaultColumns(cfg.NewEntityPointer())
		mustCols       = entity.GetMustColumns(cfg.NewEntityPointer())
	)

	cte := builder.NewCteSqlBuilder(dbType, scm, cfg.GetTableName())
	if len(cfg.Param.Columns) > 0 {
		cte.AddColumn(utils.Merge(cfg.Param.Columns, mustCols...)...)
	} else {
		cte.AddColumn(defaultColumns...)
	}
	if len(cfg.Param.ExtraColumns) > 0 {
		cte.AddExprColumn(cfg.Param.ExtraColumns...)
	}

	filterBd := buildIdListSqlBuilder(cfg)
	if filterBd != nil {
		cte.With(idsAlias, filterBd)
		cte.Join("", idsAlias, entity.IdDbName, cfg.GetTableName(), entity.IdDbName)
	}
	cte.AddOrderBy(cfg.Param.OrderBy...)
	if !cfg.Param.LoadAll {
		cte.Offset(cfg.Param.GetOffset()).Limit(cfg.Param.GetLimit())
	}
	listSql, countSql, err = cte.Build()
	return
}
func buildIdListSqlBuilder(cfg *hook.SrvConfig) *builder.CteSqlBuilder {

	var (
		scm        = dorm.GetDbSchema(cfg.OriginDB)
		dbType     = dorm.GetDbType(cfg.OriginDB)
		allSb      = make([]builder.Builder, 0, 4)
		rsAlias    = "rs"
		filterList = make([]filter.Filter, 0, 2)
		idAlias    = entity.IdDbName
	)
	cte := builder.NewCteSqlBuilder(dbType, "", rsAlias)
	cte.AddColumn(idAlias)

	if cfg.Param.SearchContent != "" {
		filterList = append(filterList, filter.KeywordToFilter(entity.RaContentDbName, cfg.Param.SearchContent)...)
	}
	filterList = append(filterList, cfg.Param.Filters...)
	if len(filterList) > 0 {
		tmpBd := builder.NewSqlBuilder(dbType, scm, cfg.GetTableName()).AddColumn(entity.IdDbName).AddFilter(filterList...)
		allSb = append(allSb, tmpBd)
	}

	if len(cfg.Param.AuthRecursiveFilters) > 0 {
		authBd := builder.NewRecursiveSqlBuilder(dbType, scm, cfg.GetTableName()).SetCteName("auth_cte")
		authBd.AddRecursiveFilter(cfg.Param.AuthRecursiveFilters...).AddColumn(entity.IdDbName)
		allSb = append(allSb, authBd)
	}

	if len(cfg.Param.RecursiveFilters) > 0 {
		tmpBd := builder.NewRecursiveSqlBuilder(dbType, scm, cfg.GetTableName()).
			SetUp(cfg.Param.IsUp).AddRecursiveFilter(cfg.Param.RecursiveFilters...).SetDepth(cfg.Param.TreeDepth)
		tmpBd.AddColumn(entity.IdDbName)
		allSb = append(allSb, tmpBd)
	}

	if len(cfg.Param.AuthFilters) > 0 {
		tmpBd := builder.NewSqlBuilder(dbType, scm, cfg.GetTableName()).AddColumn(entity.IdDbName).AddFilter(cfg.Param.AuthFilters...)
		allSb = append(allSb, tmpBd)
	}
	if len(cfg.Param.Auth2RoleFilters) > 0 {
		tmpBd := builder.NewSqlBuilder(dbType, scm, constants.TableAuth2Role).
			AddExprColumn(expr.NewAliasColumn(entity.FromIdDbName, entity.IdDbName)).
			AddFilter(cfg.Param.AuthFilters...)
		allSb = append(allSb, tmpBd)
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

func buildListFilter(cfg *hook.SrvConfig) ([]filter.Filter, error) {

	var (
		allFilters = make([]filter.Filter, 0, len(cfg.Param.Filters))
	)

	if len(cfg.Param.RecursiveFilters) > 0 {
		flt, err := convertRecursiveFilter(cfg, cfg.Param.RecursiveFilters)
		if err != nil {
			return allFilters, err
		}
		allFilters = append(allFilters, flt)
	}
	if len(cfg.Param.AuthRecursiveFilters) > 0 {
		flt, err := convertRecursiveFilter(cfg, cfg.Param.AuthRecursiveFilters)
		if err != nil {
			return allFilters, err
		}
		allFilters = append(allFilters, flt)
	}

	if len(cfg.Param.Auth2RoleFilters) > 0 {
		flt, err := convertAuth2RoleFilter(cfg, cfg.Param.AuthRecursiveFilters)
		if err != nil {
			return allFilters, err
		}
		allFilters = append(allFilters, flt)
	}

	allFilters = append(allFilters, ra.KeywordToFilters(cfg.Ctx.GetAppGorm(), cfg.GetTableName(), cfg.Param.SearchContent)...)
	allFilters = append(allFilters, cfg.Param.Filters...)
	allFilters = append(allFilters, cfg.Param.AuthFilters...)

	return allFilters, nil
}

func convertRecursiveFilter(cfg *hook.SrvConfig, filterList []filter.Filter) (filter.Filter, error) {

	var (
		scm    = dorm.GetDbSchema(cfg.OriginDB)
		dbType = dorm.GetDbType(cfg.OriginDB)
		ids    []string
		flt    = filter.Filter{}
	)
	listSql, _, err := builder.NewSqlBuilder(dbType, scm, cfg.GetTableName()).AddFilter(filterList...).AddColumn(entity.IdDbName).Build()
	if err != nil {
		return flt, err
	}
	err = dorm.RawFetch(listSql, cfg.OriginDB, &ids)
	if err != nil {
		return flt, err
	}
	if len(ids) < 1 {
		return flt, err
	}
	return filter.Filter{Column: "parent_ids", Operator: filter.IN, Value: ids}, nil
}

func convertAuth2RoleFilter(cfg *hook.SrvConfig, filterList []filter.Filter) (filter.Filter, error) {
	var (
		scm    = dorm.GetDbSchema(cfg.OriginDB)
		dbType = dorm.GetDbType(cfg.OriginDB)
		ids    []string
		flt    = filter.Filter{}
	)
	listSql, _, err := builder.NewSqlBuilder(dbType, scm, constants.TableAuth2Role).AddFilter(filterList...).AddColumn(entity.FromIdDbName).Build()
	if err != nil {
		return flt, err
	}
	err = dorm.RawFetch(listSql, cfg.OriginDB, &ids)
	if err != nil {
		return flt, err
	}
	if len(ids) < 1 {
		return flt, err
	}
	return filter.Filter{Column: "id", Operator: filter.IN, Value: ids}, nil
}
