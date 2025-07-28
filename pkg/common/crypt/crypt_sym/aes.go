package crypt_sym

import (
	"errors"
	"github.com/dromara/dongle"
)

var (
	AesIvSize = 16 // aes iv 长度
	AesIv     = []byte("0000000000000000")
)

// EncryptAesEcbPkcs7padding
// key[0] 是key
// key[1] 是iv向量，不支持
func EncryptAesEcbPkcs7padding(key [][]byte, src []byte) ([]byte, error) {
	cp := dongle.NewCipher()
	cp.SetMode(dongle.ECB)      // mode: CBC、CFB、OFB、CTR、ECB
	cp.SetPadding(dongle.PKCS7) // padding: No、Empty、Zero、PKCS5、PKCS7、AnsiX923、ISO97971
	cp.WithKey(key[0])          // key: key must be 16, 24 or 32 bytes
	e := dongle.Encrypt.FromBytes(src).ByAes(cp)
	return e.ToRawBytes(), e.Error
}

// DecryptAesEcbPkcs7padding
// key[0] 是key
// key[1] 是iv向量，不支持
func DecryptAesEcbPkcs7padding(key [][]byte, src []byte) ([]byte, error) {
	cp := dongle.NewCipher()
	cp.SetMode(dongle.ECB)      // CBC、CFB、OFB、CTR、ECB
	cp.SetPadding(dongle.PKCS7) // No、Empty、Zero、PKCS5、PKCS7、AnsiX923、ISO97971
	cp.WithKey(key[0])          // key must be 16, 24 or 32 bytes
	d := dongle.Decrypt.FromRawBytes(src).ByAes(cp)
	return d.ToBytes(), d.Error
}

// EncryptAesCbcPkcs7padding
// key[0] 是key
// key[1] 是iv向量
func EncryptAesCbcPkcs7padding(key [][]byte, src []byte) ([]byte, error) {

	cp := dongle.NewCipher()
	cp.SetMode(dongle.CBC)      // CBC、CFB、OFB、CTR、ECB
	cp.SetPadding(dongle.PKCS7) // No、Empty、Zero、PKCS5、PKCS7、AnsiX923、ISO97971
	cp.WithKey(key[0])          // key must be 16, 24 or 32 bytes

	iv := AesIv
	if len(key) > 1 {
		iv = key[1]
	}
	if len(iv) != 16 {
		return nil, errors.New("iv must 16 byte length")
	}
	cp.WithIV(iv)
	e := dongle.Encrypt.FromBytes(src).ByAes(cp)

	return e.ToRawBytes(), e.Error
}

// DecryptAesCbcPkcs7padding
// key[0] 是key
// key[1] 是iv向量
func DecryptAesCbcPkcs7padding(key [][]byte, src []byte) ([]byte, error) {
	cp := dongle.NewCipher()
	cp.SetMode(dongle.CBC)      // CBC、CFB、OFB、CTR、ECB
	cp.SetPadding(dongle.PKCS7) // No、Empty、Zero、PKCS5、PKCS7、AnsiX923、ISO97971
	cp.WithKey(key[0])          // key must be 16, 24 or 32 bytes

	iv := AesIv
	if len(key) > 1 {
		iv = key[1]
	}
	if len(iv) != 16 {
		return nil, errors.New("iv must 16 byte length")
	}
	cp.WithIV(iv)
	d := dongle.Decrypt.FromRawBytes(src).ByAes(cp)
	return d.ToBytes(), d.Error
}
