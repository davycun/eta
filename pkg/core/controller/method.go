package controller

import (
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/service/ecf"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HandlerFunc interface {
	gin.HandlerFunc | ApiConfig
}

// Publish
// 如果methodList为空，只会发布POST接口
func Publish[T HandlerFunc](tableName string, path string, handler T, methodList ...string) {
	ec, b := ecf.GetEntityConfigByTableName(tableName)
	if !b {
		logger.Errorf("can not find entity config for %s", tableName)
	}
	var fc gin.HandlerFunc
	switch hd := any(handler).(type) {
	case ApiConfig:
		fc = NewApi(tableName, hd)
	case gin.HandlerFunc:
		fc = hd
		global.GetGin().POST(utils.AbsolutePath(ec.BaseUrl, path), hd)
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
