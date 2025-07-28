package signer

import (
	"github.com/dromara/dongle"
)

func SignHmacSm3(key [][]byte, src []byte) ([]byte, error) {
	e := dongle.Encrypt.FromBytes(src).ByHmacSm3(key[0])
	return e.ToRawBytes(), e.Error
}
func SignHmacMd5(key [][]byte, src []byte) ([]byte, error) {
	e := dongle.Encrypt.FromBytes(src).ByHmacMd5(key[0])
	return e.ToRawBytes(), e.Error
}
func SignHmacSha1(key [][]byte, src []byte) ([]byte, error) {
	e := dongle.Encrypt.FromBytes(src).ByHmacSha1(key[0])
	return e.ToRawBytes(), e.Error
}
func SignHmacSha256(key [][]byte, src []byte) ([]byte, error) {
	e := dongle.Encrypt.FromBytes(src).ByHmacSha256(key[0])
	return e.ToRawBytes(), e.Error
}
