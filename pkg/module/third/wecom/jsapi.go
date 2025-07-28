package wecom

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/dromara/dongle"
)

type Jsapi struct {
	Wecom
}

type GetJsapiTicketResp struct {
	Errcode   int    `json:"errcode"`
	Errmsg    string `json:"errmsg"`
	Ticket    string `json:"ticket"`
	ExpiresIn int    `json:"expires_in"`
}

/*
GetCorpJsapiTicket 获取企业的jsapi_ticket

https://developer.work.weixin.qq.com/document/22677
*/
func (j *Jsapi) GetCorpJsapiTicket() *GetJsapiTicketResp {
	res := &GetJsapiTicketResp{}
	if j.Err != nil {
		return res
	}
	path := "/cgi-bin/get_jsapi_ticket"
	params := map[string]string{
		"access_token": j.Wecom.GetAccessToken(),
	}
	resp, err := j.client.R().
		SetQueryParams(params).
		SetError(&GetJsapiTicketResp{}).
		SetResult(&GetJsapiTicketResp{}).
		Get(path)

	if err != nil {
		j.Err = err
		return res
	}
	logger.Debugf("Wecom Jsapi.GetCorpJsapiTicket resp: %s", resp)
	// { "errcode":0, "errmsg":"ok", "ticket":"bxLdikRXVbTPdHSM05e5u5sUoXNKd8-41ZO3MhKoyN5OfkWITDGgnr2fwJ0m9E8NYzWKVZvdVtaUgWvsdshFKA", "expires_in":7200 }
	if resp.IsError() {
		return resp.Error().(*GetJsapiTicketResp)
	}
	return resp.Result().(*GetJsapiTicketResp)
}

/*
GetAppJsapiTicket 获取应用的jsapi_ticket

https://developer.work.weixin.qq.com/document/path/90506
*/
func (j *Jsapi) GetAppJsapiTicket() *GetJsapiTicketResp {
	res := &GetJsapiTicketResp{}
	if j.Err != nil {
		return res
	}
	path := "/cgi-bin/ticket/get"
	params := map[string]string{
		"type":         "agent_config",
		"access_token": j.Wecom.GetAccessToken(),
	}
	resp, err := j.client.R().
		SetQueryParams(params).
		SetError(&GetJsapiTicketResp{}).
		SetResult(&GetJsapiTicketResp{}).
		Get(path)

	if err != nil {
		j.Err = err
		return res
	}
	logger.Debugf("Wecom Jsapi.GetAppJsapiTicket resp: %s", resp)
	// { "errcode":0, "errmsg":"ok", "ticket":"bxLdikRXVbTPdHSM05e5u5sUoXNKd8-41ZO3MhKoyN5OfkWITDGgnr2fwJ0m9E8NYzWKVZvdVtaUgWvsdshFKA", "expires_in":7200 }
	if resp.IsError() {
		return resp.Error().(*GetJsapiTicketResp)
	}
	return resp.Result().(*GetJsapiTicketResp)
}

/*
GetJsapiSign JS-SDK使用权限签名算法

https://developer.work.weixin.qq.com/document/path/90506
*/
func (j *Jsapi) GetJsapiSign(ticket, nonce, timestamp, url string) string {
	str := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%s&url=%s", ticket, nonce, timestamp, url)
	return dongle.Encrypt.FromString(str).BySha1().ToHexString()
}
