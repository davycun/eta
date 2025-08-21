package service

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/dorm/expr"
	"github.com/davycun/eta/pkg/common/dorm/filter"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/common/run"
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
		err     error
		wg      = &sync.WaitGroup{}
		cfg     = hook.NewSrvConfig(iface.CurdRetrieve, method, s.SrvOptions, args, result)
		sqlList *sqlbd.SqlList
		extraRs = &sync.Map{}
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
			if s.EsRetrieveEnabled() && method != iface.MethodPartition {
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
			//把ListSql拆分出来独立处理的原因是，很多callback需要对结果集直接处理
			wg.Add(1)
			run.Go(func() {
				defer wg.Done()
				rs, err1 := s.extraSql(cfg, sqlbd.ListSql, sqlList.ListSql(), sqlList)
				err = errs.Cover(err, err1)
				result.Data = rs
			})
			return nil
		}).
		Call(func(cl *caller.Caller) error {
			for k, valSql := range sqlList.AllSql() {
				if k == sqlbd.CountSql || k == sqlbd.ListSql {
					continue
				}
				wg.Add(1)
				run.Go(func() {
					defer wg.Done()
					rs, err1 := s.extraSql(cfg, k, valSql, sqlList)
					err = errs.Cover(err, err1)
					extraRs.Store(k, rs)
				})
			}
			return nil
		}).Err
	wg.Wait()
	if err != nil {
		return err
	}
	rs := ctype.Map{}
	extraRs.Range(func(key, value any) bool {
		rs[key.(string)] = value
		return true
	})
	if len(rs) > 0 {
		rs[sqlbd.ListSql] = result.Data
		result.Data = rs
	}
	result.PageSize = args.GetPageSize()
	result.PageNum = args.GetPageNum()
	err = cfg.After()
	return err
}

// 避免extraRs可能会有并发写问题，所以用[sync.Map]
func (s *DefaultService) extraSql(cfg *hook.SrvConfig, sqlKey, valSql string, sqlList *sqlbd.SqlList) (any, error) {
	var (
		args   = cfg.Param
		method = cfg.Method
	)
	sqlRs := sqlList.NewResultSlicePointer(sqlKey)
	if sqlRs == nil && (sqlList.NeedScan || len(args.ExtraColumns) > 0) {
		colType := ctype.GetColType(s.NewResultPointer(method))
		ct := expr.ExplainColumnType(args.ExtraColumns...)
		for k1, v1 := range ct {
			colType[k1] = v1
		}
		listRs, err1 := ctype.ScanRows(cfg.OriginDB.Raw(valSql), colType)
		if sqlList.OnlyOne(sqlKey) && len(listRs) < 2 {
			if len(listRs) < 1 {
				return ctype.Map{}, err1
			} else {
				return listRs[0], err1
			}
		}
		return listRs, err1
	} else {
		listRs := sqlList.NewResultOrSlicePointer(sqlKey)
		if listRs == nil {
			if sqlList.OnlyOne(sqlKey) {
				listRs = s.NewResultPointer(method)
			} else {
				listRs = s.NewResultSlicePointer(method)
			}
		}
		if listRs == nil {
			return listRs, errs.NewServerError("the function NewResultSlicePointer return nil ")
		}
		err1 := dorm.RawFetch(valSql, cfg.OriginDB, listRs)
		return listRs, err1
	}
}
