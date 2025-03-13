package setting

import (
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/logger"
	"gorm.io/gorm"
)

//配置示例
//  st := Setting{}
//	st.Namespace = constants.NamespaceEta
//	st.Category = ConfigSmsCategory
//	st.Name = ConfigSmsCategory
//	st.Content = ctype.NewJson(
//		SmsConfig{
//			Vendor: "aliyun",
//			SmsInfoMap: map[string]SmsInfo{
//				"aliyun_sms": {
//					BaseCredentials: BaseCredentials{
//						BaseUrl:   "dysmsapi.aliyuncs.com",
//						AppKey:    "aliyun sms app key",
//						AppSecret: "aliyun sms secret",
//						ProxyUrl:  "some",
//						Debug:     true,
//					},
//					VerifyTemplate: SmsTemplate{
//						SignName:         "xxx公司",
//						TemplateCode:     "SMS_xxxx",
//						Content:          "【xxx公司】您的动态登录密码为：{code}，有效期为1分钟，请勿泄露与他人。",
//						TemplateParamKey: []string{"code"},
//						CodeKey:          "code",
//					},
//				},
//			},
//		},
//	)

type SmsConfig struct {
	Vendor     string             `json:"vendor,omitempty"`       //针对SmsInfoMap中的key，理论上也要对应SmsInfo中的Vendor
	SmsInfoMap map[string]SmsInfo `json:"sms_info_map,omitempty"` //每个key代表一个短信供应商的配置供应商
}
type SmsInfo struct {
	BaseCredentials
	Vendor         string                 `json:"vendor,omitempty"`          //短信供应商的唯一名字
	VerifyTemplate SmsTemplate            `json:"verify_template,omitempty"` //验证类型模版，单独放一个主要是避免提前定义TemplateMap里面的Key
	TemplateMap    map[string]SmsTemplate `json:"template_map,omitempty"`    //模版集合，key为模版ID ，value为模版信息
}
type SmsTemplate struct {
	SignName         string   `json:"sign_name"`
	Content          string   `json:"content"` //这里填充的是模版内容，比如："您的验证码为{code}。"
	TemplateCode     string   `json:"template_code"`
	TemplateParamKey []string `json:"template_param_key"` //短信模版（或者Content）里面由{}包括起来的key，按顺序填写
	CodeKey          string   //这个为验证码类型的模版定制的，模版参数中指定验证码的key
}

// GetSmsConfig
// 获取当前短信配置，如果在SmsConfig中Vendor没有指定，则随机返回一个（或者默认返回aliyun）
func GetSmsConfig(db *gorm.DB) (SmsInfo, bool) {
	cfg, err := GetConfig[SmsConfig](db, ConfigSmsCategory, ConfigSmsName)
	if err != nil {
		logger.Errorf("load sms config err %s", err)
		return SmsInfo{}, false
	}
	si, exists := getSmsInfo(cfg)
	if !exists && isAppDb(db) {
		return GetSmsConfig(global.GetLocalGorm())
	}
	return si, exists
}
func getSmsInfo(cfg SmsConfig) (SmsInfo, bool) {
	if cfg.SmsInfoMap == nil {
		cfg.SmsInfoMap = map[string]SmsInfo{}
	}
	if cfg.Vendor == "" && len(cfg.SmsInfoMap) > 0 {
		for k, v := range cfg.SmsInfoMap {
			logger.Infof("The randomly selected sms vendor is %s", k)
			return v, true
		}
	}
	if si, ok := cfg.SmsInfoMap[cfg.Vendor]; ok {
		return si, true
	}

	return SmsInfo{}, false
}
