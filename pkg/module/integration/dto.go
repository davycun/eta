package integration

import (
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/iface"
	jsoniter "github.com/json-iterator/go"
)

type CommandParam struct {
	Items []CommandParamItem `json:"items" binding:"required"`
}
type CommandParamItem struct {
	Command    string    `json:"command" binding:"required;oneof=create update update_by_filters delete delete_by_filters"`
	EntityCode string    `json:"entity_code" binding:"required"` // EntityConfig.Name, template.code
	Param      dto.Param `json:"param" binding:"required"`
}

type CommandResult struct {
	Items []CommandResultItem `json:"items" binding:"required"`
}
type CommandResultItem struct {
	Command    string      `json:"command"`
	EntityCode string      `json:"entity_code"`
	Result     *dto.Result `json:"result"`
}

type RefreshParam struct {
	EntityCode string `json:"entity_code" binding:"required"` // EntityConfig.Name, template.code
}
type RefreshWsResult struct {
	EntityCode string `json:"entity_code" binding:"required"` // EntityConfig.Name, template.code
	Status     string `json:"status,omitempty"`               // processing/finished
	Msg        string `json:"msg,omitempty"`                  // processing/finished
	CurrentId  string `json:"current_id,omitempty"`
}

func (i RefreshWsResult) ToString() string {
	toString, err := jsoniter.MarshalToString(i)
	if err != nil {
		return ""
	}
	return toString
}

type SyncArgs struct {
	//dsync.SyncArgs
	Srv iface.Service
}
