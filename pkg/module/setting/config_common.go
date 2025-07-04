package setting

import (
	"github.com/davycun/eta/pkg/common/logger"
	"gorm.io/gorm"
	"math/rand"
)

type BaseCredentials struct {
	BaseUrl            string                 `json:"base_url,omitempty"`       //基础URL
	ExtraBaseUrl       []string               `json:"extra_base_url,omitempty"` //额外可用的BaseUrl
	AppKey             string                 `json:"app_key,omitempty"`
	AppSecret          string                 `json:"app_secret,omitempty"`           //
	ProxyUrl           string                 `json:"proxy_url,omitempty"`            // 代理地址
	InsecureSkipVerify bool                   `json:"insecure_skip_verify,omitempty"` //是否跳过https认证
	UserName           string                 `json:"user_name,omitempty"`            //登录用户名
	Password           string                 `json:"password,omitempty"`             //登录密码
	Headers            map[string]string      `json:"headers"`                        //固定的一些headers
	Debug              bool                   `json:"debug,omitempty"`                //resty是否开启debug模式
	Timeout            int                    `json:"timeout,omitempty"`              //访问请求超时时间，单位是秒
	Extra              map[string]interface{} `json:"extra,omitempty"`                //额外的一些配置
}

func (b BaseCredentials) Valid() bool {
	return b.BaseUrl != "" && ((b.AppKey != "" && b.AppSecret != "") || (b.UserName != "" && b.Password != ""))
}
func (b BaseCredentials) RandomBaseUrl() string {
	if len(b.ExtraBaseUrl) < 1 {
		return b.BaseUrl
	}
	us := append([]string{b.BaseUrl}, b.ExtraBaseUrl...)
	return us[rand.Intn(len(us))]
}

// CommonConfig
// 前后端共用的一些配置，有可能不一定需要登录就能获取，所以针对这个配置可能需要单独接口处理
type CommonConfig struct {
	SmsNeedImageCaptcha bool `json:"sms_need_image_captcha"` //发送短信验证码的时候是否需要图片验证码，避免短信发送接口被攻击
}

func GetCommonConfig(db *gorm.DB) CommonConfig {
	cfg, err := GetConfig[CommonConfig](db, ConfigCommonCategory, ConfigCommonName)
	if err != nil {
		logger.Errorf("load common config err %s", err)
	}
	return cfg
}
