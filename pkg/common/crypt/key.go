package crypt

import (
	"errors"
	"fmt"
	"github.com/davycun/eta/pkg/common/crypt/crypt_asym"
)

var (
	keyPairFuncMap = map[string]keyPairFunc{
		AlgoASymRsaPKCS1v15:  crypt_asym.GenRsaPKCS8Key,
		AlgoASymSm2Pkcs8C132: crypt_asym.GenSm2PKCS8C132Key,
	}
)

type keyPairFunc func(bits int) (crypt_asym.KeyPair, error) //gits 表示位数512、1024、2048

// GenKeypair
// 生成非对称密钥对
func GenKeypair(algo string, bits int) (crypt_asym.KeyPair, error) {
	if f, ok := keyPairFuncMap[algo]; ok {
		return f(bits)
	}
	return crypt_asym.KeyPair{}, errors.New(fmt.Sprintf("can not found generate key pair algorithm %s", algo))
}
