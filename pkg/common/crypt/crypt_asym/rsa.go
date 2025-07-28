package crypt_asym

import (
	"github.com/dromara/dongle"
)

func EncryptRsaPKCS1v15(key [][]byte, src []byte) (ciphertext []byte, err error) {
	e := dongle.Encrypt.FromBytes(src).ByRsa(key[0])
	return e.ToRawBytes(), e.Error
}

func DecryptRsaPKCS1v15(key [][]byte, src []byte) (plaintext []byte, err error) {
	d := dongle.Decrypt.FromRawBytes(src).ByRsa(key[0])
	return d.ToBytes(), d.Error
}
