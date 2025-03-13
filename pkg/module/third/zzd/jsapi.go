package zzd

import (
	"github.com/davycun/eta/pkg/common/logger"
	"net/http"
)

type Jsapi struct {
	Zzd
}

type GetJsapiTicketResp struct {
	Success bool `json:"success"`
	Content struct {
		Data struct {
			ExpiresIn   int    `json:"expiresIn"`
			AccessToken string `json:"accessToken"`
		} `json:"data"`
		Success         bool   `json:"success"`
		RequestId       string `json:"requestId"`
		ResponseMessage string `json:"responseMessage"`
		ResponseCode    string `json:"responseCode"`
	} `json:"content"`
}

/*
GetJsapiTicket JSAPI鉴权接口

https://openplatform-portal.dg-work.cn/portal/#/helpdoc?apiType=serverapi&docKey=2674861
*/
func (j *Jsapi) GetJsapiTicket() *GetJsapiTicketResp {
	res := &GetJsapiTicketResp{}
	if j.Err != nil {
		return res
	}
	path := "/get_jsapi_token.json"
	params := buildParam(map[string]interface{}{
		"access_token": j.Zzd.GetAccessToken(),
	})
	header, query := j.signature(http.MethodPost, path, params)

	resp, err := j.client.R().
		SetHeaders(header).
		SetFormDataFromValues(query).
		SetError(&GetJsapiTicketResp{}).
		SetResult(&GetJsapiTicketResp{}).
		Post(path)

	if err != nil {
		j.Err = err
		return res
	}
	logger.Debugf("Zzd Jsapi.GetJsapiTicket resp: %s", resp)
	// {"success":true,"content":{"data":{"expiresIn":4561,"accessToken":"jsApi_009dde783756450184e897b9e2943c71"},"success":true,"requestId":"df04428415763320830387621d291f","responseMessage":"OK","responseCode":"0"}}
	// {"success":true,"content":{"success":false,"requestId":"ac16084b16836764937908989d0011","responseMessage":"accessToken not correct","responseCode":"250027","bizErrorCode":"250027"},"bizErrorCode":"250027"}
	if resp.IsError() {
		return resp.Error().(*GetJsapiTicketResp)
	}
	return resp.Result().(*GetJsapiTicketResp)
}
