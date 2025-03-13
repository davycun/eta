package setting

import (
	"github.com/davycun/eta/pkg/common/logger"
	"gorm.io/gorm"
)

type BaseCredentials struct {
	BaseUrl            string                 `json:"base_url,omitempty"`
	AppKey             string                 `json:"app_key,omitempty"`
	AppSecret          string                 `json:"app_secret,omitempty"`
	ProxyUrl           string                 `json:"proxy_url,omitempty"`
	InsecureSkipVerify bool                   `json:"insecure_skip_verify,omitempty"`
	Debug              bool                   `json:"debug,omitempty"`
	UserName           string                 `json:"user_name,omitempty"`
	Password           string                 `json:"password,omitempty"`
	Extra              map[string]interface{} `json:"extra,omitempty"`
}

func (b BaseCredentials) Valid() bool {
	return b.BaseUrl != "" && ((b.AppKey != "" && b.AppSecret != "") || (b.UserName != "" && b.Password != ""))
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
