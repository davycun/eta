package service

import (
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/core/service/ecf"
	"gorm.io/gorm"
	"reflect"
)

// NewService
// 通过表的名字创建一个针对这个表实体的服务实例
func NewService(tableName string, c *ctx.Context, db *gorm.DB) iface.Service {
	ec, b := ecf.GetEntityConfigByTableName(tableName)
	if !b {
		logger.Errorf("没有配置对应的EntityConfig[%s]", tableName)
		return nil
	}
	//需要切换一个新Context，新的contextDb和contextTable
	c1 := c.Clone()
	c1.SetContextGorm(db)
	entity.SetContextTable(c1, ec.GetTable())
	if ec.NewService == nil {
		ec.NewService = NewServiceFactory(ec.ServiceType)
	}
	return ec.NewService(c1, db, ec.GetTable())
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
