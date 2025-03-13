package service

import (
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/core/service/hook"
)

func (s *DefaultService) Create(args *dto.Param, result *dto.Result) error {
	var (
		err error
		cfg = hook.NewSrvConfig(iface.CurdModify, iface.MethodCreate, s.GetContext(), s.GetDB(), args, result)
	)
	defer func() {
		_ = cfg.CommitOrRollback(err)
	}()

	cfg.CurDB = dorm.TableWithContext(cfg.TxDB, s.GetContext(), s.GetTableName())
	//支持CreateOrUpdate
	if len(args.Conflict.Columns) > 0 || args.Conflict.OnConstraint != "" {
		cfg.CurDB = cfg.CurDB.Clauses(dto.ConvertConflict(args.Conflict))
	}

	err = caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			return cfg.Before()
		}).
		Call(func(cl *caller.Caller) error {
			cfg.CurDB = cfg.CurDB.Create(args.Data)
			return cfg.CurDB.Error
		}).
		Call(func(cl *caller.Caller) error {
			result.RowsAffected = cfg.CurDB.RowsAffected
			return cfg.After()
		}).Err
	return err
}
