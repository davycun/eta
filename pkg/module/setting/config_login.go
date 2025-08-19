package setting

import (
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/eta/constants"
	"gorm.io/gorm"
	"regexp"
)

type LoginConfig struct {
	// 密码至少8位,必须包含大小写字母、数字、符号：^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[!"#$%&'()*+,-./:;<=>?@[\]^_`{|}~])[A-Za-z\d!"#$%&'()*+,-./:;<=>?@[\]^_`{|}~]{8,}
	// 密码至少6位：.{6,}
	PwdValidateReg          string            `json:"pwd_validate_reg"`           // 密码校验正则表达式
	FailLockDurationMinutes int64             `json:"fail_lock_duration_minutes"` // 登录失败连续时间
	FailLockMaxTimes        int64             `json:"fail_lock_max_times"`        // 登录失败次数，达到后会锁定
	FailLockLockMinutes     int64             `json:"fail_lock_lock_minutes"`     // 登录锁定时间
	CaptchaExpireIn         int64             `json:"captcha_expire_in"`          // 验证码有效时间，秒
	TokenExpireIn           int64             `json:"token_expire_in"`            // token有效时间，秒
	ThirdWeChatMini         BaseCredentials   `json:"third_wechat_mini"`          // 微信小程序
	ThirdWeChat             WechatCredentials `json:"third_wechat"`               // 微信
	ThirdWeCom              BaseCredentials   `json:"third_wecom"`                // 企业微信
	ThirdZzd                ZzdCredentials    `json:"third_zzd"`                  // 浙政钉
	ThirdDing               BaseCredentials   `json:"third_ding"`                 // 钉钉
}
type WechatCredentials struct {
	BaseCredentials
	Name     string `json:"name,omitempty"`
	OriginId string `json:"origin_id,omitempty"`
}
type ZzdCredentials struct {
	BaseCredentials
	TenantId string `json:"tenant_id,omitempty"`
}

// FailLock
// 是否开启登录失败锁定功能
func (lc LoginConfig) FailLock() bool {
	return !(lc.FailLockDurationMinutes <= 0 && lc.FailLockMaxTimes <= 0 && lc.FailLockLockMinutes <= 0)
}

func (lc LoginConfig) PasswordValid(password string) bool {
	// 密码强度检验
	// 密码至少8位,必须包含大小写字母、数字、符号
	// ^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[!"#$%&'()*+,-./:;<=>?@[\]^_`{|}~])[A-Za-z\d!"#$%&'()*+,-./:;<=>?@[\]^_`{|}~]{8,}
	//pattern := `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[!"#$%&'()*+,-./:;<=>?@[\]^_` + "`" + `{|}~])[A-Za-z\d!"#$%&'()*+,-./:;<=>?@[\]^_` + "`" + `{|}~]{8,}`
	if lc.PwdValidateReg != "" {
		match, _ := regexp.MatchString(lc.PwdValidateReg, password)
		return match
	}
	return true
}

func GetLoginConfig(db *gorm.DB) (LoginConfig, bool) {
	cfg, err := GetConfig[LoginConfig](db, ConfigLoginCategory, ConfigLoginName)
	if err != nil {
		logger.Errorf("load login config err %s", err)
	}
	if cfg.TokenExpireIn == 0 {
		cfg.TokenExpireIn = DefaultTokenExpireIn
	}
	return cfg, err == nil
}

func AddDefaultLoginConfig(cf LoginConfig) {
	defaultSettingMap[ConfigLoginCategory+ConfigLoginName] = Setting{
		Namespace: constants.NamespaceEta,
		Category:  ConfigLoginCategory,
		Name:      ConfigLoginName,
		Content:   ctype.Json{Data: &cf, Valid: true},
	}
}
