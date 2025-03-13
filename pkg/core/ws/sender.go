package ws

import (
	"context"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
)

var (
	WsChannel = "eta:ws:push_message"
)

type WsMessage struct {
	Message
	MessageType int      `json:"message_type,omitempty"` // websocket.TextMessage
	UserId      []string `json:"user_id,omitempty"`      // 发送给哪些用户
}

// Message 消息
type Message struct {
	MessageKey string `json:"message_key,omitempty"` // 消息编码
	Data       any    `json:"data,omitempty"`        // 消息内容
}

/*
SendMessage 发送消息到客户端，如果UserId为空，则发送给所有客户端
把消息发送到redis，由redis的订阅者消费
*/
func SendMessage(messageKey, data string, userId ...string) {
	m := WsMessage{
		Message: Message{
			MessageKey: messageKey,
			Data:       data,
		},
		MessageType: websocket.TextMessage,
		UserId:      userId,
	}
	msgJson, err := jsoniter.Marshal(m)
	if err != nil {
		return
	}
	global.GetRedis().Publish(context.Background(), WsChannel, msgJson)
}

/*
SubscribePushMessage 订阅推送消息
*/
func SubscribePushMessage() {
	var (
		pubSub = global.GetRedis().Subscribe(context.Background(), WsChannel)
	)

	defer func() {
		err := pubSub.Close()
		if err != nil {
			logger.Errorf("close redis pubSub err %s", err)
		}
		if r := recover(); r != nil {
			logger.Errorf("redis subscribe from ws_push_message panic %v", r)
		}
	}()

	for msg := range pubSub.Channel() {
		logger.Debugf("redis subscribe receive msg %s from %s", msg.Payload, msg.Channel)
		wsMsg := &WsMessage{}
		err := jsoniter.Unmarshal(utils.StringToBytes(msg.Payload), &wsMsg)
		if err != nil {
			logger.Warnf("redis subscribe receive msg %s from %s, unmarshal error: %v", msg.Payload, msg.Channel, err)
		}
		send(wsMsg.UserId, wsMsg.MessageKey, wsMsg.Data)
	}
}

/*
send 发送消息到客户端，如果UserId为空，则发送给所有客户端

只发送给当前节点的连接
*/
func send(userId []string, messageKey string, data any) {
	HUB.Message <- &WsMessage{
		Message: Message{
			MessageKey: messageKey,
			Data:       data,
		},
		MessageType: websocket.TextMessage,
		UserId:      userId,
	}
}
