package captcha

import (
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/mojocn/base64Captcha"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestImageCaptcha(t *testing.T) {

	var (
		dd = &base64Captcha.DriverDigit{
			Length:   4,   //数字个数
			Height:   50,  //高度
			Width:    120, //宽度
			MaxSkew:  0.7,
			DotCount: 80,
		}
	)

	cp := base64Captcha.NewCaptcha(dd, store)
	id, b64s, answer, err := cp.Generate()

	assert.Nil(t, err)
	logger.Infof(b64s)
	logger.Infof(id)
	logger.Infof(answer)

}
