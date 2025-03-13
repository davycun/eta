package ding

import (
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/go-resty/resty/v2"
	"github.com/patrickmn/go-cache"
	"time"
)

var (
	ApiCache = cache.New(5*time.Minute, 10*time.Minute)
)

type Ding struct {
	BaseUrl     string `json:"base_url"`
	AppKey      string `json:"app_key"`
	AppSecret   string `json:"app_secret"`
	Err         error
	client      *resty.Client
	token       *Token
	jsapiTicket *JsapiTicket
}
type Token struct {
	AccessToken string    `json:"access_token"`
	ExpiresIn   int       `json:"expires_in"`
	ExpiredAt   time.Time `json:"expired_at"`
}
type JsapiTicket struct {
	Ticket    string    `json:"ticket"`
	ExpiresIn int       `json:"expires_in"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewDing(baseUrl, appKey, appSecret string) *Ding {
	if appKey == "" || appSecret == "" {
		logger.Errorf("创建Ding, appKey 和 appSecret 不能为空")
		return nil
	}
	if a, found := ApiCache.Get(appKey); found {
		return a.(*Ding)
	}

	if baseUrl == "" {
		baseUrl = "https://oapi.dingtalk.com"
	}
	w := &Ding{
		BaseUrl:     baseUrl,
		AppKey:      appKey,
		AppSecret:   appSecret,
		client:      resty.New(), // 创建一个resty客户端
		token:       &Token{},
		jsapiTicket: &JsapiTicket{},
	}
	w.client.SetTimeout(1 * time.Minute)
	w.client.SetBaseURL(baseUrl)
	ApiCache.Set(appKey, w, cache.NoExpiration)
	return w
}

func (d *Ding) Oauth2() *Oauth2 {
	return &Oauth2{
		Ding: *d,
	}
}
func (d *Ding) Jsapi() *Jsapi {
	return &Jsapi{
		Ding: *d,
	}
}
func (d *Ding) User() *User {
	return &User{
		Ding: *d,
	}
}

func (t *Token) IsExpired() bool {
	return t.AccessToken == "" || t.ExpiredAt.Second()-time.Now().Second() < 300
}

func (t *JsapiTicket) IsExpired() bool {
	return t.Ticket == "" || t.ExpiredAt.Second()-time.Now().Second() < 300
}

func (d *Ding) GetAccessToken() string {
	if d.token.IsExpired() {
		d.RefreshAccessToken()
	}
	return d.token.AccessToken
}

func (d *Ding) RefreshAccessToken() *Ding {
	logger.Info("Ding 刷新 AccessToken")
	accessToken := d.Oauth2().AccessToken()
	if d.Err != nil {
		logger.Errorf("Ding 刷新 AccessToken 报错，异常: %s ", d.Err)
		return d
	}
	d.token.AccessToken = accessToken.AccessToken
	d.token.ExpiresIn = accessToken.ExpireIn
	d.token.ExpiredAt = time.Now().Add(time.Duration(accessToken.ExpireIn) * time.Second)
	return d
}

func (d *Ding) GetJsapiTicket() string {
	if d.jsapiTicket.IsExpired() {
		d.RefreshJsapiTicket()
	}
	return d.jsapiTicket.Ticket
}

func (d *Ding) RefreshJsapiTicket() *Ding {
	logger.Info("Ding 刷新 JsapiTicket")
	jsapiTicket := d.Jsapi().GetJsapiTicket()
	if d.Err != nil {
		logger.Errorf("Ding 刷新 JsapiTicket 报错，异常: %s ", d.Err)
		return d
	}
	d.jsapiTicket.Ticket = jsapiTicket.Ticket
	d.jsapiTicket.ExpiresIn = jsapiTicket.ExpiresIn
	d.jsapiTicket.ExpiredAt = time.Now().Add(time.Duration(jsapiTicket.ExpiresIn) * time.Second)
	return d
}
