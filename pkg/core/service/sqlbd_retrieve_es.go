package service

import (
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/dorm/es"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/core/service/sqlbd"
)

func QueryFromEs(cfg *hook.SrvConfig, sqlList *sqlbd.SqlList) error {
	esApi, err := BuilderEsApiForQuery(cfg, sqlList)
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
	esApi, err := BuilderEsApiForAggregate(cfg, sqlList)
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

func BuilderEsApiForQuery(cfg *hook.SrvConfig, sqlList *sqlbd.SqlList) (*es.Api, error) {

	var (
		err      error
		args     = cfg.Param
		esApi    = cfg.EsApi
		obj      = cfg.NewEntityPointer()
		cols     = cfg.Param.Columns
		mustCols = entity.GetMustColumns(obj)
	)

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
		AddFilters(args.Filters...).
		AddFilters(args.AuthFilters...).
		OrderBy(cfg.Param.OrderBy...)

	return esApi, err
}
func BuilderEsApiForAggregate(cfg *hook.SrvConfig, sqlList *sqlbd.SqlList) (*es.Api, error) {

	var (
		args       = cfg.Param
		esApi, err = BuilderEsApiForQuery(cfg, sqlList)
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
