package service

import (
	"errors"
	"fmt"
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/dorm/filter"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/iface"
	"gorm.io/gorm"
	"reflect"
)

// NewService
// 通过表的名字创建一个针对这个表实体的服务实例
// 传入c 的时候如果确定返回的iface.Service 需要在协程中运行，那么传入的Context需要是Clone的，避免接口过来的Context被回收
// 比如
func NewService(tableName string, c *ctx.Context, db *gorm.DB) (iface.Service, error) {
	ec, b := iface.GetEntityConfigByTableName(tableName)
	if !b {
		return nil, errors.New(fmt.Sprintf("can not found the table[%s] service", tableName))
	}
	//需要切换一个新Context，新的contextDb和contextTable，让调用者处理
	//c1 := c.Clone()
	//c1.SetContextGorm(db)
	iface.SetContextEntityConfig(c, &ec)
	if ec.NewService == nil {
		ec.NewService = NewServiceFactory(ec)
	}
	return ec.NewService(c, db, &ec), nil
}

// NewServiceFactory
// srvType是实现服务接口的类型，如果没有指定具体的srvType或者指定srvType类没有实现iface.Service，那么返回的服务工厂将会创建默认服务DefaultService
func NewServiceFactory(ec iface.EntityConfig) iface.NewService {
	if ec.NewService != nil {
		return ec.NewService
	}
	if ec.ServiceType == nil {
		return NewDefaultService
	}
	return func(c *ctx.Context, db *gorm.DB, ec *iface.EntityConfig) iface.Service {
		val := reflect.Value{}
		if ec.ServiceType.Kind() == reflect.Pointer {
			val = reflect.New(ec.ServiceType.Elem())
		} else {
			val = reflect.New(ec.ServiceType)
		}
		valInter := val.Interface()
		if srv, ok := valInter.(iface.Service); ok {
			srv.SetContext(c)
			srv.SetDB(db)
			srv.SetEntityConfig(ec)
			if iface.GetContextEntityConfig(c) == nil {
				iface.SetContextEntityConfig(c, ec)
			}
			return srv
		}
		logger.Errorf("the service type is not a iface.Service")
		return NewDefaultService(c, db, ec)
	}
}

func NewDefaultService(c *ctx.Context, db *gorm.DB, ec *iface.EntityConfig) iface.Service {
	var (
		srv = &DefaultService{}
	)
	srv.SetContext(c)
	srv.SetDB(db)
	srv.SetEntityConfig(ec)
	if iface.GetContextEntityConfig(c) == nil {
		iface.SetContextEntityConfig(c, ec)
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
func (s *SrvWrapper) Service() iface.Service {
	return s.srv
}
func (s *SrvWrapper) Result() *dto.Result {
	return s.rs
}
func (s *SrvWrapper) Query() error {
	if s.err != nil || s.srv == nil {
		return s.err
	}
	s.err = s.srv.Query(s.param, s.rs)
	return s.err
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
func (s *SrvWrapper) DeleteByFilters() error {
	if s.err != nil || s.srv == nil {
		return s.err
	}
	s.err = s.srv.DeleteByFilters(s.param, s.rs)
	return s.err
}
