package service

import (
	"errors"
	"fmt"
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/filter"
	"github.com/davycun/eta/pkg/common/dorm/xa"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/core/service/hook"
	"strings"
)

func (s *DefaultService) Delete(args *dto.Param, result *dto.Result) error {

	var (
		err error
		cfg = hook.NewSrvConfig(iface.CurdModify, iface.MethodDelete, s.GetContext(), s.GetDB(), args, result, func(o *hook.SrvConfig) {
			//互相拷贝同步，以Service的配置优先
			o.SrvOptions.Merge(s.SrvOptions)
			s.SrvOptions.Merge(o.SrvOptions)
			o.EC = s.EC
		})
	)
	defer func() {
		_ = cfg.CommitOrRollback(err)
	}()

	cfg.CurDB = dorm.TableWithContext(cfg.TxDB, s.GetContext(), s.GetTableName())
	err = caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			return cfg.Before()
		}).
		Call(func(cl *caller.Caller) error {
			var (
				dbType = dorm.GetDbType(cfg.CurDB)
				sqs    = make([]string, 0, len(cfg.Values))
			)
			for _, v := range cfg.Values {
				id := v.FieldByName(entity.IdFieldName).String()
				updatedAt := v.FieldByName(entity.UpdatedAtFieldName).Int()
				sqs = append(sqs, fmt.Sprintf(`(%s='%s' and %s=%d)`, dorm.Quote(dbType, entity.IdDbName), id, dorm.Quote(dbType, entity.UpdatedAtDbName), updatedAt))
			}
			if len(sqs) < 1 {
				return nil
			}
			if len(args.AuthFilters) > 0 {
				cfg.CurDB = cfg.CurDB.Where(filter.ResolveWhereTable(cfg.GetTableName(), args.AuthFilters, dbType))
			}
			cfg.CurDB = cfg.CurDB.Where(strings.Join(sqs, " or ")).Delete(s.NewEntityPointer())
			return cfg.CurDB.Error
		}).
		Call(func(cl *caller.Caller) error {
			result.RowsAffected = cfg.CurDB.RowsAffected
			return cfg.After()
		}).Err

	return err
}
func (s *DefaultService) DeleteByFilters(args *dto.Param, result *dto.Result) error {

	var (
		err error
		cfg = hook.NewSrvConfig(iface.CurdModify, iface.MethodDeleteByFilters, s.GetContext(), s.GetDB(), args, result)
	)
	defer func() {
		if !dorm.InTransaction(s.GetDB()) {
			xa.CommitOrRollback(cfg.TxDB, err)
		}
	}()
	//TODO 是否允许全量删除，应该通过比对带filters的统计结果与全量总量是否一致来确定是否是全量更新
	//不严格的做法是判断filter的size，否则统计可能会影响性能
	if len(args.Filters) < 1 {
		return errors.New("不允许删除全量数据")
	}

	cfg.CurDB = dorm.TableWithContext(cfg.TxDB, s.GetContext(), s.GetTableName())
	err = caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			return cfg.Before()
		}).
		Call(func(cl *caller.Caller) error {
			wh := filter.ResolveWhereTable(cfg.GetTableName(), args.Filters, s.GetDbType())
			if wh == "" {
				return errors.New("不允许删除全量数据")
			}
			if len(args.AuthFilters) > 0 {
				cfg.CurDB = cfg.CurDB.Where(filter.ResolveWhereTable(cfg.GetTableName(), args.AuthFilters, s.GetDbType()))
			}
			cfg.CurDB = cfg.CurDB.Where(wh).Delete(s.NewEntityPointer())
			return cfg.CurDB.Error
		}).
		Call(func(cl *caller.Caller) error {
			result.RowsAffected = cfg.CurDB.RowsAffected
			return cfg.After()
		}).Err

	return err
}
