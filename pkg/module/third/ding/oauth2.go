package ding

import (
	"github.com/davycun/eta/pkg/common/logger"
)

type Oauth2 struct {
	Ding
}

type AccessTokenResp struct {
	AccessToken string `json:"accessToken"`
	ExpireIn    int    `json:"expireIn"`
}
type UserAccessTokenResp struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpireIn     int    `json:"expireIn"`
	CorpId       string `json:"corpId"`
}

/*
AccessToken 获取企业内部应用的accessToken

https://open.dingtalk.com/document/orgapp/obtain-the-access_token-of-an-internal-app
*/
func (o *Oauth2) AccessToken() *AccessTokenResp {
	res := &AccessTokenResp{}
	if o.Err != nil {
		return res
	}
	path := "/v1.0/oauth2/accessToken"
	params := map[string]string{
		"appKey":    o.AppKey,
		"appSecret": o.AppSecret,
	}
	resp, err := o.client.R().
		SetBody(params).
		SetHeader("Content-Type", "application/json").
		SetError(&AccessTokenResp{}).
		SetResult(&AccessTokenResp{}).
		Post(path)

	if err != nil {
		o.Err = err
		return res
	}
	logger.Debugf("Ding Oauth2.AccessToken resp: %s", resp)
	// {"accessToken":"fw8ef8we8f76e6f7s8dxxxx","expireIn":7200}
	if resp.IsError() {
		return resp.Error().(*AccessTokenResp)
	}
	return resp.Result().(*AccessTokenResp)
}

/*
UserAccessToken 获取用户token

https://open.dingtalk.com/document/orgapp/obtain-user-token
*/
func (o *Oauth2) UserAccessToken(code string) *UserAccessTokenResp {
	res := &UserAccessTokenResp{}
	if o.Err != nil {
		return res
	}
	path := "/v1.0/oauth2/userAccessToken"
	params := map[string]string{
		"clientId":     o.AppKey,
		"clientSecret": o.AppSecret,
		"code":         code,
		//"refreshToken": "abcd",
		"grantType": "authorization_code",
	}
	resp, err := o.client.R().
		SetBody(params).
		SetHeader("Content-Type", "application/json").
		SetError(&UserAccessTokenResp{}).
		SetResult(&UserAccessTokenResp{}).
		Post(path)

	if err != nil {
		o.Err = err
		return res
	}
	logger.Debugf("Ding Oauth2.UserAccessToken resp: %s", resp)
	// {"accessToken":"fw8ef8we8f76e6f7s8dxxxx","expireIn":7200}
	if resp.IsError() {
		return resp.Error().(*UserAccessTokenResp)
	}
	return resp.Result().(*UserAccessTokenResp)
}
