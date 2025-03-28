package service

import (
	"errors"
	"fmt"
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/dorm/filter"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/core/service/ecf"
	"gorm.io/gorm"
	"reflect"
)

// NewService
// 通过表的名字创建一个针对这个表实体的服务实例
func NewService(tableName string, c *ctx.Context, db *gorm.DB) (iface.Service, error) {
	ec, b := ecf.GetEntityConfigByTableName(tableName)
	if !b {
		return nil, errors.New(fmt.Sprintf("can not found the table[%s] service", tableName))
	}
	//需要切换一个新Context，新的contextDb和contextTable
	c1 := c.Clone()
	c1.SetContextGorm(db)
	entity.SetContextTable(c1, ec.GetTable())
	if ec.NewService == nil {
		ec.NewService = NewServiceFactory(ec.ServiceType)
	}
	return ec.NewService(c1, db, ec.GetTable()), nil
}

// NewServiceFactory
// srvType是实现服务接口的类型，如果没有指定具体的srvType或者指定srvType类没有实现iface.Service，那么返回的服务工厂将会创建默认服务DefaultService
func NewServiceFactory(srvType reflect.Type) iface.NewService {
	if srvType == nil {
		return NewDefaultService
	}
	return func(c *ctx.Context, db *gorm.DB, tb *entity.Table) iface.Service {
		val := reflect.Value{}
		if srvType.Kind() == reflect.Pointer {
			val = reflect.New(srvType.Elem())
		} else {
			val = reflect.New(srvType)
		}
		valInter := val.Interface()
		if srv, ok := valInter.(iface.Service); ok {
			srv.SetContext(c)
			srv.SetDB(db)
			srv.SetTable(tb)
			if entity.GetContextTable(c) == nil {
				entity.SetContextTable(c, tb)
			}
			return srv
		}
		logger.Errorf("the service type is not a iface.Service")
		return NewDefaultService(c, db, tb)
	}
}

func NewDefaultService(c *ctx.Context, db *gorm.DB, tb *entity.Table) iface.Service {
	srv := &DefaultService{}
	srv.SetContext(c)
	srv.SetDB(db)
	srv.SetTable(tb)
	if entity.GetContextTable(c) == nil {
		entity.SetContextTable(c, tb)
	}
	return srv
}

type SrvWrapper struct {
	err   error
	srv   iface.Service
	param *dto.Param
	rs    *dto.Result
}

func NewSrvWrapper(tableName string, c *ctx.Context, db *gorm.DB) *SrvWrapper {
	wp := &SrvWrapper{
		param: &dto.Param{
			ModifyParam: dto.ModifyParam{
				SingleTransaction: true,
			},
		},
		rs: &dto.Result{},
	}
	wp.srv, wp.err = NewService(tableName, c, db)
	return wp
}

// SetData
// 这里必须传入实体的切片
func (s *SrvWrapper) SetData(dt any) *SrvWrapper {
	tp := reflect.TypeOf(dt)
	if tp.Kind() == reflect.Pointer {
		tp = tp.Elem()
	}
	if tp.Kind() != reflect.Slice {
		s.err = errors.New("the data must be a slice")
		return s
	}
	s.param.Data = dt
	return s
}
func (s *SrvWrapper) SetParam(fc func(param *dto.Param)) *SrvWrapper {
	fc(s.param)
	return s
}
func (s *SrvWrapper) AddFilters(flt ...filter.Filter) *SrvWrapper {
	s.param.Filters = append(s.param.Filters, flt...)
	return s
}
func (s *SrvWrapper) Create() error {
	if s.err != nil || s.srv == nil {
		return s.err
	}
	s.err = s.srv.Create(s.param, s.rs)
	return s.err
}
func (s *SrvWrapper) Update() error {
	if s.err != nil || s.srv == nil {
		return s.err
	}
	s.err = s.srv.Update(s.param, s.rs)
	return s.err
}
func (s *SrvWrapper) UpdateByFilters() error {
	if s.err != nil || s.srv == nil {
		return s.err
	}
	s.err = s.srv.UpdateByFilters(s.param, s.rs)
	return s.err
}
func (s *SrvWrapper) Delete() error {
	if s.err != nil || s.srv == nil {
		return s.err
	}
	s.err = s.srv.Delete(s.param, s.rs)
	return s.err
}
