package zzd

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/go-resty/resty/v2"
	"github.com/golang-module/dongle"
	"github.com/patrickmn/go-cache"
	"math/rand"
	"net/url"
	"sort"
	"strings"
	"time"
)

var (
	ApiCache = cache.New(5*time.Minute, 10*time.Minute)
)

type Zzd struct {
	BaseUrl     string `json:"base_url"`
	AppKey      string `json:"app_key"`
	AppSecret   string `json:"app_secret"`
	TenantId    string `json:"tenant_id"`
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

func NewZzd(baseUrl, appKey, appSecret, tenantId string) *Zzd {
	if appKey == "" || appSecret == "" {
		logger.Errorf("创建Zzd, appKey 和 appSecret 不能为空")
		return nil
	}
	if a, found := ApiCache.Get(appKey); found {
		return a.(*Zzd)
	}
	if baseUrl == "" {
		baseUrl = "https://openplatform-pro.ding.zj.gov.cn"
	}
	z := &Zzd{
		BaseUrl:     baseUrl,
		AppKey:      appKey,
		AppSecret:   appSecret,
		TenantId:    tenantId,
		client:      resty.New(), // 创建一个resty客户端
		token:       &Token{},
		jsapiTicket: &JsapiTicket{},
	}
	z.client.SetTimeout(1 * time.Minute)
	z.client.SetBaseURL(baseUrl)
	ApiCache.Set(appKey, z, cache.NoExpiration)
	return z
}

func (z *Zzd) Emp() *Emp {
	return &Emp{
		Zzd: *z,
	}
}
func (z *Zzd) Message() *Message {
	return &Message{
		Zzd: *z,
	}
}
func (z *Zzd) Oauth2() *Oauth2 {
	return &Oauth2{
		Zzd: *z,
	}
}
func (z *Zzd) Org() *Org {
	return &Org{
		Zzd: *z,
	}
}
func (z *Zzd) Jsapi() *Jsapi {
	return &Jsapi{
		Zzd: *z,
	}
}

func (t *Token) IsExpired() bool {
	return t.AccessToken == "" || t.ExpiredAt.Second()-time.Now().Second() < 300
}

func (t *JsapiTicket) IsExpired() bool {
	return t.Ticket == "" || t.ExpiredAt.Second()-time.Now().Second() < 300
}

func (z *Zzd) GetAccessToken() string {
	if z.token.IsExpired() {
		z.RefreshAccessToken()
	}
	return z.token.AccessToken
}

func (z *Zzd) RefreshAccessToken() *Zzd {
	logger.Info("Wecom 刷新 AccessToken")
	accessToken := z.Oauth2().GetAccessToken()
	if z.Err != nil || !accessToken.Success || !accessToken.Content.Success {
		logger.Errorf("zzd 刷新 AccessToken 报错，异常: %s ", z.Err)
		return z
	}
	z.token.AccessToken = accessToken.Content.Data.AccessToken
	z.token.ExpiresIn = accessToken.Content.Data.ExpiresIn
	z.token.ExpiredAt = time.Now().Add(time.Duration(accessToken.Content.Data.ExpiresIn) * time.Second)
	return z
}

func (z *Zzd) signature(method string, path string, params map[string]interface{}) (header map[string]string, query url.Values) {
	rand.Seed(time.Now().UnixNano())
	var (
		timestamp = time.Now().Format("2006-01-02T15:04:05.000+08:00")
		nonce     = fmt.Sprintf("%d%d", time.Now().UnixMilli(), rand.Intn(9000)+1000)
	)
	method = strings.ToUpper(method)

	u := url.URL{}
	u.Query()
	bQuery := buildQuery(params)
	q, err := url.QueryUnescape(bQuery)
	if err != nil {
		z.Err = err
		logger.Errorf("签名失败, 错误信息: %s", err.Error())
		return
	}

	var toSignStr string
	if q == "" {
		toSignStr = fmt.Sprintf("%s\n%s\n%s\n%s", method, timestamp, nonce, path)
	} else {
		toSignStr = fmt.Sprintf("%s\n%s\n%s\n%s\n%s", method, timestamp, nonce, path, q)
	}

	sign := dongle.Encrypt.FromString(toSignStr).ByHmacSha256(z.AppSecret).ToBase64String()

	header = map[string]string{
		"X-Hmac-Auth-Timestamp": timestamp,
		"X-Hmac-Auth-Version":   "1.0",
		"X-Hmac-Auth-Nonce":     nonce,
		"apiKey":                z.AppKey,
		"X-Hmac-Auth-Signature": sign,
	}
	query, err = url.ParseQuery(bQuery)
	if err != nil {
		z.Err = err
	}
	return
}

func (z *Zzd) GetJsapiTicket() string {
	if z.jsapiTicket.IsExpired() {
		z.RefreshJsapiTicket()
	}
	return z.jsapiTicket.Ticket
}

func (z *Zzd) RefreshJsapiTicket() *Zzd {
	logger.Info("Zzd 刷新 JsapiTicket")
	jsapiTicket := z.Jsapi().GetJsapiTicket()
	if z.Err != nil || !jsapiTicket.Success || !jsapiTicket.Content.Success {
		logger.Errorf("Zzd 刷新 JsapiTicket 报错，异常: %s ", z.Err)
		return z
	}
	z.jsapiTicket.Ticket = jsapiTicket.Content.Data.AccessToken
	z.jsapiTicket.ExpiresIn = jsapiTicket.Content.Data.ExpiresIn
	z.jsapiTicket.ExpiredAt = time.Now().Add(time.Duration(jsapiTicket.Content.Data.ExpiresIn) * time.Second)
	return z
}

func (z *Zzd) url(path string) string {
	return fmt.Sprintf("%s%s", z.BaseUrl, path)
}

func buildQuery(params map[string]interface{}) string {
	query := make([]string, 0)
	keys := make([]string, 0)
	for k, _ := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := params[k]
		switch value := v.(type) {
		case []string:
			vArr := v.([]string)
			sort.Strings(vArr)
			for _, item := range vArr {
				query = append(query, fmt.Sprintf("%s=%s", k, item))
			}
		case string:
			query = append(query, fmt.Sprintf("%s=%s", k, url.QueryEscape(value)))
		case bool:
			query = append(query, fmt.Sprintf("%s=%t", k, v))
		case int, int8, int16, int32, int64:
			query = append(query, fmt.Sprintf("%s=%d", k, v))
		default:
			query = append(query, fmt.Sprintf("%s=%s", k, v))
		}
	}
	return strings.Join(query, "&")
}

// 使用键值对构建map，并移除value为空的元素
func buildParam(data map[string]interface{}) map[string]interface{} {
	n := make(map[string]interface{}, len(data))
	for key, value := range data {
		if value != nil && value != "" {
			n[key] = value
		}
	}
	return n
}
