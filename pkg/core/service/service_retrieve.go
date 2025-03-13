package service

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/dorm/expr"
	"github.com/davycun/eta/pkg/common/dorm/filter"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/run"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/core/service/sqlbd"
	"sync"
)

func (s *DefaultService) Query(args *dto.Param, result *dto.Result) error {
	return s.Retrieve(args, result, iface.MethodQuery)
}
func (s *DefaultService) Detail(args *dto.Param, result *dto.Result) error {
	return s.Retrieve(args, result, iface.MethodDetail)
}
func (s *DefaultService) DetailById(id int64, result *dto.Result) error {
	args := dto.Param{}
	args.Filters = []filter.Filter{{Column: entity.IdDbName, Operator: filter.Eq, Value: id}}
	return s.Detail(&args, result)
}

func (s *DefaultService) Retrieve(args *dto.Param, result *dto.Result, method iface.Method) error {
	var (
		err error
		wg  = &sync.WaitGroup{}
		cfg = hook.NewSrvConfig(iface.CurdRetrieve, method, s.GetContext(), s.GetDB(), args, result, func(o *hook.SrvConfig) {
			o.SrvOptions = s.SrvOptions
		})
		sqlList *sqlbd.SqlList
	)

	err = caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			return cfg.Before()
		}).
		Call(func(cl *caller.Caller) error {
			sqlList, err = sqlbd.Build(cfg, cfg.GetTableName(), method)
			return err
		}).
		Call(func(cl *caller.Caller) error {
			if s.RetrieveEnableEs() {
				cl.Stop()
				return s.RetrieveFromEs(cfg, sqlList)
			}
			return nil
		}).
		Call(func(cl *caller.Caller) error {
			if args.OnlyCount {
				cl.Stop()
				countSql := sqlList.CountSql()
				if countSql == "" {
					return errs.NewServerError(fmt.Sprintf("CountSql[%s] is empty", method))
				}
				return dorm.RawFetch(countSql, cfg.OriginDB, &result.Total)
			}
			return nil
		}).
		Call(func(cl *caller.Caller) error {
			if args.AutoCount {
				wg.Add(1)
				countSql := sqlList.CountSql()
				if countSql == "" {
					return errs.NewServerError(fmt.Sprintf("CountSql[%s] is empty", method))
				}
				run.Go(func() {
					defer wg.Done()
					err = errs.Cover(err, dorm.RawFetch(countSql, cfg.OriginDB, &result.Total))
				})
			}
			return nil
		}).
		Call(func(cl *caller.Caller) error {
			wg.Add(1)
			listSql := sqlList.ListSql()
			if listSql == "" {
				return errs.NewServerError(fmt.Sprintf("ListSql[%s] is empty", method))
			}
			if len(args.ExtraColumns) > 0 || sqlList.NeedScan {
				run.Go(func() {
					defer wg.Done()
					colType := ctype.GetColType(s.NewRsDataPointer())
					ct := expr.ExplainColumnType(args.ExtraColumns...)
					for k, v := range ct {
						colType[k] = v
					}
					dt, err1 := ctype.ScanRows(cfg.OriginDB.Raw(listSql), colType)
					err = errs.Cover(err, err1)
					result.Data = dt
				})
			} else {
				run.Go(func() {
					defer wg.Done()
					listRs := s.NewRsDataSlicePointer()
					err = errs.Cover(err, dorm.RawFetch(listSql, cfg.OriginDB, listRs))
					result.Data = listRs
				})
			}
			return nil
		}).Err
	wg.Wait()
	if err != nil {
		return err
	}
	err = cfg.After()
	return err
}

func (s *DefaultService) RetrieveFromEs(cfg *hook.SrvConfig, sqlList *sqlbd.SqlList) error {

	if sqlList == nil || sqlList.EsFilter == nil {
		return errs.NewServerError(fmt.Sprintf("[%s]没有指定BuildEsFilter函数", cfg.Method))
	}
	var (
		esApi         = cfg.EsApi
		cols          = entity.GetDefaultColumns(s.NewEntityPointer())
		esFilter, err = sqlList.EsFilter(cfg)
	)

	if err != nil {
		return err
	}
	s.GetContext().Set(ctx.OpFromEsContextKey, "1")

	if len(cfg.Param.ExtraColumns) > 0 {
		logger.Error("当前 Query 查询有 ExtraColumns，ES 暂时不支持")
	}

	if len(cfg.Param.Columns) > 0 {
		cols = utils.Merge(cfg.Param.Columns, entity.GetMustColumns(s.NewEntityPointer())...)
	}
	if cfg.Param.OnlyCount {
		esApi.Offset(0).Limit(0)
	} else {
		esApi.Offset(cfg.Param.GetOffset()).Limit(cfg.Param.GetLimit())
	}

	if sqlList.IsAgg {
		aggRs, err1 := esApi.AddFilters(esFilter...).
			WithCount(cfg.Param.AutoCount || cfg.Param.OnlyCount).
			AddGroupCol(cfg.Param.GroupColumns...).
			AddGroupAggCol(cfg.Param.AggregateColumns...).
			AddHaving(cfg.Param.Having...).
			AddAggCol(sqlList.EsAggCol(cfg)...).
			OrderBy(cfg.Param.OrderBy...).
			Offset(cfg.Param.GetOffset()).
			Limit(cfg.Param.GetLimit()).
			Aggregate()
		if err1 != nil {
			return err1
		}
		cfg.Result.Total = esApi.Total
		cfg.Result.Data = aggRs.Group
	} else {
		listRs := s.NewRsDataSlicePointer()
		if sqlList.NeedScan {
			listRs = make([]ctype.Map, 0, 10)
		}

		esApi = esApi.WithCount(cfg.Param.AutoCount || cfg.Param.OnlyCount).
			AddColumn(cols...).
			AddFilters(esFilter...).
			OrderBy(cfg.Param.OrderBy...)
		if cfg.Param.LoadAll && !cfg.Param.OnlyCount {
			esApi = esApi.LoadAll(&listRs)
		} else {
			esApi = esApi.Find(&listRs)
		}
		cfg.Result.Total = esApi.Total
		cfg.Result.Data = listRs
	}

	return esApi.Err
}
