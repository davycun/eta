package crypt_sym

import (
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/deatil/go-cryptobin/cryptobin/crypto"
	"github.com/duke-git/lancet/v2/random"
	"github.com/golang-module/dongle"
	"testing"
)

func TestSm4CbcPkcs7paddingEnc(t *testing.T) {
	//key := genKey()
	key := "8eadb267ecd6e860"
	logger.Infof("key: %s", key)
	logger.Infof("key bytes: %v", utils.StringToBytes(key))
	//iv := []byte("1234567890abcdef")
	plaintext := []byte("Hello World!")

	ciphertext, err := EncryptSm4CbcPkcs7padding([][]byte{[]byte(key)}, plaintext)
	ciphertext1, err1 := EncryptSm4CbcPkcs7padding([][]byte{[]byte(key)}, plaintext)
	if err != nil || err1 != nil {
		t.Errorf("Encryption failed: %v", err)
	}
	//logger.Infof("ciphertext: %s", ciphertext)
	//logger.Infof("ciphertext1: %s", ciphertext1)
	//logger.Infof("ciphertext hex: %s", dongle.Encode.FromBytes(ciphertext).ByHex().ToString())
	logger.Infof("ciphertext base64: %s", dongle.Encode.FromBytes(ciphertext).ByBase64().ToString())
	logger.Infof("ciphertext1 base64: %s", dongle.Encode.FromBytes(ciphertext1).ByBase64().ToString())

	cipher1 := crypto.
		FromBytes(plaintext).
		SetKey(key).
		WithIv(Sm4Iv).
		SM4().
		CBC().
		PKCS7Padding().
		Encrypt().
		ToBase64String()
	logger.Infof("cipher1: %s", cipher1)
	cipher2 := crypto.
		FromBytes(plaintext).
		SetKey(key).
		WithIv(Sm4Iv).
		SM4().
		CBC().
		PKCS7Padding().
		Encrypt().
		ToBase64String()
	logger.Infof("cipher2: %s", cipher2)
	plain1 := crypto.
		FromBase64String(cipher1).
		SetKey(key).
		WithIv(Sm4Iv).
		SM4().
		CBC().
		PKCS7Padding().
		Decrypt().
		ToString()
	logger.Infof("plain1: %s", plain1)

	cfbEnc1 := crypto.FromBytes(plaintext).SetKey(key).WithIv(Sm4Iv).SM4().CFB().PKCS7Padding().Encrypt().ToBase64String()
	cfbDec1 := crypto.FromBase64String(cfbEnc1).SetKey(key).WithIv(Sm4Iv).SM4().CFB().PKCS7Padding().Decrypt().ToString()
	logger.Infof("cfbEnc1: %s", cfbEnc1)
	logger.Infof("cfbDec1: %s", cfbDec1)

	ofbEnc1 := crypto.FromBytes(plaintext).SetKey(key).WithIv(Sm4Iv).SM4().OFB().PKCS7Padding().Encrypt().ToBase64String()
	ofbDec1 := crypto.FromBase64String(ofbEnc1).SetKey(key).WithIv(Sm4Iv).SM4().OFB().PKCS7Padding().Decrypt().ToString()
	logger.Infof("ofbEnc1: %s", ofbEnc1)
	logger.Infof("ofbDec1: %s", ofbDec1)

	gcmEnc1 := crypto.FromBytes(plaintext).SetKey(key).WithIv(Sm4Iv).SM4().GCM().PKCS7Padding().Encrypt().ToBase64String()
	gcmDec1 := crypto.FromBase64String(gcmEnc1).SetKey(key).WithIv(Sm4Iv).SM4().GCM().PKCS7Padding().Decrypt().ToString()
	logger.Infof("gcmEnc1: %s", gcmEnc1)
	logger.Infof("gcmDec1: %s", gcmDec1)

}

func genKey() string {
	// SM4 requires a 128-bit (16-byte) key
	return random.RandNumeralOrLetter(16)
}
