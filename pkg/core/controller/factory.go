package controller

import (
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/core/service"
	"reflect"
)

func NewControllerFactory(ec iface.EntityConfig) iface.NewController {
	if ec.ControllerType == nil {
		return NewDefaultController
	}

	return func(srv iface.NewService) iface.Controller {
		var (
			val reflect.Value
		)
		if ec.ControllerType.Kind() == reflect.Pointer {
			val = reflect.New(ec.ControllerType.Elem())
		} else {
			val = reflect.New(ec.ControllerType)
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

func newController(ec iface.EntityConfig) iface.Controller {
	if ec.NewService == nil {
		ec.NewService = service.NewServiceFactory(ec)
	}
	if ec.NewController == nil {
		ec.NewController = NewControllerFactory(ec)
	}
	return ec.NewController(ec.NewService)
}
