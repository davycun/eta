package zzd

import (
	"github.com/davycun/eta/pkg/common/logger"
	"net/http"
)

type Oauth2 struct {
	Zzd
}

type GetAccessTokenResp struct {
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
		BizErrorCode    string `json:"bizErrorCode"`
	} `json:"content"`
	BizErrorCode string `json:"bizErrorCode"`
}
type AuthCodeResp struct {
	Success bool `json:"success"`
	Content struct {
		Success         bool   `json:"success"`
		ResponseMessage string `json:"responseMessage"`
		ResponseCode    string `json:"responseCode"`
		BizErrorCode    string `json:"bizErrorCode"`
		Data            struct {
			LastName     string `json:"lastName"`
			RealmId      int    `json:"realmId"`
			ClientId     string `json:"clientId"`
			Openid       string `json:"openid"`
			RealmName    string `json:"realmName"`
			NickNameCn   string `json:"nickNameCn"`
			TenantUserId string `json:"tenantUserId"`
			Avatar       string `json:"avatar"`
			EmployeeCode string `json:"employeeCode"`
			AccountId    int    `json:"accountId"`
			TenantName   string `json:"tenantName"`
			ReferId      string `json:"referId"`
			Namespace    string `json:"namespace"`
			TenantId     int    `json:"tenantId"`
			Account      string `json:"account"`
		} `json:"data"`
	} `json:"content"`
}
type QrCodeResp struct {
	Success bool `json:"success"`
	Content struct {
		Success         bool   `json:"success"`
		ResponseMessage string `json:"responseMessage"`
		ResponseCode    string `json:"responseCode"`
		BizErrorCode    string `json:"bizErrorCode"`
		Data            struct {
			LastName     string `json:"lastName"`
			RealmId      int    `json:"realmId"`
			ClientId     string `json:"clientId"`
			RealmName    string `json:"realmName"`
			NickNameCn   string `json:"nickNameCn"`
			TenantUserId string `json:"tenantUserId"`
			Avatar       string `json:"avatar"`
			EmployeeCode string `json:"employeeCode"`
			AccountId    int    `json:"accountId"`
			TenantName   string `json:"tenantName"`
			Namespace    string `json:"namespace"`
			TenantId     int    `json:"tenantId"`
			Account      string `json:"account"`
		} `json:"data"`
	} `json:"content"`
}

func (o *Oauth2) GetAccessToken() *GetAccessTokenResp {
	res := &GetAccessTokenResp{}
	if o.Err != nil {
		return res
	}
	path := "/gettoken.json"
	params := map[string]interface{}{}
	header, query := o.signature(http.MethodGet, path, params)

	resp, err := o.client.R().
		SetHeaders(header).
		SetQueryParamsFromValues(query).
		SetError(&GetAccessTokenResp{}).
		SetResult(&GetAccessTokenResp{}).
		Get(path)

	if err != nil {
		o.Err = err
		return res
	}
	logger.Debugf("Oauth2.GetAccessToken resp: %s", resp)
	// {"success":true,"content":{"data":{"expiresIn":7200,"accessToken":"app_298e692864a640feb529814c4b20f4a1"},"success":true,"requestId":"0aa0734417036651733746211d0011","responseMessage":"OK","responseCode":"0","bizErrorCode":"0"},"bizErrorCode":"0"}
	if resp.IsError() {
		return resp.Error().(*GetAccessTokenResp)
	}
	return resp.Result().(*GetAccessTokenResp)
}

// GetUserInfo 免密登录, 服务端通过authCode获取授权用户的个人信息
func (o *Oauth2) GetUserInfo(authCode string) *AuthCodeResp {
	res := &AuthCodeResp{}
	if o.Err != nil {
		return res
	}
	path := "/rpc/oauth2/dingtalk_app_user.json"
	params := buildParam(map[string]interface{}{
		"access_token": o.Zzd.GetAccessToken(),
		"auth_code":    authCode,
	})
	header, query := o.signature(http.MethodPost, path, params)

	resp, err := o.client.R().
		SetHeaders(header).
		SetFormDataFromValues(query).
		SetError(&AuthCodeResp{}).
		SetResult(&AuthCodeResp{}).
		Post(path)

	if err != nil {
		o.Err = err
		return res
	}
	logger.Debugf("Zzd Oauth2.GetUserInfo resp: %s", resp)
	if resp.IsError() {
		return resp.Error().(*AuthCodeResp)
	}
	return resp.Result().(*AuthCodeResp)
}

// GetUserInfoByCode 扫码登录, 服务端通过临时授权码获取授权用户的个人信息
func (o *Oauth2) GetUserInfoByCode(code string) *QrCodeResp {
	res := &QrCodeResp{}
	if o.Err != nil {
		return res
	}
	path := "/rpc/oauth2/getuserinfo_bycode.json"
	params := buildParam(map[string]interface{}{
		"access_token": o.Zzd.GetAccessToken(),
		"code":         code,
	})
	header, query := o.signature(http.MethodPost, path, params)

	resp, err := o.client.R().
		SetHeaders(header).
		SetFormDataFromValues(query).
		SetError(&QrCodeResp{}).
		SetResult(&QrCodeResp{}).
		Post(path)

	if err != nil {
		o.Err = err
		return res
	}
	logger.Debugf("Zzd Oauth2.GetUserInfoByCode resp: %s", resp)
	if resp.IsError() {
		return resp.Error().(*QrCodeResp)
	}
	return resp.Result().(*QrCodeResp)
}
