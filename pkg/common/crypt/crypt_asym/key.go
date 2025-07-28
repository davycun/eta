package crypt_asym

import (
	"crypto/rand"
	"fmt"
	"github.com/dromara/dongle/openssl"
	sm2_key "github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/x509"
	"strings"
)

const (
	algoASymSm2Pkcs8C132 = "sm2_pkcs8_c1c3c2" // 非对称加密
	algoASymRsaPKCS1v15  = "rsa_pkcs1_v1.5"
)

type KeyPair struct {
	Algo       string //算法
	PublicKey  string
	PrivateKey string
}

// Valid
// TODO 其实应该认证校验内容是否合法
func (k KeyPair) Valid() bool {
	//PKCS#8
	return (strings.HasPrefix(k.PrivateKey, "-----BEGIN PRIVATE KEY-----") &&
		strings.HasPrefix(k.PublicKey, "-----BEGIN PUBLIC KEY-----")) ||

		//PKCS#8加密密钥
		(strings.HasPrefix(k.PrivateKey, "-----BEGIN ENCRYPTED PRIVATE KEY-----") &&
			strings.HasPrefix(k.PublicKey, "-----BEGIN ENCRYPTED PUBLIC KEY-----")) ||

		//PKCS#1
		(strings.HasPrefix(k.PrivateKey, "-----BEGIN RSA PRIVATE KEY-----") &&
			strings.HasPrefix(k.PublicKey, "-----BEGIN RSA PUBLIC KEY-----"))
}

// GenRsaPKCS8Key
// bits 1024,2048,3072,4096
func GenRsaPKCS8Key(bits int) (KeyPair, error) {
	if err := checkKeySize(bits); err != nil {
		return KeyPair{}, err
	}
	publicKey, privateKey := openssl.RSA.GenKeyPair(openssl.PKCS8, bits)
	return KeyPair{
		Algo:       algoASymRsaPKCS1v15,
		PublicKey:  string(publicKey),
		PrivateKey: string(privateKey),
	}, nil
}

// GenerateKey returns an error if a key of less than 1024 bits is requested, and all Sign, Verify, Encrypt, and Decrypt methods return an error if used with a key smaller than 1024 bits.
// Such keys are insecure and should not be used.
func checkKeySize(size int) error {
	if size >= 1024 {
		return nil
	}
	return fmt.Errorf("crypto/rsa: %d-bit keys are insecure (see https://go.dev/pkg/crypto/rsa#hdr-Minimum_key_size)", size)
}

// GenSm2PKCS8C132Key
// 暂时支持256 bits
func GenSm2PKCS8C132Key(bits int) (KeyPair, error) {

	key, err := sm2_key.GenerateKey(rand.Reader)
	if err != nil {
		return KeyPair{}, err
	}
	priKeyBytes, err := x509.WritePrivateKeyToPem(key, nil)
	if err != nil {
		return KeyPair{}, err
	}
	pubKeyBytes, err := x509.WritePublicKeyToPem(&key.PublicKey)
	if err != nil {
		return KeyPair{}, err
	}
	return KeyPair{
		Algo:       algoASymSm2Pkcs8C132,
		PublicKey:  string(pubKeyBytes),
		PrivateKey: string(priKeyBytes),
	}, nil
}
