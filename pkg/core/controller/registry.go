package controller

import (
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/core/service"
	"github.com/davycun/eta/pkg/core/service/ecf"
	"github.com/gin-gonic/gin"
	"path"
)

const (
	ApiPathCreate          = "create"
	ApiPathUpdate          = "update"
	ApiPathUpdateByFilters = "update_by_filters"
	ApiPathDelete          = "delete"
	ApiPathDeleteByFilters = "delete_by_filters"
	ApiPathQuery           = "query"
	ApiPathDetail          = "detail"
	ApiPathCount           = "count"
	ApiPathAggregate       = "aggregate"
	ApiPathPartition       = "partition"
	ApiPathImport          = "import"
	ApiPathExport          = "export"
)

var (
	controllerMap = map[string]ControlConfig{}
)

type ControlConfig struct {
	entityConfig ecf.EntityConfig
	control      iface.Controller
}

func (cc ControlConfig) GetEntityConfig() ecf.EntityConfig {
	return cc.entityConfig
}

func (cc ControlConfig) GetController() iface.Controller {
	if cc.control != nil {
		return cc.control
	}
	if cc.entityConfig.GetTable() == nil {
		return nil
	}
	if cc.entityConfig.NewService == nil {
		cc.entityConfig.NewService = service.NewServiceFactory(cc.entityConfig.ServiceType)
	}
	if cc.entityConfig.NewController == nil {
		cc.entityConfig.NewController = NewControllerFactory(cc.entityConfig.ControllerType)
	}
	cc.control = cc.entityConfig.NewController(cc.entityConfig.NewService)
	return cc.control
}

func registry(tableName string, cc ControlConfig) {
	if tableName == "" {
		logger.Warnf("registry ControlConfig err,because the tableName is empty,the baseUrl is %s ", cc.entityConfig.BaseUrl)
		return
	}
	controllerMap[tableName] = cc
}

func LoadController(tableName string) ControlConfig {
	if cc, ok := controllerMap[tableName]; ok {
		return cc
	}
	ec, b := ecf.GetEntityConfigByTableName(tableName)
	if !b {
		return ControlConfig{}
	}
	if ec.NewService == nil {
		ec.NewService = service.NewServiceFactory(ec.ServiceType)
	}
	if ec.NewController == nil {
		ec.NewController = NewDefaultController
	}
	handler := ec.NewController(ec.NewService)
	return ControlConfig{control: handler, entityConfig: ec}

}

func Registry(ec ecf.EntityConfig) *gin.RouterGroup {

	if ec.DisableApi {
		return nil
	}

	var (
		group   = global.GetGin().Group(ec.BaseUrl)
		handler iface.Controller
	)

	handler = newController(ec)
	tb := ec.GetTable()
	if tb != nil {
		registry(tb.GetTableName(), ControlConfig{entityConfig: ec, control: handler})
	}

	if len(ec.EnableMethod) < 1 && len(ec.DisableMethod) < 1 {
		group.POST("/create", handler.Create)
		group.POST("/update", handler.Update)
		group.POST("/update_by_filters", handler.UpdateByFilters)
		group.POST("/delete", handler.Delete)
		group.POST("/delete_by_filters", handler.DeleteByFilters)
		group.POST("/query", handler.Query)
		group.POST("/detail/:id", handler.Detail)
		group.POST("/count", handler.Count)
		group.POST("/aggregate", handler.Aggregate)
		group.POST("/partition", handler.Partition)
		group.POST("/import", handler.Import)
		group.POST("/export", handler.Export)
	} else {

		mp := make(map[iface.Method]gin.HandlerFunc)
		mp[iface.MethodCreate] = handler.Create
		mp[iface.MethodUpdate] = handler.Create
		mp[iface.MethodUpdateByFilters] = handler.Create
		mp[iface.MethodDelete] = handler.Create
		mp[iface.MethodDeleteByFilters] = handler.Create
		mp[iface.MethodQuery] = handler.Create
		mp[iface.MethodDetail] = handler.Create
		mp[iface.MethodCount] = handler.Create
		mp[iface.MethodAggregate] = handler.Create
		mp[iface.MethodPartition] = handler.Create
		mp[iface.MethodImport] = handler.Create
		mp[iface.MethodExport] = handler.Create

		for k, v := range mp {
			//如果允许的方法中不包含当前方法，则跳过
			if len(ec.EnableMethod) > 0 && !utils.ContainAny(ec.EnableMethod, k) {
				continue
			}
			//如果不允许的方法中包含当前方法，则跳过
			if len(ec.DisableMethod) > 0 && utils.ContainAny(ec.DisableMethod, k) {
				continue
			}
			group.POST(path.Join("/", k.String()), v)
		}
	}

	return group
}
