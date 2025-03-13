package ws

import (
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/duke-git/lancet/v2/strutil"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 30 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer.
	maxMessageSize = 128 * 1024 * 1024
)

// Client 客户端实例
type Client struct {
	Hub                 *Hub
	UserId              string
	LatestHeartBeatTime int64 // 暂时没用
	Conn                *websocket.Conn
	Send                chan *WsMessage
}

// read 读取数据通道
func (c *Client) Read() {
	defer func() {
		c.Hub.Unregister <- c
		err := c.Conn.Close()
		logger.Warnf("ws client[%s] closed!", c.UserId)
		if err != nil {
			logger.Warnf("ws client[%s] close error: %v", c.UserId, err)
			return
		}
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	err := c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		logger.Warnf("ws client[%s] SetReadDeadline error: %v", c.UserId, err)
		return
	}
	c.Conn.SetPongHandler(func(string) error {
		err = c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			logger.Warnf("ws client[%s] SetPongHandler.SetReadDeadline error: %v", c.UserId, err)
			return err
		}
		return nil
	})

	for {
		messageType, message, err := c.Conn.ReadMessage()
		if err != nil {
			// 客户端断开连接
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNoStatusReceived, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Errorf("ws client[%s] read error: %v", c.UserId, err)
			}
			break
		}
		// 暂时不对消息进行处理
		logger.Infof("ws client[%s] receive message: messageType:%d, messageContent:%v", c.UserId, messageType, string(message))
	}
}

// write 发送数据通道
func (c *Client) Write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		c.Hub.Unregister <- c
		ticker.Stop()
		err := c.Conn.Close()
		if err != nil {
			return
		}
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				// The Hub closed the channel.
				err := c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					return
				}
				return
			}
			err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				logger.Warnf("ws client[%s] SetWriteDeadline error: %v", c.UserId, err)
				return
			}

			var data = dto.ControllerResponse{
				Message: "success",
				Result: Message{
					MessageKey: message.MessageKey,
					Data:       message.Data,
				},
				Code:    "200",
				Success: true,
			}
			dataJson, _ := jsoniter.Marshal(data)
			err = c.Conn.WriteMessage(message.MessageType, dataJson)
			if err != nil {
				logger.Errorf("ws client[%s] send message error: %v", c.UserId, err)
				return
			}
			logger.Debugf("ws client[%s] send message: %v", c.UserId, strutil.BytesToString(dataJson))

		case <-ticker.C:
			err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				logger.Errorf("ws client[%s] SetWriteDeadline error: %v", c.UserId, err)
				return
			}
			if err = c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				logger.Errorf("ws client[%s] send PingMessage error: %v", c.UserId, err)
				return
			}
		}
	}
}
