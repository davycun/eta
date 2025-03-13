package ws_api

import (
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/ws"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

var (
	// ws upGrader
	upGrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func createWs(c *gin.Context) {
	conn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Errorf("create ws error: %v", err)
		return
	}
	client := &ws.Client{
		Hub:                 ws.HUB,
		UserId:              ctx.GetContext(c).GetContextUserId(),
		LatestHeartBeatTime: time.Now().UnixMilli(),
		Conn:                conn,
		Send:                make(chan *ws.WsMessage),
	}
	go client.Read()
	go client.Write()
	client.Hub.Register <- client
}
