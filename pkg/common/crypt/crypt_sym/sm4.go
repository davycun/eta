package crypt_sym

import (
	"errors"
	"github.com/deatil/go-cryptobin/cryptobin/crypto"
)

// 参考文档：
// sm4: https://github.com/emmansun/gmsm/blob/main/docs/sm4.md

var (
	Sm4IvSize = 16 // sm4 iv 长度
	Sm4Iv     = []byte("0000000000000000")
)

func EncryptSm4EcbPkcs7padding(key [][]byte, src []byte) ([]byte, error) {
	e := crypto.FromBytes(src).WithKey(key[0]).SM4().ECB().PKCS7Padding().Encrypt()
	return e.ToBytes(), e.Error()
}
func DecryptSm4EcbPkcs7padding(key [][]byte, src []byte) ([]byte, error) {
	d := crypto.FromBytes(src).WithKey(key[0]).SM4().ECB().PKCS7Padding().Decrypt()
	return d.ToBytes(), d.Error()
}

func EncryptSm4CbcPkcs7padding(key [][]byte, src []byte) ([]byte, error) {
	iv := Sm4Iv
	if len(key) > 1 {
		iv = key[1]
	}
	if len(iv) != 16 {
		return nil, errors.New("iv must 16 byte length")
	}
	e := crypto.FromBytes(src).WithKey(key[0]).WithIv(iv).SM4().CBC().PKCS7Padding().Encrypt()
	return e.ToBytes(), e.Error()
}
func DecryptSm4CbcPkcs7padding(key [][]byte, src []byte) ([]byte, error) {
	iv := Sm4Iv
	if len(key) > 1 {
		iv = key[1]
	}
	if len(iv) != 16 {
		return nil, errors.New("iv must 16 byte length")
	}
	d := crypto.FromBytes(src).WithKey(key[0]).WithIv(iv).SM4().CBC().PKCS7Padding().Decrypt()
	return d.ToBytes(), d.Error()
}

func EncryptSm4Cfb(key [][]byte, src []byte) ([]byte, error) {
	iv := Sm4Iv
	if len(key) > 1 {
		iv = key[1]
	}
	if len(iv) != 16 {
		return nil, errors.New("iv must 16 byte length")
	}
	e := crypto.FromBytes(src).WithKey(key[0]).WithIv(iv).SM4().CFB().PKCS7Padding().Encrypt()
	return e.ToBytes(), e.Error()
}
func DecryptSm4Cfb(key [][]byte, src []byte) ([]byte, error) {
	iv := Sm4Iv
	if len(key) > 1 {
		iv = key[1]
	}
	if len(iv) != 16 {
		return nil, errors.New("iv must 16 byte length")
	}
	d := crypto.FromBytes(src).WithKey(key[0]).WithIv(iv).SM4().CFB().PKCS7Padding().Decrypt()
	return d.ToBytes(), d.Error()
}

func EncryptSm4Ofb(key [][]byte, src []byte) ([]byte, error) {
	iv := Sm4Iv
	if len(key) > 1 {
		iv = key[1]
	}
	if len(iv) != 16 {
		return nil, errors.New("iv must 16 byte length")
	}
	e := crypto.FromBytes(src).WithKey(key[0]).WithIv(iv).SM4().OFB().PKCS7Padding().Encrypt()
	return e.ToBytes(), e.Error()
}
func DecryptSm4Ofb(key [][]byte, src []byte) ([]byte, error) {
	iv := Sm4Iv
	if len(key) > 1 {
		iv = key[1]
	}
	if len(iv) != 16 {
		return nil, errors.New("iv must 16 byte length")
	}
	d := crypto.FromBytes(src).WithKey(key[0]).WithIv(iv).SM4().OFB().PKCS7Padding().Decrypt()
	return d.ToBytes(), d.Error()
}

func EncryptSm4Gcm(key [][]byte, src []byte) ([]byte, error) {
	iv := Sm4Iv
	if len(key) > 1 {
		iv = key[1]
	}
	if len(iv) != 16 {
		return nil, errors.New("iv must 16 byte length")
	}
	e := crypto.FromBytes(src).WithKey(key[0]).WithIv(iv).SM4().GCM().PKCS7Padding().Encrypt()
	return e.ToBytes(), e.Error()
}
func DecryptSm4Gcm(key [][]byte, src []byte) ([]byte, error) {
	iv := Sm4Iv
	if len(key) > 1 {
		iv = key[1]
	}
	if len(iv) != 16 {
		return nil, errors.New("iv must 16 byte length")
	}
	d := crypto.FromBytes(src).WithKey(key[0]).WithIv(iv).SM4().GCM().PKCS7Padding().Decrypt()
	return d.ToBytes(), d.Error()
}
