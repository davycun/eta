package ws_api

import (
	"github.com/davycun/eta/pkg/core/controller"
	"github.com/davycun/eta/pkg/core/ws"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	controller.DefaultController
}
type Param struct {
	MessageType int      `json:"message_type,omitempty"`
	UserId      []string `json:"user_id,omitempty"`
	MessageKey  string   `json:"message_key,omitempty"`
	Data        string   `json:"data,omitempty"`
}
type Result struct {
	Data string `json:"data"`
}

func (handler Controller) pushTestMsg(c *gin.Context) {
	param := new(Param)
	err := controller.BindBody(c, param)
	if err != nil {
		handler.ProcessResult(c, &Result{Data: "fail"}, err)
		return
	}
	ws.SendMessage(param.MessageKey, param.Data, param.UserId...)
	handler.ProcessResult(c, &Result{Data: "ok"}, nil)
}
