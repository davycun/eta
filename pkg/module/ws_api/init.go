package ws_api

import (
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/core/ws"
	"github.com/gin-gonic/gin"
)

func InitModule() {
	go ws.HUB.Run()
	go ws.SubscribePushMessage()

	handler := Controller{}

	group := global.GetGin().Group("/ws")
	group.GET("/push", func(ctx *gin.Context) { createWs(ctx) }) // 公开
	group.POST("/push_test_msg", handler.pushTestMsg)            // 公开
}
