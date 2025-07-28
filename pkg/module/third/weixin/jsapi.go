package weixin

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/dromara/dongle"
)

type Jsapi struct {
	WeiXin
}

type GetJsapiTicketResp struct {
	Errcode   int    `json:"errcode"`
	Errmsg    string `json:"errmsg"`
	Ticket    string `json:"ticket"`
	ExpiresIn int    `json:"expires_in"`
}

/*
GetJsapiTicket 获取应用的jsapi_ticket

https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/JS-SDK.html#62
*/
func (j *Jsapi) GetJsapiTicket() (res *GetJsapiTicketResp, err error) {
	path := "https://api.weixin.qq.com/cgi-bin/ticket/getticket"
	params := map[string]string{
		"access_token": j.WeiXin.GetAccessToken(),
		"type":         "jsapi",
	}

	resp, err := j.client.R().
		SetQueryParams(params).
		SetError(&GetJsapiTicketResp{}).
		SetResult(&GetJsapiTicketResp{}).
		Get(path)

	if err != nil {
		return
	}
	logger.Debugf("WeiXin Jsapi.GetJsapiTicket resp: %s", resp)
	// {"errcode":0,"errmsg":"ok","ticket":"bxLdikRXVbTPdHSM05e5u5sUoXNKd8-41ZO3MhKoyN5OfkWITDGgnr2fwJ0m9E8NYzWKVZvdVtaUgWvsdshFKA","expires_in":7200}
	if resp.IsError() {
		return resp.Error().(*GetJsapiTicketResp), err
	}
	return resp.Result().(*GetJsapiTicketResp), err
}

/*
GetJsapiSign JS-SDK使用权限签名算法

https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/JS-SDK.html#62
*/
func (j *Jsapi) GetJsapiSign(nonce, timestamp, url string) string {
	ticket := j.WeiXin.GetJsapiTicket()
	if ticket == "" {
		logger.Errorf("WeiXin Jsapi.GetJsapiSign ticket is empty")
		return ""
	}
	str := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%s&url=%s", ticket, nonce, timestamp, url)
	return dongle.Encrypt.FromString(str).BySha1().ToHexString()
}
