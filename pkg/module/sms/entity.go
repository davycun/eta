package sms

import (
	"encoding/json"
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/eta/constants"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"strings"
)

// Task
// Task里面的Content是公共的内容，Target里面的是针对Target的Mobile的独立的内容
// 最终发送给Target.Mobile的内容是 Target.Content + Task.Content
// 比如Target.Content = "尊敬的张先生: "，Task.Content = "我们诚挚的邀请您参加2025年人工智能大会"，那么最终被发送的内容是 "尊敬的张先生: 我们诚挚的邀请您参加2025年人工智能大会"
// TODO 暂时的延时发送是通过短信供应商支持的
type Task struct {
	entity.BaseEntity
	Status      *ctype.String    `json:"status,omitempty" gorm:"column:status;comment:任务状态"` // 待发送 、已发送
	TargetDesc  *ctype.String    `json:"target_desc,omitempty" gorm:"column:target_desc;comment:接收短信的人的描述"`
	PlainTime   *ctype.LocalTime `json:"plain_time,omitempty" gorm:"column:plain_time;comment:计划短信发送时间"`
	SendReason  *ctype.String    `json:"send_reason,omitempty" gorm:"column:send_reason;comment:为什要发送短信"`
	Content     *ctype.String    `json:"content,omitempty" gorm:"column:content;comment:短信内容"`
	TargetTotal *ctype.Integer   `json:"target_total,omitempty" gorm:"target_total;comment:发送短信总数"`
	TargetList  []Target         `json:"target_list,omitempty" gorm:"-:all"` //对应任务发送的内容
}

func (o Task) TableName(namer schema.Namer) string {
	if namer != nil {
		return namer.TableName(constants.TableSmsTask)
	}
	return constants.TableSmsTask
}

// Target
// TaskId  task任务的ID
// Content //当前手机号发送的内容，最终发送的内容会拼接上Task里面的Content。
// 如果是模版发送
// 1）如果短信平台不支持模版发送，那么需要自己填写Content，Content示例：您的登录验证码为:{code} 其中code的值从TemplateParam获取
// 2)如果短信平台支持模版发送，那么就不需要填写Content，需要填写短信平台的TemplateCode和TemplateParam
type Target struct {
	entity.BaseEntity
	TaskId  string        `json:"task_id,omitempty" gorm:"column:task_id;comment:对应的任务ID"`
	Mobile  string        `json:"mobile,omitempty" gorm:"column:mobile;comment:手机号"`
	Content *ctype.String `json:"content,omitempty" gorm:"column:content;comment:短信内容"`
	Code    string        `json:"code,omitempty" gorm:"column:code;comment:短信供应商发送短信对应的ID"` //可以根据这个code去短信供应商查询短信的发送状态等

	//====下面这几个参数是在短信供应商只能根据模版发送短信的情况下使用====
	SignName      string    `json:"sign_name,omitempty" gorm:"column:sign_name;comment:发送消息的title"`      //针对不可以自定义发送短信内容的平台，只能通过模版发送短信，必填
	TemplateCode  string    `json:"template_code,omitempty" gorm:"column:template_code;comment:短信模版的编码"` //针对不可以自定义发送短信内容的平台，只能通过模版发送短信，必填
	TemplateParam ctype.Map `json:"template_param,omitempty" gorm:"column:template_param;serializer:json;comment:短信模板的参数"`
	TemplateKey   string    `json:"template_key,omitempty" gorm:"-:all"` //这是在backend.SMSConfig中templateMap中指定的key
}

func (t Target) ParamJson() string {
	pm, err := json.Marshal(t.TemplateParam)
	if err != nil {
		logger.Errorf("target json serialize err %s", err)
	}
	return string(pm)
}

func (t Target) TableName(namer schema.Namer) string {
	if namer != nil {
		return namer.TableName(constants.TableSmsTarget)
	}
	return constants.TableSmsTarget
}

func (t Target) ResolveContent() (string, error) {
	if !ctype.IsValid(t.Content) {
		return "", errs.NewClientError("the sms template content can not be empty")
	}
	if t.TemplateParam == nil {
		t.TemplateParam = ctype.Map{}
	}
	var (
		str       = ctype.ToString(t.Content)
		bts       = []byte(str)
		vb        = strings.Builder{}
		isVar     = false
		rsContent = strings.Builder{}
	)
	for _, v := range bts {
		if v == '{' {
			isVar = true
			continue
		}
		if v == '}' {
			isVar = false
			varStr := vb.String()
			if x, ok := t.TemplateParam[varStr]; ok {
				rsContent.WriteString(ctype.ToString(x))
			} else {
				return "", errs.NewClientError(fmt.Sprintf("can not found the param {%s}", varStr))
			}
			continue
		}

		if isVar {
			vb.WriteByte(v)
		} else {
			rsContent.WriteByte(v)
		}
	}
	return rsContent.String(), nil
}

type TargetList []Target

func (d TargetList) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return dorm.JsonGormDBDataType(db, field)
}

func (d TargetList) GormDataType() string {
	return dorm.JsonGormDataType()
}
