package oauth2_test

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/davycun/eta/pkg/common/crypt"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/http_tes"
	"github.com/davycun/eta/pkg/common/id/nanoid"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/user"
	"github.com/davycun/eta/pkg/module/user/login/oauth2"
	"github.com/davycun/eta/pkg/module/user/userkey"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
	"time"
)

func TestGetAccessToken(t *testing.T) {

	key, secure := newKeySecure(t)

	lr := &oauth2.LoginResult{}
	rd := &dto.ControllerResponse{
		Result: lr,
	}
	openapiTestcase := http_tes.HttpCase{
		Desc:         "openapi获取access_token",
		Method:       "GET",
		Path:         "/oauth2/access_token",
		Headers:      map[string]string{"Content-Type": "application/json"},
		ResponseDest: rd,
		ShowBody:     true,
		Code:         "200",
	}

	algo := crypt.AlgoSignHmacSha256
	nowTs := time.Now().UTC().Unix()
	nonce := convertor.ToString(nowTs) + convertor.ToString(gofakeit.Number(10000, 99999))
	calcSign, err := crypt.NewEncrypt(algo, secure).FromRawString(fmt.Sprintf("%s%s", convertor.ToString(nowTs), nonce)).ToHexString()
	assert.Nil(t, err)

	uri := url.URL{}
	uri.Path = "/oauth2/access_token"
	query := uri.Query()
	query.Set("algo", algo)
	query.Set("access_key", key)
	query.Set("nonce", nonce)
	query.Set("ts", convertor.ToString(nowTs))
	query.Set("sign", calcSign)
	uri.RawQuery = query.Encode()
	openapiTestcase.Path = uri.String()

	http_tes.Call(t, openapiTestcase)
	assert.NotEmpty(t, lr.Authorization)

	// 重复调用接口，返回“请求重复”
	openapiRepeatTestcase := openapiTestcase
	openapiRepeatTestcase.Code = "400"
	openapiRepeatTestcase.HttpCode = 400
	http_tes.Call(t, openapiRepeatTestcase)

	// token有效期内再次获取 access token, 返回当前有效的 access token
	lr2 := &oauth2.LoginResult{}
	rd2 := &dto.ControllerResponse{
		Result: lr2,
	}

	nonce = convertor.ToString(nowTs) + convertor.ToString(gofakeit.Number(10000, 99999))
	calcSign, err = crypt.NewEncrypt(algo, key).FromRawString(fmt.Sprintf("%s%s", convertor.ToString(nowTs), nonce)).ToHexString()
	query.Set("nonce", nonce)
	query.Set("sign", calcSign)
	uri.RawQuery = query.Encode()
	openapiTestcase1 := openapiTestcase
	openapiTestcase1.ResponseDest = rd2
	openapiTestcase1.Path = uri.String()
	http_tes.Call(t, openapiTestcase1)
	assert.Equal(t, lr.Authorization, lr2.Authorization)
}

func TestFixedToken(t *testing.T) {
	key := newFixToken(t)
	http_tes.Call(t, http_tes.HttpCase{
		Desc:   "TestFixedToken",
		Method: "GET",
		Path:   "/user/current",
		Headers: map[string]string{
			constants.HeaderAuthorization: key,
		},
		//Headers:  map[string]string{"Content-Type": "application/json"},
		Body:     map[string]any{},
		ShowBody: true,
	})
}

func newKeySecure(t *testing.T) (string, string) {
	key, secure := nanoid.New()[:16], nanoid.New()

	var (
		localDB = global.GetLocalGorm()
	)

	us, err := user.LoadDefaultUser(localDB)
	assert.Nil(t, err)

	u2kList := []userkey.UserKey{
		{
			UserId:       us.ID,
			AccessKey:    ctype.NewStringPrt(key),
			AccessSecure: ctype.NewStringPrt(secure),
		},
	}
	err = dorm.Table(localDB, constants.TableUserKey).Create(&u2kList).Error
	assert.Nil(t, err)
	return key, secure
}
func newFixToken(t *testing.T) string {
	fixToken := constants.LoginTypeFixToken + "_" + nanoid.New()[:16]

	var (
		localDB = global.GetLocalGorm()
	)

	us, err := user.LoadDefaultUser(localDB)
	assert.Nil(t, err)

	u2kList := []userkey.UserKey{
		{
			UserId:     us.ID,
			FixedToken: ctype.NewStringPrt(fixToken),
		},
	}
	err = dorm.Table(localDB, constants.TableUser2App).Create(&u2kList).Error
	assert.Nil(t, err)
	return fixToken
}
