package controller

import (
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/core/service"
	"github.com/davycun/eta/pkg/core/service/ecf"
	"reflect"
)

func NewControllerFactory(controllerType reflect.Type) iface.NewController {
	if controllerType == nil {
		return NewDefaultController
	}

	return func(srv iface.NewService) iface.Controller {
		var (
			val reflect.Value
		)
		if controllerType.Kind() == reflect.Pointer {
			val = reflect.New(controllerType.Elem())
		} else {
			val = reflect.New(controllerType)
		}
		valInter := val.Interface()
		if handler, ok := valInter.(iface.Controller); ok {
			handler.SetNewService(srv)
			return handler
		}
		logger.Errorf("the controller type is not a iface.Controller")
		return NewDefaultController(srv)
	}

}

func newController(ec ecf.EntityConfig) iface.Controller {
	if ec.NewService == nil {
		ec.NewService = service.NewServiceFactory(ec.ServiceType)
	}
	if ec.NewController == nil {
		ec.NewController = NewControllerFactory(ec.ControllerType)
	}
	return ec.NewController(ec.NewService)
}
