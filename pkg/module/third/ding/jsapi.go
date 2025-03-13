package ding

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/golang-module/dongle"
)

type Jsapi struct {
	Ding
}

type GetJsapiTicketResp struct {
	Ticket    string `json:"jsapiTicket"`
	ExpiresIn int    `json:"expireIn"`
}

/*
GetJsapiTicket 获取jsapiTicket

https://open.dingtalk.com/document/orgapp/create-a-jsapi-ticket
*/
func (j *Jsapi) GetJsapiTicket() *GetJsapiTicketResp {
	res := &GetJsapiTicketResp{}
	if j.Err != nil {
		return res
	}
	path := "/v1.0/oauth2/jsapiTickets"
	resp, err := j.client.R().
		SetHeader("x-acs-dingtalk-access-token", j.Ding.GetAccessToken()).
		SetError(&GetJsapiTicketResp{}).
		SetResult(&GetJsapiTicketResp{}).
		Get(path)

	if err != nil {
		j.Err = err
		return res
	}
	logger.Debugf("Ding Jsapi.GetJsapiTicket resp: %s", resp)
	// { "errcode":0, "errmsg":"ok", "ticket":"bxLdikRXVbTPdHSM05e5u5sUoXNKd8-41ZO3MhKoyN5OfkWITDGgnr2fwJ0m9E8NYzWKVZvdVtaUgWvsdshFKA", "expires_in":7200 }
	if resp.IsError() {
		return resp.Error().(*GetJsapiTicketResp)
	}
	return resp.Result().(*GetJsapiTicketResp)
}

/*
GetJsapiSign JS-SDK使用权限签名算法

https://open.dingtalk.com/document/orgapp/develop-webapp-frontend
*/
func (j *Jsapi) GetJsapiSign(nonce, timestamp, url string) string {
	str := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%s&url=%s", j.Ding.GetJsapiTicket(), nonce, timestamp, url)
	return dongle.Encrypt.FromString(str).BySha1().ToHexString()
}
