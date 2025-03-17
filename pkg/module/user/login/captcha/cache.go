package captcha

import (
	"github.com/davycun/eta/pkg/common/cache"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/eta/constants"
	"time"
)

func LoadCaptcha(code string) (Captcha, error) {
	var (
		cpt = Captcha{}
	)
	_, err := cache.Get(constants.RedisKey(constants.CaptchaCodeKey, code), &cpt)
	return cpt, err
}

func StoreCaptcha(code string, cpt Captcha, expiration time.Duration) (err error) {
	if expiration == 0 {
		expiration = time.Second * time.Duration(1*60)
	}
	return cache.SetEx(constants.RedisKey(constants.CaptchaCodeKey, code), cpt, expiration)
}

func DelCaptcha(code string) error {
	_, err := cache.Del(constants.RedisKey(constants.CaptchaCodeKey, code))
	return err
}

func Verify(cpt Captcha) bool {
	if cpt.Code == "" {
		return false
	}
	target, err := LoadCaptcha(cpt.Code)
	if err != nil {
		logger.Errorf("load captcha[%s] err %s", cpt.Code, err)
		return false
	}

	if cpt.Verify(target) {
		err = DelCaptcha(cpt.Code)
		if err != nil {
			logger.Errorf("delete captcha[%s] err %s", cpt.Code, err)
		}
		return true
	}
	return false
}
