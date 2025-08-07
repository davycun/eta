package controller

import (
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/gin-gonic/gin"
	"net/http"
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

type HandlerFunc interface {
	gin.HandlerFunc | ApiConfig
}

func Registry(ec iface.EntityConfig) *gin.RouterGroup {

	if ec.DisableApi {
		return nil
	}

	var (
		group   = global.GetGin().Group(ec.BaseUrl)
		handler iface.Controller
	)

	handler = newController(ec)

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

// Publish
// 如果methodList为空，只会发布POST接口
func Publish[T HandlerFunc](tableName string, path string, handler T, methodList ...string) {
	ec, b := iface.GetEntityConfigByTableName(tableName)
	if !b {
		logger.Errorf("can not find entity config for %s", tableName)
	}
	var fc gin.HandlerFunc
	switch hd := any(handler).(type) {
	case ApiConfig:
		fc = NewApi(tableName, hd)
	case gin.HandlerFunc:
		fc = hd
	}

	if len(methodList) < 1 {
		global.GetGin().POST(utils.AbsolutePath(ec.BaseUrl, path), fc)
	}
	if utils.ContainAny(methodList, http.MethodPost) {
		global.GetGin().POST(utils.AbsolutePath(ec.BaseUrl, path), fc)
	}
	if utils.ContainAny(methodList, http.MethodGet) {
		global.GetGin().GET(utils.AbsolutePath(ec.BaseUrl, path), fc)
	}
	if utils.ContainAny(methodList, http.MethodPut) {
		global.GetGin().PUT(utils.AbsolutePath(ec.BaseUrl, path), fc)
	}
	if utils.ContainAny(methodList, http.MethodDelete) {
		global.GetGin().DELETE(utils.AbsolutePath(ec.BaseUrl, path), fc)
	}
	if utils.ContainAny(methodList, http.MethodHead) {
		global.GetGin().HEAD(utils.AbsolutePath(ec.BaseUrl, path), fc)
	}
	if utils.ContainAny(methodList, http.MethodOptions) {
		global.GetGin().OPTIONS(utils.AbsolutePath(ec.BaseUrl, path), fc)
	}
	if utils.ContainAny(methodList, http.MethodPatch) {
		global.GetGin().PATCH(utils.AbsolutePath(ec.BaseUrl, path), fc)
	}
}
