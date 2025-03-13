package weixin

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/httpclient"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/go-resty/resty/v2"
	"github.com/patrickmn/go-cache"
	"net/http"
	"time"
)

type WxToken struct {
	AccessToken string    `json:"access_token"`
	ExpiresIn   int       `json:"expires_in"`
	ExpiredAt   time.Time `json:"expired_at"`
}

func (w *WxToken) IsExpired() bool {
	return w.AccessToken == "" || w.ExpiredAt.Second()-time.Now().Second() < 300
}

func (w *WxToken) ReLoad(key, secret string) {
	logger.Infof("开始为%s重新加载AccessToken", key)
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", key, secret)
	var tk WxToken
	err := httpclient.DefaultHttpClient.Url(url).Method(http.MethodGet).Do(&tk).Error
	if err != nil {
		logger.Errorf("invoke weixin accesstoken error %s ", err.Error())
	} else {
		w.ExpiredAt = time.Now().Add(time.Duration(tk.ExpiresIn) * time.Second)
		w.ExpiresIn = tk.ExpiresIn
		w.AccessToken = tk.AccessToken
	}
}

type JsapiTicket struct {
	Ticket    string    `json:"ticket"`
	ExpiresIn int       `json:"expires_in"`
	ExpiredAt time.Time `json:"expired_at"`
}

func (j *JsapiTicket) IsExpired() bool {
	return j.Ticket == "" || j.ExpiredAt.Second()-time.Now().Second() < 300
}

type WeiXin struct {
	wxk         *WxKey
	token       *WxToken
	jsapiTicket *JsapiTicket
	client      *resty.Client
}

func (w *WeiXin) GetAccessToken() string {
	if w.token.IsExpired() {
		w.token.ReLoad(w.wxk.AppId, w.wxk.AppSecret)
	}
	return w.token.AccessToken
}

func (w *WeiXin) GetJsapiTicket() string {
	if w.jsapiTicket.IsExpired() {
		w.RefreshJsapiTicket()
	}
	return w.jsapiTicket.Ticket
}

func GetWeiXin(key string) *WeiXin {
	if key == "" {
		logger.Errorf("指定key[%s]的weixin不存在", key)
		return nil
	}
	if a, found := ApiCache.Get(key); found {
		return a.(*WeiXin)
	}
	if wk, ok := WK[key]; ok {
		return NewWeiXin(wk)
	}
	return nil
}

func NewWeiXin(wk *WxKey) *WeiXin {
	if wk.Key == "" {
		logger.Errorf("WxKey 的Key不能为空")
		return nil
	}
	if a, found := ApiCache.Get(wk.Key); found {
		return a.(*WeiXin)
	}
	if wk.AppId == "" || wk.AppSecret == "" {
		logger.Errorf("创建WeiXin[%s], appId和Secret不能为空", wk.Key)
		return nil
	}
	wx := &WeiXin{
		wxk:         wk,
		token:       &WxToken{},
		jsapiTicket: &JsapiTicket{},
		client:      resty.New(), // 创建一个resty客户端
	}
	wx.client.SetTimeout(1 * time.Minute)
	ApiCache.Set(wk.Key, wx, cache.NoExpiration)
	return wx
}

func (w *WeiXin) Oauth2() *Oauth2 {
	return &Oauth2{
		WeiXin: *w,
	}
}

func (w *WeiXin) Jsapi() *Jsapi {
	return &Jsapi{
		WeiXin: *w,
	}
}

func (w *WeiXin) RefreshJsapiTicket() *WeiXin {
	logger.Info("WeiXin 刷新 JsapiTicket")
	jsapiTicket, err := w.Jsapi().GetJsapiTicket()
	if err != nil {
		logger.Errorf("Wecom 刷新 JsapiTicket 报错，异常: %s ", err)
		return w
	}
	w.jsapiTicket.Ticket = jsapiTicket.Ticket
	w.jsapiTicket.ExpiresIn = jsapiTicket.ExpiresIn
	w.jsapiTicket.ExpiredAt = time.Now().Add(time.Duration(jsapiTicket.ExpiresIn) * time.Second)
	return w
}
