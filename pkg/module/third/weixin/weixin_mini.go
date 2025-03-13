package weixin

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/httpclient"
	"net/http"
)

type Code2SessionResponse struct {
	WxError
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
	OpenId     string `json:"open_id"`
}
type GetPhoneNumberResponse struct {
	WxError
	PhoneInfo struct {
		PhoneNumber     string `json:"phoneNumber"`
		PurePhoneNumber string `json:"purePhoneNumber"`
		CountryCode     string `json:"countryCode"`
		Watermark       struct {
			Timestamp int64  `json:"timestamp"`
			Appid     string `json:"appid"`
		} `json:"watermark"`
	} `json:"phone_info"`
}

func (w *WeiXin) Code2Session(code string) (Code2SessionResponse, error) {

	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", w.wxk.AppId, w.wxk.AppSecret, code)
	var rs Code2SessionResponse
	err := httpclient.DefaultHttpClient.Url(url).
		Method(http.MethodGet).Do(&rs).Error
	return rs, err
}
func (w *WeiXin) GetPhoneNumber(code string) (GetPhoneNumberResponse, error) {

	var req = struct {
		Code string `json:"code"`
	}{
		Code: code,
	}
	url := fmt.Sprintf("https://api.weixin.qq.com/wxa/business/getuserphonenumber?access_token=%s", w.wxk.AppId)
	var rs GetPhoneNumberResponse
	err := httpclient.DefaultHttpClient.Url(url).
		Method(http.MethodPost).
		Body(httpclient.MIMEJSON, req).
		Do(&rs).
		Error
	return rs, err
}
