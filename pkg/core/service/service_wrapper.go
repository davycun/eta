package service

import (
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/core/service/hook"
)

func (s *DefaultService) RetrieveWrapper(args *dto.Param, result *dto.Result, method iface.Method, fc HookFunc) error {
	var (
		err error
		cfg = hook.NewSrvConfig(iface.CurdRetrieve, method, s.GetContext(), s.GetDB(), args, result, func(o *hook.SrvConfig) {
			o.SrvOptions = s.SrvOptions
		})
	)

	err = caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			return cfg.Before()
		}).
		Call(func(cl *caller.Caller) error {
			return fc(cfg)
		}).Err
	if err != nil {
		return err
	}
	err = cfg.After()
	return err
}
func (s *DefaultService) ModifyWrapper(method iface.Method, args *dto.Param, result *dto.Result, fc HookFunc) error {
	var (
		err error
		cfg = hook.NewSrvConfig(iface.CurdModify, method, s.GetContext(), s.GetDB(), args, result, func(o *hook.SrvConfig) {
			o.SrvOptions = s.SrvOptions
		})
	)

	err = caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			return cfg.Before()
		}).
		Call(func(cl *caller.Caller) error {
			return fc(cfg)
		}).Err
	if err != nil {
		return err
	}
	err = cfg.After()
	return err
}

type HookFunc func(cfg *hook.SrvConfig) error
