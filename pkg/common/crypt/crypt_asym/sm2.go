package crypt_asym

import (
	"github.com/deatil/go-cryptobin/cryptobin/sm2"
	"github.com/dromara/dongle"
	"strings"
)

// EncryptSm2PKCS8 sm2 C1C3C2 加密
// @param key: base64文本字节
// @param src: 明文字节
func EncryptSm2PKCS8(key [][]byte, src []byte) (rawBytes []byte, err error) {

	k := string(key[0])
	e := sm2.New().FromBytes(src)
	//看是不是PEM格式
	if strings.HasPrefix(k, "-----") {
		e = e.FromPublicKey(key[0])
	} else {
		e = e.FromPublicKeyBytes(dongle.Decode.FromBytes(key[0]).ByBase64().ToBytes())
	}
	e = e.WithMode(sm2.C1C3C2).Encrypt()
	return e.ToBytes(), e.Error()
}

func DecryptSm2PKCS8(key [][]byte, src []byte) (rawBytes []byte, err error) {

	k := string(key[0])
	d := sm2.New().FromBytes(src)

	if strings.HasPrefix(k, "-----") {
		d = d.FromPrivateKey(key[0])
	} else {
		d = d.FromPrivateKeyBytes(dongle.Decode.FromBytes(key[0]).ByBase64().ToBytes())
	}
	d = d.WithMode(sm2.C1C3C2).Decrypt()
	return d.ToBytes(), d.Error()
}
