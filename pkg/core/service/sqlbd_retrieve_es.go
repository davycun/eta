package service

import (
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/dorm/es"
	"github.com/davycun/eta/pkg/common/dorm/filter"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/builder"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/ra"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/core/service/sqlbd"
	"github.com/davycun/eta/pkg/eta/constants"
	"gorm.io/gorm"
)

func QueryFromEs(cfg *hook.SrvConfig, sqlList *sqlbd.SqlList) error {
	esApi, err := BuildEsApiForQuery(cfg, sqlList)
	if err != nil {
		return err
	}

	var (
		rs = cfg.NewResultSlicePointer(cfg.Method)
	)
	if rs == nil {
		rs = &[]ctype.Map{}
	}

	esApi.Find(rs)
	cfg.Result.Data = rs
	cfg.Result.Total = esApi.Total

	return esApi.Err
}
func AggregateFromEs(cfg *hook.SrvConfig, sqlList *sqlbd.SqlList) error {
	esApi, err := BuildEsApiForAggregate(cfg, sqlList)
	if err != nil {
		return err
	}
	rs, err := esApi.Aggregate()
	if err != nil {
		return err
	}
	cfg.Result.Data = rs.Group
	cfg.Result.Total = rs.GroupTotal
	return err
}

func BuildEsApiForQuery(cfg *hook.SrvConfig, sqlList *sqlbd.SqlList) (*es.Api, error) {

	var (
		args          = cfg.Param
		esApi         = cfg.EsApi
		obj           = cfg.NewEntityPointer()
		cols          = ResolveColumns(cfg.Param, cfg.GetEntityConfig())
		mustCols      = entity.GetMustColumns(obj)
		esObj         = cfg.NewEsEntityPointer()
		parentIdsName = entity.GetParentIdsName(esObj)
		fltList, err  = ConvertParamFilterToEsFilters(cfg.GetDB(), args, cfg.GetTableName(), parentIdsName)
	)

	if err != nil {
		return esApi, err
	}

	if esApi == nil {
		esApi = es.NewApi(global.GetES(), cfg.GetEsIndexName())
	}

	if len(args.ExtraColumns) > 0 {
		logger.Error("当前 Query 查询有 ExtraColumns，ES 暂时不支持")
	}

	//ES如果不指定就获取所有字段
	if len(cols) > 0 {
		cols = utils.Merge(cols, mustCols...)
	}
	if args.OnlyCount {
		esApi.Offset(0).Limit(0)
	} else if !args.LoadAll {
		esApi.Offset(cfg.Param.GetOffset()).Limit(cfg.Param.GetLimit())
	}

	esApi = esApi.WithCount(cfg.Param.AutoCount || cfg.Param.OnlyCount).
		AddColumn(cols...).
		AddFilters(fltList...).
		OrderBy(cfg.Param.OrderBy...)

	return esApi, err
}
func BuildEsApiForAggregate(cfg *hook.SrvConfig, sqlList *sqlbd.SqlList) (*es.Api, error) {

	var (
		args       = cfg.Param
		esApi, err = BuildEsApiForQuery(cfg, sqlList)
	)
	if err != nil {
		return esApi, err
	}

	if esApi == nil {
		esApi = es.NewApi(global.GetES(), cfg.GetEsIndexName())
	}

	esApi.AddHaving(args.Having...).
		AddGroupCol(args.GroupColumns...).
		AddGroupAggCol(args.AggregateColumns...).
		OrderBy(args.OrderBy...)

	return esApi, err
}

// ConvertParamFilterToEsFilters
// 把dto.Param参数转换成适合ES查询的filters
// 当通过ES查询树结构表的时候，通常ES宽表会需要把父节点的ID都存储在"parent_ids"中，针对递归的filter就可以先从数据库顶点的ID，然后再作为ES的parent_ids字段的筛选
// parentIdsName 指定ESEntity中存储所有父节点ID的字段名称，默认为"parent_ids"
func ConvertParamFilterToEsFilters(db *gorm.DB, args *dto.Param, tableName string, parentIdsName string) ([]filter.Filter, error) {
	var (
		dbType     = dorm.GetDbType(db)
		allFilters = make([]filter.Filter, 0, len(args.Filters))
	)
	if parentIdsName == "" {
		parentIdsName = entity.ParentIdsDbName
	}

	if len(args.RecursiveFilters) > 0 {
		flt, err := ConvertRecursiveFilterToEsFilter(db, args.RecursiveFilters, tableName, parentIdsName)
		if err != nil {
			return allFilters, err
		}
		allFilters = append(allFilters, flt)
	}
	if len(args.AuthRecursiveFilters) > 0 {
		flt, err := ConvertRecursiveFilterToEsFilter(db, args.AuthRecursiveFilters, tableName, parentIdsName)
		if err != nil {
			return allFilters, err
		}
		allFilters = append(allFilters, flt)
	}
	if len(args.Auth2RoleFilters) > 0 {
		flt, err := ConvertAuth2RoleFilter(db, args.AuthRecursiveFilters)
		if err != nil {
			return allFilters, err
		}
		allFilters = append(allFilters, flt)
	}
	allFilters = append(allFilters, ra.KeywordToFilters(db, tableName, args.SearchContent, dbType)...)
	allFilters = append(allFilters, args.Filters...)
	allFilters = append(allFilters, args.AuthFilters...)

	return allFilters, nil
}

func ConvertRecursiveFilterToEsFilter(db *gorm.DB, filters []filter.Filter, tableName string, parentIdsName string) (filter.Filter, error) {

	var (
		scm    = dorm.GetDbSchema(db)
		dbType = dorm.GetDbType(db)
		ids    []string
		flt    = filter.Filter{}
	)
	listSql, _, err := builder.NewSqlBuilder(dbType, scm, tableName).AddFilter(filters...).AddColumn(entity.IdDbName).Build()
	if err != nil {
		return flt, err
	}
	err = dorm.RawFetch(listSql, db, &ids)
	if err != nil {
		return flt, err
	}
	if len(ids) < 1 {
		return flt, err
	}
	return filter.Filter{Column: parentIdsName, Operator: filter.IN, Value: ids}, nil
}

func ConvertAuth2RoleFilter(db *gorm.DB, filterList []filter.Filter) (filter.Filter, error) {
	var (
		scm    = dorm.GetDbSchema(db)
		dbType = dorm.GetDbType(db)
		ids    []string
		flt    = filter.Filter{}
	)
	listSql, _, err := builder.NewSqlBuilder(dbType, scm, constants.TableAuth2Role).AddFilter(filterList...).AddColumn(entity.FromIdDbName).Build()
	if err != nil {
		return flt, err
	}
	err = dorm.RawFetch(listSql, db, &ids)
	if err != nil {
		return flt, err
	}
	if len(ids) < 1 {
		return flt, err
	}
	return filter.Filter{Column: entity.IdDbName, Operator: filter.IN, Value: ids}, nil
}
