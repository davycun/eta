package wecom

import (
	"github.com/davycun/eta/pkg/common/logger"
)

type Oauth2 struct {
	Wecom
}

type GetAccessTokenResp struct {
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}
type GetUserInfoResp struct {
	Errcode        int    `json:"errcode"`
	Errmsg         string `json:"errmsg"`
	Userid         string `json:"userid"`          // 成员UserID。若需要获得用户详情信息，可调用通讯录接口：读取成员。如果是互联企业/企业互联/上下游，则返回的UserId格式如：CorpId/userid
	UserTicket     string `json:"user_ticket"`     // 成员票据，最大为512字节，有效期为1800s。scope为snsapi_privateinfo，且用户在应用可见范围之内时返回此参数。 后续利用该参数可以获取用户信息或敏感信息，参见"获取访问用户敏感信息"。暂时不支持上下游或/企业互联场景
	Openid         string `json:"openid"`          // 非企业成员的标识，对当前企业唯一。不超过64字节
	ExternalUserid string `json:"external_userid"` // 外部联系人id，当且仅当用户是企业的客户，且跟进人在应用的可见范围内时返回。如果是第三方应用调用，针对同一个客户，同一个服务商不同应用获取到的id相同
}

/*
GetAccessToken 获取access_token

https://developer.work.weixin.qq.com/document/path/91039
*/
func (o *Oauth2) GetAccessToken() *GetAccessTokenResp {
	res := &GetAccessTokenResp{}
	if o.Err != nil {
		return res
	}
	path := "/cgi-bin/gettoken"
	params := map[string]string{
		"corpid":     o.AppKey,
		"corpsecret": o.AppSecret,
	}
	resp, err := o.client.R().
		SetQueryParams(params).
		SetError(&GetAccessTokenResp{}).
		SetResult(&GetAccessTokenResp{}).
		Get(path)

	if err != nil {
		o.Err = err
		return res
	}
	logger.Debugf("Wecom Oauth2.GetAccessToken resp: %s", resp)
	// {"errcode":0,"errmsg":"ok","access_token":"accesstoken000001","expires_in":7200}
	if resp.IsError() {
		return resp.Error().(*GetAccessTokenResp)
	}
	return resp.Result().(*GetAccessTokenResp)
}

/*
GetUserInfo 获取访问用户身份

https://developer.work.weixin.qq.com/document/path/91023
*/
func (o *Oauth2) GetUserInfo(code string) *GetUserInfoResp {
	res := &GetUserInfoResp{}
	if o.Err != nil {
		return res
	}
	path := "/cgi-bin/auth/getuserinfo"
	params := map[string]string{
		"access_token": o.Wecom.GetAccessToken(),
		"code":         code,
	}

	resp, err := o.client.R().
		SetQueryParams(params).
		SetError(&GetUserInfoResp{}).
		SetResult(&GetUserInfoResp{}).
		Get(path)

	if err != nil {
		o.Err = err
		return res
	}
	logger.Debugf("GetUser resp: %s", resp)
	// {"errcode":0,"errmsg":"ok","openid":"OPENID","external_userid":"EXTERNAL_USERID"}
	// {"errcode":0,"errmsg":"ok","userid":"USERID","user_ticket":"USER_TICKET"}
	// {"errcode":40029,"errmsg":"invalid code"}
	if resp.IsError() {
		return resp.Error().(*GetUserInfoResp)
	}
	return resp.Result().(*GetUserInfoResp)
}
