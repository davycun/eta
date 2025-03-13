package subscribe

import (
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/eta/constants"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Subscriber
// 回调接口的要求
// Content-Type：application/json
// Method为：POST
// 请求body：{"data":[{},{}]}
type Subscriber struct {
	entity.BaseEntity
	Method string  `json:"method,omitempty" gorm:"column:method;comment:请求方法" binding:"oneof=POST"`
	Header Headers `json:"header,omitempty" gorm:"column:header;serializer:json;comment:请求头"`
	Url    string  `json:"url,omitempty" gorm:"column:url;comment:请求地址"`
	Target string  `json:"target,omitempty" gorm:"column:target;comment:订阅目标" binding:"required"` //需要填写表名
}

func (s Subscriber) TableName(namer schema.Namer) string {
	if namer != nil {
		return namer.TableName(constants.TableSubscriber)
	}
	return constants.TableSubscriber
}

type Header struct {
	Key    string
	Values []string
}
type Headers []Header

func (h Headers) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return dorm.JsonGormDBDataType(db, field)
}

func (h Headers) GormDataType() string {
	return dorm.JsonGormDataType()
}

func (s Subscriber) AfterMigrator(db *gorm.DB, c *ctx.Context) error {
	return dorm.CreateUniqueIndex(db, constants.TableSubscriber, "url", "target")
}
