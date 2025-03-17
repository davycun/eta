package sms_sender

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/cache"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/module/setting"
	"github.com/davycun/eta/pkg/module/sms"
	"gorm.io/gorm"
	"time"
)

func SendVerifyCode(txDb *gorm.DB, target sms.Target, code string) error {

	if target.Mobile == "" {
		return errs.NewClientError("the mobile is empty")
	}

	sd, err := NewSender(txDb)
	if err != nil {
		return err
	}

	cfg, ok := setting.GetSmsConfig(txDb)
	//if !ok || cfg.VerifyTemplate.TemplateCode == "" || cfg.AppKey == "" || cfg.AppSecret == "" {
	//	return errs.NewServerError("没有配置短信供应商!")
	//}
	if !ok {
		logger.Warnf("没有配置短信供应商!")
	}

	vt := cfg.VerifyTemplate
	if cfg.TemplateMap != nil && target.TemplateKey != "" {
		if x, ok1 := cfg.TemplateMap[target.TemplateKey]; ok1 {
			vt = x
		}
	}
	if !ctype.IsValid(target.Content) {
		target.Content = ctype.NewStringPrt(vt.Content)
	}
	if target.TemplateCode == "" {
		target.TemplateCode = vt.TemplateCode
	}
	if target.SignName == "" {
		target.SignName = vt.SignName
	}

	if vt.CodeKey == "" {
		vt.CodeKey = "code"
		logger.Errorf("没有配置验证码的Code的Key")
	}
	if target.TemplateParam == nil {
		target.TemplateParam = ctype.Map{}
	}
	target.TemplateParam[vt.CodeKey] = code

	_, err = sd.SendTemplate(target)

	if err != nil {
		return err
	}
	err = saveVerifyCode(target.Mobile, code)
	return err
}
func VerifyVerifyCode(mobile string, code string) bool {
	var (
		key     = getSmsVerifyCodeRedisKey(mobile)
		srcCode = ""
	)

	b, err := cache.Get(key, &srcCode)
	if err != nil {
		logger.Errorf("verify code err %s", err)
		return false
	}
	return b && srcCode == code
}

func saveVerifyCode(phone string, code string) error {
	key := getSmsVerifyCodeRedisKey(phone)
	return cache.SetEx(key, code, time.Second*300)
}

func getSmsVerifyCodeRedisKey(phone string) string {
	return fmt.Sprintf("sms:verify_code:%s", phone)
}
