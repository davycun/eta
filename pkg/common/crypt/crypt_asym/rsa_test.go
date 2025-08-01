package crypt_asym_test

import (
	"github.com/davycun/eta/pkg/common/crypt/crypt_asym"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/dromara/dongle"
	"github.com/dromara/dongle/openssl"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	RsaPublicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAlT/vEvLOLmH3KM5a1r09
OveVHORnfSDOrMloY1H0nm6MI/VieEl3rEgtWf4MRRqDS3wRiju7z8bLT5bc8c0i
O6q8fePpCGO9eAJNeiKnQC/dkW6HPeQTQiaYZt5Peem4BzUooDEOIRjsb9EraiYF
57evgGvjh4gSQMhAdJFPAskqs6fABdjfegt+vYl3KujBoW8IIORTiPAawWut14PS
RaaD+3FsxY/H4kRByTkY0dq6E6qRBqm6NHAspkRGd5oRMoXuZ9jrpV/aQ2hkIIXu
eTMawQD2zrOcfy143Xqrp/XMHdm24d9Zqzj/6igb/ZjDqjPkIdcK7PfxsJsb1LRT
uwIDAQAB
-----END PUBLIC KEY-----`
	RsaPrivateKey = `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCVP+8S8s4uYfco
zlrWvT0695Uc5Gd9IM6syWhjUfSebowj9WJ4SXesSC1Z/gxFGoNLfBGKO7vPxstP
ltzxzSI7qrx94+kIY714Ak16IqdAL92Rboc95BNCJphm3k956bgHNSigMQ4hGOxv
0StqJgXnt6+Aa+OHiBJAyEB0kU8CySqzp8AF2N96C369iXcq6MGhbwgg5FOI8BrB
a63Xg9JFpoP7cWzFj8fiREHJORjR2roTqpEGqbo0cCymREZ3mhEyhe5n2OulX9pD
aGQghe55MxrBAPbOs5x/LXjdequn9cwd2bbh31mrOP/qKBv9mMOqM+Qh1wrs9/Gw
mxvUtFO7AgMBAAECggEARHTqx5ovD/9HSqQ77jsmlqFw96ub/DzMD1ziUIwK05CJ
HwUygVHRXmhMxPZN0nRyvHDP6sOzRX49Sug7t30LsqqBgozDGmIFScJknxy98icC
Te6QgcbXPoRNawGVGqolCslLNQ7LGEtCR5d9flaqZrpN+W2DE2tKVASF6/Gqd+/y
IZ154qyrvgfVYGdm8VQChftC0blQLcGDyF6jen0SNodlN8E4yyZ+rn1VjCH+tlfM
yaavq6XaPb5upNEEh8U3SdEXf6S4FGVQfcAtnQiToCCxQ0BV7E9rT/Gz8d7l2OMe
R0syCqQPqRDrl5Bg5suuyl8eRY0AfjWT3FCGGejRAQKBgQDDkw8eaD7lWzZ1uhL2
cArHI9ccptOg+KfMcSdjLHXjRMqGn7qGQeTwrG+kJyixRqyvYyDAPtwQumiVv+Ax
wjh6fGBqc2FAeTOCDt9Uz/wZz/UiixKQQKnbBa1Oazri1zV6bMH6fDV8oZfq6OV9
TNeSq/AtYxK4QZY+RcjOxWZhGwKBgQDDXNOrqeCS2IZzq26JgzamxoMH8mqdCtWp
1KhXq+bArZccq+IXKAECkSAJKhcgCbXvbNI6FxgcynOsz28SjkJtVJ1Ye0L43aQc
OW9hTJotKBBpnTzGTn69J1xdPbsALFgnooAQ9yNauoZejC2qczGgig1FrVTpjQvF
UcfBGWWh4QKBgDKweOiuqC6V94WH1sZcv85hca2RZ6R/Di6k3UqNCXkAKWW/HH8T
sRzX9I+dPqTD5poGnUR2hl8nkVfOxXLgHfdRKUQt53TodPsuk5/N0E94YNa1KPiO
affEfuimTjrhAJFcguJDMzG8SD2wY1qYgf8X5UY+OWncRe6Z87Pz4dclAoGACRqU
SpWZ/33TliRQ/Ft++nqZtI8ZZMQSfN2KErvR/vyX5CAmYwncMjBtG8A4X6fUMJoT
md1lpEHS7iSkemrisZGV23+y+UHq2d3bUN9u99e8HA/VuzABO/NPnJC53CI04XPz
H9dEcH/srw89OYowr9h/EdYn9NI70DAlbNwwHkECgYEAptsCdPX4EYi6f6pFoEGS
w+PkqcqmM6ErMl86vn96nbb0l+hZG4sqAzW10mcHuYb5DLpYdHeq3gcwN/BS2qYB
1ZED6iDJ98ZBI56JXxfrUCWZOP9oKWkIatM8/bhonHgYZjIbDfs8nKK/wnUWdH7j
T9+7G/08Uh75Y65iNcnDycA=
-----END PRIVATE KEY-----`
	RsaPlaintext  = "CBC模式引入一个新的概念：初始向量IV。IV的作用和MD5的\"加盐\"有些类似，目的是防止同样的明文块始终加密成相同的密文块。\nCBC模式原理:在每个明文块加密前会让那个明文块和IV向量先做异或操作。IV作为初始化变量，参与第一个明文块的异或，后续的每个明文块和它前一个明文块所加密出的密文块相异或，这样相同的明文块加密出来的密文块显然不一样。AES 算法在对明文加密的时候，并不是把整个明文加密成一整段密文，而是把明文拆分成几组独立的明文块，每一个明文块的长度128bit（16B），最后不足128bit(16B),会根据不同的Padding 填充模式进行填充，然后进行加密。\n总结：加密过程是先处理pading，后加密。解密过程是先进行分块解密，最后在处理Padding。AES是加密算法其中的一种，它是属于对称加密，对称加密的意思就是，加密以及解密用的都是同一个Key。相比于非对称加密RSA，SM2等，它的优点就是快。\n"
	RsaCiphertext = "CiN93TbV/EzrTNmmZnjX382UNwFMtZtftiexh/pg/2tXm2kype10D5HgSthzSQnh94HRnkSyMoyNvvtzOfG+zX9nPW8c2OzxohMpb3JZlA/uIO6+4Cm+/81AVFl+xmHjFu7gg5Uik6kcPxyRnl9WZXha8FstpvSBb+PMamKVMrBwpmswGuSM/FyEsrntxZHxP7xAOuUpdnW8Jg2MsoYgHKRXgE0d8FbMfFdOT5FlITghmLY5LJJjKSucESYoKT3UJYdZ6/4tJq19jsEGkjC1UFBa4Am4AwIgybSnlkQ0Yo1B6vFk7m/JdfEqB7AMOWptht19Rti9LvQTn0p+nu0sjI4QXLSGK5xjCwDtICyojRqFG5sf23Ji+2w0nqghzPtPOpcYNmi+UM+3LJIM2iEBtBZ7D1EDS6kDMms8vEzjFHAEmN1RiV2kTjRveKXySFWBHnxLijVPfFj/WfhCeUEAE8GTzl3uPOY4t9q1IFqfgPrDuYqCS6XemAGMRv8T3UQzKEbO2Zxnf3rUgsuR30WwZ7kyvCPapjl9BtiqYUsmW9zoHXp8ieS8+GDDFt5KZhrywcUgiMzkBDsZwKq4JA2famc3f+sxCTZux/P+QQRcYQk2lJ3lb826o/7m2yd0kGlF+udt9q3Zp3nCdrhqYIUTKWt62U98T/+grDISLyQ+hlFcNkrHwYB/URFdvWGSIY0yVE7I/4I3AMUeAGxSqDGLP3WuaIB7NbBXYiyTjXJWAyeRhTWgHEXniRVP4YE1xmRDvcOvgPkl9tdKrG/ZGh6+NCfRITqOEHyWxDTQ39BVVeYircWBwsGXGfOrteohLxVqrG37fA89nOlbixsV89ErdrYu8qFZGI4vjWU7wX+NDJacZvhZBGec1TFrpuAvJyPbf+XKWSbYqc6vu+h0nUhYgMcXRiO7CEipod0ZptqhUYckhwgRkly7XZk4xFIQFGDRLBN4JJbYcWXixHXYX1oZX06yiM+PPgVIN1DPMMzeDQEA3PUtiYQb7j5QHFbiOv9xEf4c4htKe382EDCKBdT4y6+3a5jl3kaLNb0Klm0Wuc7M9ash0NLoLdyFYVJjvoQE31+GOAo/zdos+RfksoO4KYHu+VKG8Gja4F2dDoV1rt/pF3B4Pbowfcd7bNmxsFa9JceQwX0iTEiSFsvWKKRUmLndezLClWgDEK6+pdOquwOrJXqb4oU6m91eLqTogiTXh5LEZajbldclWYeegAMHSlknQy2jomzl8hWnvFugjtO54XwlO/ySftgZeR759ToAClad/GdzD9xpWh6fmj1lxBvvWkYvFQanIfq4zN9b8KWMagMNc/0MS4ZU5gEb0YcV5qTrXQlTBHFx7R19sirW4QpRC92Rqucz6mvQEfdwJdPq7rxN+4xa4KUezOeP6XhKOQ6k1PfgvTAa47Kv33EXegpOPWxvauX/YGIYSfuUMzKt2yRzHxfNBlKcGsRz83lT7RhKv/Xca/iY1A671r7VGihFgSYqAX1f6EqyzjEMYLkwsbtS+MQuFyPR4YfvU6kj9JdC87riebdXKBhFKk1vbM0jCGyVDh/zybNwtkO6twg2/sU3fdLK65LUEttGtdTr+vfsn2ML/YvUfASjVIGXj47zu8Cyku+loBfyhcSloQ3SRA8fUuzSHvMM2iNA1ytneDtWCBP2+DVgZVbpQEY+0AxZOsjTb0zR2N+r9IkYDpY="
)

func TestRsaKenGen(t *testing.T) {
	publicKey8, privateKey8 := openssl.RSA.GenKeyPair(openssl.PKCS8, 2048)
	logger.Infof("生成 2048 字节 PKCS#8 格式公钥：\n%s", publicKey8)
	logger.Infof("生成 2048 字节 PKCS#8 格式私钥：\n%s", privateKey8)
}

func TestRsaEnc(t *testing.T) {
	ciphertext := dongle.Encrypt.FromString(RsaPlaintext).ByRsa(utils.StringToBytes(RsaPublicKey)).ToBase64String()
	assert.Equal(t, RsaCiphertext, ciphertext)
	logger.Infof("加密结果: %s", ciphertext)
}

func TestRsaDec(t *testing.T) {
	plaintext := dongle.Decrypt.FromBase64String(RsaCiphertext).ByRsa(utils.StringToBytes(RsaPrivateKey)).ToString()
	logger.Infof("解密结果: %s", plaintext)
}

func TestRsaDec1(t *testing.T) {
	encStr := dongle.Encrypt.FromString(RsaPlaintext).ByRsa(utils.StringToBytes(RsaPublicKey)).ToBase64String()
	logger.Infof("ciphertext: %s", encStr)
	decStr := dongle.Decrypt.FromBase64String(encStr).ByRsa(utils.StringToBytes(RsaPrivateKey)).ToString()
	logger.Infof("ciphertext: %s", decStr)
}

func TestRsa(t *testing.T) {
	var (
		key, err = crypt_asym.GenRsaPKCS8Key(2048)
		ds       = []string{
			"CBC模式引入一个新的概念：初始向量IV。IV的作用和MD5的\"加盐\"有些类似，目的是防止同样的明文块始终加密成相同的密文块。\nCBC模式原理:在每个明文块加密前会让那个明文块和IV向量先做异或操作。IV作为初始化变量，参与第一个明文块的异或，后续的每个明文块和它前一个明文块所加密出的密文块相异或，这样相同的明文块加密出来的密文块显然不一样。AES 算法在对明文加密的时候，并不是把整个明文加密成一整段密文，而是把明文拆分成几组独立的明文块，每一个明文块的长度128bit（16B），最后不足128bit(16B),会根据不同的Padding 填充模式进行填充，然后进行加密。\n总结：加密过程是先处理pading，后加密。解密过程是先进行分块解密，最后在处理Padding。AES是加密算法其中的一种，它是属于对称加密，对称加密的意思就是，加密以及解密用的都是同一个Key。相比于非对称加密RSA，SM2等，它的优点就是快。\n",
			"这是￥9857！@",
			"/@%$^@&*!()_+:?><SHDETYwle函数",
		}
	)
	assert.Nil(t, err)
	priKey := [][]byte{[]byte(key.PrivateKey)}
	pubKey := [][]byte{[]byte(key.PublicKey)}
	for _, v := range ds {
		ciphertext, err1 := crypt_asym.EncryptRsaPKCS1v15(pubKey, []byte(v))
		assert.Nil(t, err1)
		plaintext, err1 := crypt_asym.DecryptRsaPKCS1v15(priKey, ciphertext)
		assert.Nil(t, err1)
		assert.Equal(t, v, string(plaintext))
	}
}
func TestSm2(t *testing.T) {
	var (
		key, err = crypt_asym.GenSm2PKCS8C132Key(2048)
		ds       = []string{
			"CBC模式引入一个新的概念：初始向量IV。IV的作用和MD5的\"加盐\"有些类似，目的是防止同样的明文块始终加密成相同的密文块。\nCBC模式原理:在每个明文块加密前会让那个明文块和IV向量先做异或操作。IV作为初始化变量，参与第一个明文块的异或，后续的每个明文块和它前一个明文块所加密出的密文块相异或，这样相同的明文块加密出来的密文块显然不一样。AES 算法在对明文加密的时候，并不是把整个明文加密成一整段密文，而是把明文拆分成几组独立的明文块，每一个明文块的长度128bit（16B），最后不足128bit(16B),会根据不同的Padding 填充模式进行填充，然后进行加密。\n总结：加密过程是先处理pading，后加密。解密过程是先进行分块解密，最后在处理Padding。AES是加密算法其中的一种，它是属于对称加密，对称加密的意思就是，加密以及解密用的都是同一个Key。相比于非对称加密RSA，SM2等，它的优点就是快。\n",
			"这是￥9857！@",
			"/@%$^@&*!()_+:?><SHDETYwle函数",
		}
	)
	assert.Nil(t, err)
	priKey := [][]byte{[]byte(key.PrivateKey)}
	pubKey := [][]byte{[]byte(key.PublicKey)}
	for _, v := range ds {
		ciphertext, err1 := crypt_asym.EncryptSm2PKCS8(pubKey, []byte(v))
		assert.Nil(t, err1)
		plaintext, err1 := crypt_asym.DecryptSm2PKCS8(priKey, ciphertext)
		assert.Nil(t, err1)
		assert.Equal(t, v, string(plaintext))
	}
}
