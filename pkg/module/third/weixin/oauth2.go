package weixin

import (
	"github.com/davycun/eta/pkg/common/logger"
)

type Oauth2 struct {
	WeiXin
}

type AccessTokenResp struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`

	AccessToken    string `json:"access_token"`
	ExpiresIn      int    `json:"expires_in"`
	RefreshToken   string `json:"refresh_token"`
	Openid         string `json:"openid"`
	Scope          string `json:"scope"`
	IsSnapshotuser int    `json:"is_snapshotuser"`
	Unionid        string `json:"unionid"`
}

/*
AccessToken 通过code换取网页授权access_token

https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/Wechat_webpage_authorization.html
*/
func (o *Oauth2) AccessToken(code string) (res *AccessTokenResp, err error) {
	path := "https://api.weixin.qq.com/sns/oauth2/access_token"
	params := map[string]string{
		"appid":      o.wxk.AppId,
		"secret":     o.wxk.AppSecret,
		"code":       code,
		"grant_type": "authorization_code",
	}

	resp, err := o.client.R().
		SetQueryParams(params).
		SetError(&AccessTokenResp{}).
		SetResult(&AccessTokenResp{}).
		Post(path)

	if err != nil {
		return
	}
	logger.Debugf("WeiXin Oauth2.GetAccessToken resp: %s", resp)
	// {"access_token":"ACCESS_TOKEN","expires_in":7200,"refresh_token":"REFRESH_TOKEN","openid":"OPENID","scope":"SCOPE","is_snapshotuser":1,"unionid":"UNIONID"}
	if resp.IsError() {
		return resp.Error().(*AccessTokenResp), err
	}
	return resp.Result().(*AccessTokenResp), err
}
