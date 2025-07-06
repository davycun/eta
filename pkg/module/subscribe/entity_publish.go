package subscribe

import (
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/eta/constants"
	"gorm.io/gorm/schema"
)

type Record struct {
	entity.BaseEntity
	SubId    string       `json:"sub_id,omitempty" gorm:"column:sub_id;comment:订阅ID"`
	Status   string       `json:"status,omitempty" gorm:"column:status;comment:推送状态"` //调用别人HTTP返回的状态码
	Method   iface.Method `json:"method,omitempty" gorm:"column:method;comment:当前的操作方法"`
	Request  string       `json:"request,omitempty" gorm:"column:request;comment:请求内容"`
	Response string       `json:"response,omitempty" gorm:"column:response;serializer:json;comment:调用的响应信息"`
	Count    int          `json:"count,omitempty" gorm:"column:count;comment:推送次数"`
}

func (s Record) TableName(namer schema.Namer) string {
	if namer != nil {
		return namer.TableName(constants.TablePublishRecord)
	}
	return constants.TablePublishRecord
}
