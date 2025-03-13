package crypt

import (
	"errors"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/golang-module/dongle"
)

type Crypt struct {
	src       []byte
	dst       []byte
	secretKey []byte
	iv        []byte
	Error     error
	algo      string //加解密算法
	isEncrypt bool   //是否是解密
}

func (c *Crypt) SetIv(iv string) *Crypt {
	c.iv = []byte(iv)
	return c
}

func (c *Crypt) cryptCheck() *Crypt {
	if c.Error != nil {
		return c
	}
	if len(c.src) < 1 {
		c.Error = errors.New("加解密内容为空")
		return c
	}
	if c.algo == "" {
		c.Error = errors.New("没有指定加解密算法")
		return c
	}

	return c
}

func (c *Crypt) crypt() *Crypt {
	if len(c.dst) > 0 {
		return c
	}
	if c.cryptCheck().Error != nil {
		return c
	}
	var (
		fc cryptFunc
		ok bool
	)
	if c.isEncrypt {
		fc, ok = encryptFuncMap[c.algo]
	} else {
		fc, ok = decryptFuncMap[c.algo]
	}
	if !ok {
		c.Error = errors.New("不支持指定加解密算法")
		return c
	}
	keys := [][]byte{c.secretKey}
	if len(c.iv) > 0 {
		keys = append(keys, c.iv)
	}
	c.dst, c.Error = fc(keys, c.src)
	return c
}

func (c *Crypt) FromRawString(s string) *Crypt {
	c.src = utils.StringToBytes(s)
	return c
}
func (c *Crypt) FromHexString(s string) *Crypt {
	c.src = dongle.Decode.FromString(s).ByHex().ToBytes()
	return c
}
func (c *Crypt) FromBase64String(s string) *Crypt {
	c.src = dongle.Decode.FromString(s).ByBase64().ToBytes()
	return c
}

func (c *Crypt) FromRawBytes(b []byte) *Crypt {
	c.src = b
	return c
}
func (c *Crypt) FromHexBytes(b []byte) *Crypt {
	c.src = dongle.Decode.FromBytes(b).ByHex().ToBytes()
	return c
}
func (c *Crypt) FromBase64Bytes(b []byte) *Crypt {
	c.src = dongle.Decode.FromBytes(b).ByBase62().ToBytes()
	return c
}

func (c *Crypt) ToRawString() (string, error) {
	if c.Error != nil {
		return "", c.Error
	}
	return utils.BytesToString(c.crypt().dst), c.Error
}

func (c *Crypt) ToHexString() (string, error) {
	if c.Error != nil {
		return "", c.Error
	}
	return dongle.Encode.FromBytes(c.crypt().dst).ByHex().ToString(), c.Error
}

func (c *Crypt) ToBase64String() (string, error) {
	if c.Error != nil {
		return "", c.Error
	}
	return dongle.Encode.FromBytes(c.crypt().dst).ByBase64().ToString(), c.Error
}

func (c *Crypt) ToRawBytes() ([]byte, error) {
	if c.Error != nil {
		return []byte(""), c.Error
	}
	return c.crypt().dst, c.Error
}

func (c *Crypt) ToHexBytes() ([]byte, error) {
	if c.Error != nil {
		return []byte(""), c.Error
	}
	return dongle.Encode.FromBytes(c.crypt().dst).ByHex().ToBytes(), c.Error
}

func (c *Crypt) ToBase64Bytes() ([]byte, error) {
	if c.Error != nil {
		return []byte(""), c.Error
	}
	return dongle.Encode.FromBytes(c.crypt().dst).ByBase64().ToBytes(), c.Error
}

func NewEncrypt(algo string, key string) *Crypt {
	return &Crypt{
		secretKey: []byte(key),
		algo:      algo,
		isEncrypt: true,
	}
}
func NewDecrypt(algo string, key string) *Crypt {
	return &Crypt{
		secretKey: []byte(key),
		algo:      algo,
		isEncrypt: false,
	}
}
