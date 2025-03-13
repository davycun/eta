package crypt

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCrypt(t *testing.T) {

	var (
		key     = "1234567890123456"
		srcList = []string{
			"hello world!",
			`/%429{;/.,<>\nasdiej193973409kziwem_-+=\"1111*&!~#$%^&*()_+}{":><?`,
		}
		ecMap = map[string]string{
			AlgoSymAesEcbPkcs7padding: key,
			AlgoSymAesCbcPkcs7padding: key,
			AlgoSymSm4EcbPkcs7padding: key,
			AlgoSymSm4CbcPkcs7padding: key,
			AlgoSymSm4Cfb:             key,
			AlgoSymSm4Ofb:             key,
			AlgoSymSm4Gcm:             key,
			AlgoAsymSm2Pkcs8C132:      GetPublicKey(AlgoAsymSm2Pkcs8C132),
			AlgoAsymRsaPKCS1v15:       GetPublicKey(AlgoAsymRsaPKCS1v15),
		}
		decMap = map[string]string{
			AlgoSymAesEcbPkcs7padding: key,
			AlgoSymAesCbcPkcs7padding: key,
			AlgoSymSm4EcbPkcs7padding: key,
			AlgoSymSm4CbcPkcs7padding: key,
			AlgoSymSm4Cfb:             key,
			AlgoSymSm4Ofb:             key,
			AlgoSymSm4Gcm:             key,
			AlgoAsymSm2Pkcs8C132:      GetPrivateKey(AlgoAsymSm2Pkcs8C132),
			AlgoAsymRsaPKCS1v15:       GetPrivateKey(AlgoAsymRsaPKCS1v15),
		}
	)

	for _, src := range srcList {
		for algo, k := range ecMap {
			dst, err := NewEncrypt(algo, k).FromRawString(src).ToBase64String()
			assert.Nil(t, err)
			src1, err := NewDecrypt(algo, decMap[algo]).FromBase64String(dst).ToRawString()
			assert.Nil(t, err)
			assert.Equal(t, src, src1)
		}
	}

}

func TestAZDG(t *testing.T) {
	var (
		key      = "chinagdn"
		userName = "admin"
	)
	encodeUserName, err := EncryptBase64(AlgoSymAZDG, key, userName)
	assert.Nil(t, err)

	decodeUserName, err := DecryptBase64(AlgoSymAZDG, key, encodeUserName)
	assert.Nil(t, err)
	assert.Equal(t, userName, string(decodeUserName))
}

func TestSm4(t *testing.T) {

	sm4Key1 := "BEcjIfflBrd8nrCp"
	//publicKey1 := "BPbu6FcFt8zD2Omfh+EECaoay4XsjdDSgh0mEm5P6WLzjaKZBQaDfAVW+fzApmdsW6shscqP9OPjeCle8sbLtEc="
	//privateKey1 := "YfQpfFLG/X8ZOZzlp5n2FxrC3pMbFBrjbTywg1R3lxk="

	enc, err := EncryptBase64(AlgoAsymSm2Pkcs8C132, GetPublicKey(AlgoAsymSm2Pkcs8C132), sm4Key1)
	assert.Nil(t, err)
	bs, err := DecryptBase64(AlgoAsymSm2Pkcs8C132, GetPrivateKey(AlgoAsymSm2Pkcs8C132), enc)
	assert.Nil(t, err)

	assert.Equal(t, sm4Key1, string(bs))
}
