package crypt_sym_test

import (
	"crypto/md5"
	"fmt"
	"github.com/davycun/eta/pkg/common/crypt/crypt_sym"
	"github.com/golang-module/dongle"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAZDG(t *testing.T) {

	key := "chinadgn"
	username := "admin"
	enCode, err := crypt_sym.EncryptAZDG([][]byte{[]byte(key)}, []byte(username))
	assert.Nil(t, err)
	deCode, err := crypt_sym.DecryptAZDG([][]byte{[]byte(key)}, enCode)
	assert.Equal(t, username, string(deCode))
}

func TestX(t *testing.T) {
	key := "chinadgn"
	h := md5.New()
	h.Write([]byte(key))
	bs := h.Sum(nil)

	keyHash := fmt.Sprintf("%x", h.Sum(nil))
	keyHash2 := dongle.Encode.FromBytes(bs).ByHex().ToString()
	assert.Equal(t, keyHash2, keyHash)
}
