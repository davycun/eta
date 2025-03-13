package wecom

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/go-resty/resty/v2"
	"github.com/patrickmn/go-cache"
	"time"
)

var (
	ApiCache = cache.New(5*time.Minute, 10*time.Minute)
)

type Wecom struct {
	BaseUrl         string `json:"base_url"`
	AppKey          string `json:"app_key"`
	AppSecret       string `json:"app_secret"`
	Err             error
	client          *resty.Client
	token           *Token
	corpJsapiTicket *JsapiTicket
	appJsapiTicket  *JsapiTicket
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

func NewWecom(baseUrl, appKey, appSecret string) *Wecom {
	if appKey == "" || appSecret == "" {
		logger.Errorf("创建Wecom, appKey 和 appSecret 不能为空")
		return nil
	}
	if a, found := ApiCache.Get(appKey); found {
		return a.(*Wecom)
	}
	if baseUrl == "" {
		baseUrl = "https://qyapi.weixin.qq.com"
	}
	w := &Wecom{
		BaseUrl:         baseUrl,
		AppKey:          appKey,
		AppSecret:       appSecret,
		client:          resty.New(), // 创建一个resty客户端
		token:           &Token{},
		corpJsapiTicket: &JsapiTicket{},
		appJsapiTicket:  &JsapiTicket{},
	}
	w.client.SetTimeout(1 * time.Minute)
	w.client.SetBaseURL(baseUrl)
	ApiCache.Set(appKey, w, cache.NoExpiration)
	return w
}

func (w *Wecom) Oauth2() *Oauth2 {
	return &Oauth2{
		Wecom: *w,
	}
}
func (w *Wecom) Jsapi() *Jsapi {
	return &Jsapi{
		Wecom: *w,
	}
}

func (t *Token) IsExpired() bool {
	return t.AccessToken == "" || t.ExpiredAt.Second()-time.Now().Second() < 300
}

func (t *JsapiTicket) IsExpired() bool {
	return t.Ticket == "" || t.ExpiredAt.Second()-time.Now().Second() < 300
}

func (w *Wecom) GetAccessToken() string {
	if w.token.IsExpired() {
		w.RefreshAccessToken()
	}
	return w.token.AccessToken
}

func (w *Wecom) RefreshAccessToken() *Wecom {
	logger.Info("Wecom 刷新 AccessToken")
	accessToken := w.Oauth2().GetAccessToken()
	if w.Err != nil {
		logger.Errorf("Wecom 刷新 AccessToken 报错，异常: %s ", w.Err)
		return w
	}
	w.token.AccessToken = accessToken.AccessToken
	w.token.ExpiresIn = accessToken.ExpiresIn
	w.token.ExpiredAt = time.Now().Add(time.Duration(accessToken.ExpiresIn) * time.Second)
	return w
}

func (w *Wecom) GetCorpJsapiTicket() string {
	if w.corpJsapiTicket.IsExpired() {
		w.RefreshCorpJsapiTicket()
	}
	return w.corpJsapiTicket.Ticket
}

func (w *Wecom) RefreshCorpJsapiTicket() *Wecom {
	logger.Info("Wecom 刷新 corp jsapiTicket")
	jsapiTicket := w.Jsapi().GetCorpJsapiTicket()
	if w.Err != nil {
		logger.Errorf("Wecom 刷新 corp jsapiTicket 报错，异常: %s ", w.Err)
		return w
	}
	w.corpJsapiTicket.Ticket = jsapiTicket.Ticket
	w.corpJsapiTicket.ExpiresIn = jsapiTicket.ExpiresIn
	w.corpJsapiTicket.ExpiredAt = time.Now().Add(time.Duration(jsapiTicket.ExpiresIn) * time.Second)
	return w
}

func (w *Wecom) GetAppJsapiTicket() string {
	if w.appJsapiTicket.IsExpired() {
		w.RefreshAppJsapiTicket()
	}
	return w.appJsapiTicket.Ticket
}

func (w *Wecom) RefreshAppJsapiTicket() *Wecom {
	logger.Info("Wecom 刷新 app jsapiTicket")
	jsapiTicket := w.Jsapi().GetAppJsapiTicket()
	if w.Err != nil {
		logger.Errorf("Wecom 刷新 app jsapiTicket 报错，异常: %s ", w.Err)
		return w
	}
	w.appJsapiTicket.Ticket = jsapiTicket.Ticket
	w.appJsapiTicket.ExpiresIn = jsapiTicket.ExpiresIn
	w.appJsapiTicket.ExpiredAt = time.Now().Add(time.Duration(jsapiTicket.ExpiresIn) * time.Second)
	return w
}

func (w *Wecom) url(path string) string {
	return fmt.Sprintf("%s%s", w.BaseUrl, path)
}
