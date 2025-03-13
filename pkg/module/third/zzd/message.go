package zzd

import (
	"github.com/davycun/eta/pkg/common/logger"
	"net/http"
)

type Message struct {
	Zzd
}

type WorkNotificationResp struct {
	Success bool `json:"success"`
	Content struct {
		MsgId string `json:"msgId"`
	} `json:"content"`
}

/*
WorkNotification

工作通知

	https://openplatform-portal.dg-work.cn/portal/#/helpdoc?apiType=serverapi&docKey=2674860

工作通知接口限制说明：

	同一个应用相同消息的内容同一个用户一天只能接收一次。
	同一个应用给同一个用户发送消息，ISV应用开发方式一天不得超过100次。
	为保障系统稳定性，接口调用说明：
	1）单应用单次仅支持给1000人以内的工作通知发送,发送间隔不得低于10min ;禁止使用多线程方式调用。
	2）单应用发送工作通知总量当天不能超过2万条。
	3）当天通知人数超过2万人建议使用公告。
	4）当业务系统在短时间内发送过大，对通知上下游系统（IM，通讯录）引起稳定性隐患，会触发平台监控预警，针对可能引发系统风险的业务应用，平台会做下线处理，请合理调用。
	举例：当天3500人，可分4次发送，每次间隔10min。

参数:
  - organizationCodes	String	否	接收者的部门id列表， 接收者是部门id下(包括子部门下)的所有用户，与receiverIds都为空时不发送，最大长度列表跟receiverIds加起来不大于1000
  - receiverIds	String	否	接收人用户ID(accountId)， 多个人时使用半角逗号分隔，与organizationCodes都为空时不发送，最大长度列表跟organizationCodes加起来不大于1000
  - tenantId	String	是	租户ID
  - bizMsgId	String	是	业务消息id，自定义，有去重功能 调用者的业务数据ID，同样的ID调用多次会提示"重复"错误
  - msg String	是	json对象 必须 {"msgtype":"text","text":{"content":"消息内容"}} 消息内容，目前支持：文本消息：text, 链接消息：link, Markdown：markdown，OA消息：oa, 卡片消息：action_card。最长不超过2048个字节
*/
func (o *Message) WorkNotification(organizationCodes, receiverIds, tenantId, bizMsgId, msg string) *WorkNotificationResp {
	res := &WorkNotificationResp{}
	if o.Err != nil {
		return res
	}
	path := "/message/workNotification"
	params := buildParam(map[string]interface{}{
		"organizationCodes": organizationCodes,
		"receiverIds":       receiverIds,
		"tenantId":          tenantId,
		"bizMsgId":          bizMsgId,
		"msg":               msg,
	})
	header, query := o.signature(http.MethodPost, path, params)

	resp, err := o.client.R().
		SetHeaders(header).
		SetFormDataFromValues(query).
		SetError(&WorkNotificationResp{}).
		SetResult(&WorkNotificationResp{}).
		Post(path)

	if err != nil {
		o.Err = err
		return res
	}
	logger.Debugf("Zzd Message.WorkNotification resp: %s", resp)
	// {"success":true,"content":{"msgId":"080cdd6f-4cf3-4204-a583-5298adbc7b8c"}}
	if resp.IsError() {
		return resp.Error().(*WorkNotificationResp)
	}
	return resp.Result().(*WorkNotificationResp)
}
