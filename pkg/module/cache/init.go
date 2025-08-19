package cache

import (
	"github.com/davycun/eta/pkg/common/global"
)

func InitModule() {
	handler := &Controller{}

	group := global.GetGin().Group("/cache")
	group.POST("/scan", handler.scan)         // 获取 key
	group.GET("/detail/:key", handler.detail) // 查询一个key
	group.POST("/set", handler.set)           // 赋值
	group.POST("/del", handler.del)           // 删除
}
