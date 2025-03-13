package oauth2_test

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/crypt"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/http_tes"
	"github.com/davycun/eta/pkg/common/id/nanoid"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/app"
	"github.com/davycun/eta/pkg/module/user"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestAccountTest(t *testing.T) {

	var (
		us, err1 = user.LoadDefaultUser(global.GetLocalGorm())
		ap, err  = app.LoadDefaultApp(global.GetLocalGorm())
	)
	assert.Nil(t, err1)
	assert.Nil(t, err)
	appId, userId, token, err := http_tes.Login(user.GetRootUser().Account.Data, user.GetRootUser().Password)
	assert.Nil(t, err)
	assert.NotEmpty(t, token)
	assert.Equal(t, ap.ID, appId)
	assert.Equal(t, us.ID, userId)
}

func TestAccountCryptTest(t *testing.T) {

	var (
		pk       = crypt.GetPublicKey(crypt.AlgoAsymSm2Pkcs8C132)
		us, err1 = user.LoadDefaultUser(global.GetLocalGorm())
		ap, err  = app.LoadDefaultApp(global.GetLocalGorm())

		account  = ctype.ToString(user.GetRootUser().Account)
		password = user.GetRootUser().Password
		body     = fmt.Sprintf(`{"username":"%s","password":"%s"}`, account, password)

		symKey    = nanoid.New() //对称加密的
		encSymKey = ""           //把symkey进行非对称加密
		algoSym   = crypt.AlgoSymSm4EcbPkcs7padding
		algoASym  = crypt.AlgoAsymSm2Pkcs8C132
	)

	type CryptRs struct {
		Content string `json:"content"`
	}
	cr := &CryptRs{}

	//对body进行加密
	body, err = crypt.EncryptBase64(algoSym, symKey, body)
	assert.Nil(t, err)
	//对密钥进行加密
	encSymKey, err = crypt.EncryptBase64(algoASym, pk, symKey)
	assert.Nil(t, err)

	cs := http_tes.HttpCase{
		Method:       http.MethodPost,
		Path:         "/oauth2/login",
		HttpCode:     http.StatusOK,
		ResponseDest: cr,
		Headers: map[string]string{
			"Content-Type":                           "application/json",
			constants.HeaderCryptAsymmetricAlgorithm: algoASym,
			constants.HeaderCryptSymmetryAlgorithm:   algoSym,
			constants.HeaderCryptSymmetryKey:         encSymKey,
		},
	}
	http_tes.Call(t, cs)

	assert.Nil(t, err1)
	assert.Nil(t, err)
	appId, userId, token, err := http_tes.Login(us.Account.Data, us.Password)
	assert.Nil(t, err)
	assert.NotEmpty(t, token)
	assert.Equal(t, ap.ID, appId)
	assert.Equal(t, us.ID, userId)
}
