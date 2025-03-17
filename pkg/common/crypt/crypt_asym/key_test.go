package crypt_asym_test

import (
	"github.com/davycun/eta/pkg/common/crypt/crypt_asym"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSm2Gen2(t *testing.T) {
	key, err := crypt_asym.GenSm2PKCS8C132Key(256)
	assert.Nil(t, err)
	assert.Contains(t, key.PrivateKey, "-----BEGIN PRIVATE KEY-----")
	assert.Contains(t, key.PublicKey, "-----BEGIN PUBLIC KEY-----")
}
func TestRSAGen2(t *testing.T) {
	key, err := crypt_asym.GenRsaPKCS8Key(2048)
	assert.Nil(t, err)
	assert.Contains(t, key.PrivateKey, "-----BEGIN PRIVATE KEY-----")
	assert.Contains(t, key.PublicKey, "-----BEGIN PUBLIC KEY-----")
}
