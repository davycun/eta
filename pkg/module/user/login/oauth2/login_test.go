package oauth2_test

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/crypt"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/http_tes"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/security"
	"github.com/davycun/eta/pkg/module/user"
	"github.com/davycun/eta/pkg/module/user/login/captcha"
	"github.com/davycun/eta/pkg/module/user/login/oauth2"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	captchaId = ulid.Make().String()
	smsPhone  = "13012345678"
	smsCode   = "123456"
	testcase  = []http_tes.HttpCase{
		{
			Desc:   "用户名密码登录",
			Method: "POST",
			Path:   "/oauth2/login",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: oauth2.LoginByUsernameParam{
				Username: user.RootUserAccount,
				Password: user.RootUserPassword,
			},
			ShowBody: true,
			Code:     "200",
			ValidateFunc: []func(t *testing.T, resp *http_tes.Resp){
				func(t *testing.T, resp *http_tes.Resp) {
					res := resp.Result.(map[string]interface{})
					assert.NotNil(t, res["authorization"])
				},
			},
		},
		{
			Desc:    "短信验证码登录",
			Method:  "POST",
			Path:    "/oauth2/login_by_code",
			Headers: map[string]string{"Content-Type": "application/json"},
			Body: oauth2.LoginByCodeParam{
				LoginType: constants.LoginTypeSmsService,
				Code:      smsCode,
				Phone:     smsPhone,
			},
			ShowBody: true,
			Code:     "200",
			ValidateFunc: []func(t *testing.T, resp *http_tes.Resp){
				func(t *testing.T, resp *http_tes.Resp) {
					res := resp.Result.(map[string]interface{})
					assert.NotNil(t, res["authorization"])
				},
			},
		},
	}
)

func TestSmsLogin(t *testing.T) {

	var (
		code = "123456"
	)

	us := user.NewTestData()
	us.Phone = ctype.NewStringPrt("13900001111")
	us.Password = "123456"
	http_tes.Create(t, "/user/create", []user.User{us})

	captcha.StoreCaptcha(code, captcha.Captcha{Code: code}, time.Second*60)

}

func TestLogin(t *testing.T) {
	err := captcha.StoreCaptcha(captchaId, captcha.Captcha{Code: smsCode, Phone: smsPhone}, time.Second*60)
	assert.NoError(t, err)
	db := global.GetLocalGorm()
	u := &user.User{}
	err = db.Model(u).Where(fmt.Sprintf(`%s=?`, dorm.Quote(dorm.GetDbType(db), "account")), user.RootUserAccount).Updates(map[string]interface{}{"phone": smsPhone}).Error
	assert.NoError(t, err)
	http_tes.Call(t, testcase...)
}

func TestCryptoLogin(t *testing.T) {
	httpCase := []http_tes.HttpCase{
		genCryptTestCase(t, crypt.AlgoASymSm2Pkcs8C132, crypt.AlgoSymSm4CbcPkcs7padding),
		genCryptTestCase(t, crypt.AlgoASymRsaPKCS1v15, crypt.AlgoSymSm4EcbPkcs7padding),
	}

	http_tes.Call(t, httpCase...)
}

func genCryptTestCase(t *testing.T, aSymAlgo string, symAlgo string) http_tes.HttpCase {

	key := http_tes.TransferKey
	//aSymAlgo := crypt.AlgoASymSm2Pkcs8C132
	//symAlgo := crypt.AlgoSymSm4CbcPkcs7padding
	bs64Key, err := crypt.EncryptBase64(aSymAlgo, security.GetPublicKey(aSymAlgo), key)
	assert.Nil(t, err)
	bodyStr := fmt.Sprintf(`{"username":"%s","password":"%s"}`, user.RootUserAccount, user.RootUserPassword)
	bodyB64, err := crypt.EncryptBase64(symAlgo, key, bodyStr)
	assert.Nil(t, err)

	hc := http_tes.HttpCase{
		Desc:   "用户名密码登录-加密",
		Method: "POST",
		Path:   "/oauth2/login",
		Headers: map[string]string{
			"Content-Type":                           "application/json",
			constants.HeaderCryptSymmetryAlgorithm:   symAlgo,
			constants.HeaderCryptAsymmetricAlgorithm: aSymAlgo,
			constants.HeaderCryptSymmetryKey:         bs64Key,
		},
		Body:     bodyB64,
		ShowBody: true,
		Code:     "200",
		ValidateFunc: []func(t *testing.T, resp *http_tes.Resp){
			func(t *testing.T, resp *http_tes.Resp) {
				res := resp.Result.(map[string]interface{})
				assert.NotNil(t, res["authorization"])
			},
		},
	}

	return hc
}
