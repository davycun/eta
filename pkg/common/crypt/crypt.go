package crypt

//采用策略设计模式进行架构

import (
	"github.com/davycun/eta/pkg/common/crypt/crypt_asym"
	"github.com/davycun/eta/pkg/common/crypt/crypt_sym"
	"github.com/davycun/eta/pkg/common/crypt/signer"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/utils"
)

const (
	AlgoSignHmacSm3    = "hmac_sm3" //签名
	AlgoSignHmacMd5    = "hmac_md5"
	AlgoSignHmacSha1   = "hmac_sha1"
	AlgoSignHmacSha256 = "hmac_sha256"

	AlgoSymAesEcbPkcs7padding = "aes_ecb_pkcs7padding" // 对称加密
	AlgoSymAesCbcPkcs7padding = "aes_cbc_pkcs7padding"
	AlgoSymSm4EcbPkcs7padding = "sm4_ecb_pkcs7padding"
	AlgoSymSm4CbcPkcs7padding = "sm4_cbc_pkcs7padding"
	AlgoSymSm4Cfb             = "sm4_cfb"
	AlgoSymSm4Ofb             = "sm4_ofb"
	AlgoSymSm4Gcm             = "sm4_gcm"
	AlgoSymAZDG               = "sym_azdg"

	AlgoASymSm2Pkcs8C132 = "sm2_pkcs8_c1c3c2" // 非对称加密
	AlgoASymRsaPKCS1v15  = "rsa_pkcs1_v1.5"
)

type (
	cryptFunc func(key [][]byte, src []byte) ([]byte, error) // key[0] 是key，key[1] 是iv向量
)

var (
	AlgorithmList = []string{
		AlgoSignHmacSm3,
		AlgoSignHmacMd5,
		AlgoSignHmacSha1,
		AlgoSignHmacSha256,
		AlgoSymAesEcbPkcs7padding,
		AlgoSymAesCbcPkcs7padding,
		AlgoSymSm4EcbPkcs7padding,
		AlgoSymSm4CbcPkcs7padding,
		AlgoSymSm4Cfb,
		AlgoSymSm4Ofb,
		AlgoSymSm4Gcm,
		AlgoASymSm2Pkcs8C132,
		AlgoASymRsaPKCS1v15,
	}
	decryptFuncMap = map[string]cryptFunc{
		AlgoSymAesEcbPkcs7padding: crypt_sym.DecryptAesEcbPkcs7padding,
		AlgoSymAesCbcPkcs7padding: crypt_sym.DecryptAesCbcPkcs7padding,
		AlgoSymSm4EcbPkcs7padding: crypt_sym.DecryptSm4EcbPkcs7padding,
		AlgoSymSm4CbcPkcs7padding: crypt_sym.DecryptSm4CbcPkcs7padding,
		AlgoSymSm4Cfb:             crypt_sym.DecryptSm4Cfb,
		AlgoSymSm4Ofb:             crypt_sym.DecryptSm4Ofb,
		AlgoSymSm4Gcm:             crypt_sym.DecryptSm4Gcm,
		AlgoSymAZDG:               crypt_sym.DecryptAZDG,
		AlgoASymSm2Pkcs8C132:      crypt_asym.DecryptSm2PKCS8,
		AlgoASymRsaPKCS1v15:       crypt_asym.DecryptRsaPKCS1v15,
	}
	encryptFuncMap = map[string]cryptFunc{
		AlgoSymAesEcbPkcs7padding: crypt_sym.EncryptAesEcbPkcs7padding,
		AlgoSymAesCbcPkcs7padding: crypt_sym.EncryptAesCbcPkcs7padding,
		AlgoSymSm4EcbPkcs7padding: crypt_sym.EncryptSm4EcbPkcs7padding,
		AlgoSymSm4CbcPkcs7padding: crypt_sym.EncryptSm4CbcPkcs7padding,
		AlgoSymSm4Cfb:             crypt_sym.EncryptSm4Cfb,
		AlgoSymSm4Ofb:             crypt_sym.EncryptSm4Ofb,
		AlgoSymSm4Gcm:             crypt_sym.EncryptSm4Gcm,
		AlgoSymAZDG:               crypt_sym.EncryptAZDG,
		AlgoASymSm2Pkcs8C132:      crypt_asym.EncryptSm2PKCS8,
		AlgoASymRsaPKCS1v15:       crypt_asym.EncryptRsaPKCS1v15,
		AlgoSignHmacSha256:        signer.SignHmacSha256,
		AlgoSignHmacSha1:          signer.SignHmacSha1,
		AlgoSignHmacMd5:           signer.SignHmacMd5,
		AlgoSignHmacSm3:           signer.SignHmacSm3,
	}
)

func ExistsAlgo(algo string) bool {
	return utils.ContainAll(AlgorithmList, algo)
}

func GetAesDefaultIv() string {
	return string(crypt_sym.AesIv)
}
func GetSm4DefaultIv() string {
	return string(crypt_sym.Sm4Iv)
}

func EncryptBase64(algo, key string, plaintext string) (string, error) {
	return NewEncrypt(algo, key).FromRawString(plaintext).ToBase64String()
}

func DecryptBase64(algo, key string, base64Text string) ([]byte, error) {
	return NewDecrypt(algo, key).FromBase64String(base64Text).ToRawBytes()
}

func RegistryEncrypt(algo string, fc cryptFunc) {
	if fc == nil {
		return
	}
	if _, ok := encryptFuncMap[algo]; ok {
		logger.Errorf("registry algo encrypt func has exists")
		return
	}
	encryptFuncMap[algo] = fc
}
func RegistryDecrypt(algo string, fc cryptFunc) {
	if fc == nil {
		return
	}
	if _, ok := decryptFuncMap[algo]; ok {
		logger.Errorf("registry algo decrypt func has exists")
		return
	}
	decryptFuncMap[algo] = fc
}
