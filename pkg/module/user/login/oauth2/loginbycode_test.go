package oauth2_test

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/crypt"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/module/user/login/oauth2"
	"github.com/stretchr/testify/assert"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestLoginByCodeShzb(t *testing.T) {
	shzbKey := "zhoupu@secretkey" // 对接上海宗保的key，sm4 ecb
	msg := fmt.Sprintf("%s/%d", "13188888889", time.Now().Unix()*1000)
	hexString, err := crypt.NewEncrypt(crypt.AlgoSymSm4EcbPkcs7padding, shzbKey).FromRawString(msg).ToHexString()
	assert.NoError(t, err)
	logger.Infof("hexString: %s", hexString)
}

func TestLoginByCodeShzbParse(t *testing.T) {
	shzbKey := "zhoupu@secretkey" // 对接上海宗保的key，sm4 ecb
	hexString := "9896a6ea7756b59409c507b895063a9f8f89c4c39b55a5547cfdc3657877c642"
	plainText, err := crypt.NewDecrypt(crypt.AlgoSymSm4EcbPkcs7padding, shzbKey).FromHexString(hexString).ToRawString()
	assert.NoError(t, err)
	if err != nil {
		return
	}
	logger.Infof("明文: %s", plainText)

	info := strings.Split(plainText, "/")
	if len(info) != 2 {
		logger.Errorf("授权码格式错误: %s", plainText)
		return
	}

	// 时间戳
	ms, err1 := strconv.ParseInt(info[1], 10, 64)
	assert.NoError(t, err)
	if err1 != nil {
		return
	}
	nowTs := time.Now().UTC().Unix()
	dur := nowTs - ms/1000
	if dur < -5 || dur > int64(oauth2.RequestValidIn)+5 { // 前后兼容5秒
		logger.Infof("请求过期, dur: %ds", dur)
	}
}
