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
		rsaKey, _ = GenKeypair(AlgoASymRsaPKCS1v15, 2048)
		sm2Key, _ = GenKeypair(AlgoASymSm2Pkcs8C132, 2048)
		ecMap     = map[string]string{
			AlgoSymAesEcbPkcs7padding: key,
			AlgoSymAesCbcPkcs7padding: key,
			AlgoSymSm4EcbPkcs7padding: key,
			AlgoSymSm4CbcPkcs7padding: key,
			AlgoSymSm4Cfb:             key,
			AlgoSymSm4Ofb:             key,
			AlgoSymSm4Gcm:             key,
			AlgoASymSm2Pkcs8C132:      sm2Key.PublicKey,
			AlgoASymRsaPKCS1v15:       rsaKey.PrivateKey,
		}
		decMap = map[string]string{
			AlgoSymAesEcbPkcs7padding: key,
			AlgoSymAesCbcPkcs7padding: key,
			AlgoSymSm4EcbPkcs7padding: key,
			AlgoSymSm4CbcPkcs7padding: key,
			AlgoSymSm4Cfb:             key,
			AlgoSymSm4Ofb:             key,
			AlgoSymSm4Gcm:             key,
			AlgoASymSm2Pkcs8C132:      sm2Key.PrivateKey,
			AlgoASymRsaPKCS1v15:       rsaKey.PrivateKey,
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

	keypair, err := GenKeypair(AlgoASymSm2Pkcs8C132, 2048)
	assert.Nil(t, err)

	enc, err := EncryptBase64(AlgoASymSm2Pkcs8C132, keypair.PublicKey, sm4Key1)
	assert.Nil(t, err)
	bs, err := DecryptBase64(AlgoASymSm2Pkcs8C132, keypair.PrivateKey, enc)
	assert.Nil(t, err)

	assert.Equal(t, sm4Key1, string(bs))
}
