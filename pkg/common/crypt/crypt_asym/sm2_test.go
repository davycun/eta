package crypt_asym_test

import (
	"crypto/rand"
	"github.com/davycun/eta/pkg/common/crypt/crypt_asym"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/utils"
	csm2 "github.com/deatil/go-cryptobin/cryptobin/sm2"
	"github.com/dromara/dongle"
	"github.com/stretchr/testify/assert"
	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/x509"
	"testing"
)

var (
	Plaintext  = "CBC模式引入一个新的概念：初始向量IV。IV的作用和MD5的\"加盐\"有些类似，目的是防止同样的明文块始终加密成相同的密文块。\nCBC模式原理:在每个明文块加密前会让那个明文块和IV向量先做异或操作。IV作为初始化变量，参与第一个明文块的异或，后续的每个明文块和它前一个明文块所加密出的密文块相异或，这样相同的明文块加密出来的密文块显然不一样。AES 算法在对明文加密的时候，并不是把整个明文加密成一整段密文，而是把明文拆分成几组独立的明文块，每一个明文块的长度128bit（16B），最后不足128bit(16B),会根据不同的Padding 填充模式进行填充，然后进行加密。\n总结：加密过程是先处理pading，后加密。解密过程是先进行分块解密，最后在处理Padding。AES是加密算法其中的一种，它是属于对称加密，对称加密的意思就是，加密以及解密用的都是同一个Key。相比于非对称加密RSA，SM2等，它的优点就是快。\n"
	Ciphertext = "BPvnKO4UP68ioE8HKU+fPZuTGCBnnm0aC+od7jCAE23wjeBptY5GhbC0HGrfzHsso7qxZc8ecHMvPWwwGv9J2+aUq8gq95euvvaOFma9rYBG5m0CBCVhHvnPL/unXsbk8hdWPgqbTnNPefS4+2vn/yLAj8rtMdwRgUOSM224XpfA+nXH3NjUJCWEFHnz8kN5MS7JawYEtfHPjEily//Ojc5/cf3OUL97KAcNYDKlGmaVfqoe6rsWVLGtiZf6A5Q3f97CPH6GuqGTorabvYByvxuzzzQXGhADl9Cf+N6mvvc4wZC8SWz671FqfrHFstxcMLFlPUN5iO8LikWckeNW7b6gl5k5YZMKpYwoJcVXxZNHuNmXDPTpBjS0Rh+eIAdGOPhlHF74r+1mmU4cMQb3tvDHlsDulnHRP8EKlVU25R/vQdJjnHC/mVK0+CkVB5Jiru8uSh2SGUC+mtm8kwWBePtOwAoiITTRw856xcKYDyI19G3T552vUDGxPAHt2lvrX/Pzn66B062VDLE2xgejMYZ84uKCEjiS9SwzPcEVXDbC/abCIJdLt0ZdFS9xdGxT4bwjSbgMyiNguG1EooWU8t92d53S6Mh9vluAp4ZRPbTBMm+TURKt9CTI7xACZw8nxmVUFRp47Z59RKq4o4mhuRlKHqJbgHQCuCDz4pJcsbesa7NRg6bsyhfhg5EfBtBauidacOPJHNHVaTsHU6OQBKnuKThK8sjiqAtG7FH9oWR8kErH5pAwV+H0XiCIxjjVS7wmAyvQrElkQJDfgUdyICEk+Gm++Vw/kpXDXvQD7hjGApdeKMy9cmeTXYvnYs+1H/UhJW9Bftmn7YMmKVOePLOz4MNg6rMiBt1KxODsyzAKqCfJueU3q/s5uDlmGOF8fjLDtEqoOzLOOv9Lghf3XMXiWmuSr1RyG7GerzC3QqArt/XYOblOYZtxBzBRPCllQb6T3gHsu8uD1/lfYZgGezNbm65/la4O6MJ3/WoCh1FUVRzKlWrmLXuWOgpTHk/vVugPQFKtzS6QamwJfoiQbB9vIpJkZLE8FgCRJjIXjaYtMUIwHudOhB9mCul7XKjdmFNVlrcNFFsf4DRzMbtgQCBigLf3q8rFk20PIwOzcjKbSWnULS0XQrAXUjkRQWjireeesT3Jiq/PVfKF7qVX+xuCN8Al0WPfx9B1bGcA3lMlrPD+eKefBXpm1CNQCul43GitvvBuHhxjnw9CT6gSiA2GlSnXpA9n7FEqv+nd3v/DNK6zxJ667XA3iMuy/YMMfI8Vx+SAsq3OUS/xQqFbSRHPjsQWcRoBiknIUkPxW6oYA3rcb5Yv6T8FTjnKCtyE6YSGzwoXK0M5TR33v0+7ma/bjBHN6No8Ko7+nIc9nfZpWEJ6LGZryjp9bxThI2vA56ZBk8cs8SE6lH/77xmBuNeu6aTQkeLz+FnBuAbaaQZsybcKY1MvJNEAwE0z5KOIcOXNZCurjoCscK6QHdrKjecXShraMS11YRLFdRuci4H3EIDk8uTlPKLgV1vuR8m2v2Cme9FzS0wiEphTkPktJMPjxUVDrgJkvg=="

	Sm2PrivateKey = `-----BEGIN PRIVATE KEY-----
MIGTAgEAMBMGByqGSM49AgEGCCqBHM9VAYItBHkwdwIBAQQgB/IIvg4ijmPUOR8p
wKVpVTRXRqXFPKrJj2hg8fqBWTSgCgYIKoEcz1UBgi2hRANCAAQQOlEWQAsN1DDi
zfIKg9xCgWb7BGHJHc+YN+rkhNAEMV/7Qfg8AwIAqEI++tkb8bgVX7knmMPjBbL8
N57dip/u
-----END PRIVATE KEY-----`
	Sm2PublicKey = `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoEcz1UBgi0DQgAEEDpRFkALDdQw4s3yCoPcQoFm+wRh
yR3PmDfq5ITQBDFf+0H4PAMCAKhCPvrZG/G4FV+5J5jD4wWy/Dee3Yqf7g==
-----END PUBLIC KEY-----`
)

func TestSm2KenGen(t *testing.T) {
	pri, err := sm2.GenerateKey(rand.Reader) // 生成密钥对
	assert.Nil(t, err)
	pub := &pri.PublicKey

	priPem, err := x509.WritePrivateKeyToPem(pri, nil)
	assert.Nil(t, err)
	pubPem, err := x509.WritePublicKeyToPem(pub)
	assert.Nil(t, err)

	logger.Infof("private key: \n%s", string(priPem))
	logger.Infof("public key: \n%s", string(pubPem))
}

func TestSm2PKCS8Enc(t *testing.T) {
	key := []byte(Sm2PublicKey)
	msg := []byte(Plaintext)
	ciphertext, err := crypt_asym.EncryptSm2PKCS8([][]byte{key}, msg)
	assert.Nil(t, err)
	logger.Infof("加密结果: %s", dongle.Encode.FromBytes(ciphertext).ByBase64().ToString())
	msg1, err := crypt_asym.DecryptSm2PKCS8([][]byte{[]byte(Sm2PrivateKey)}, ciphertext)
	assert.Nil(t, err)
	assert.Equal(t, Plaintext, string(msg1))
}

func TestSm2PKCS8Dec(t *testing.T) {
	key := utils.StringToBytes(Sm2PrivateKey)
	msg := dongle.Decode.FromString(Ciphertext).ByBase64().ToBytes()
	plaintext, err := crypt_asym.DecryptSm2PKCS8([][]byte{key}, msg)
	assert.Nil(t, err)
	logger.Infof("解密结果: %s", plaintext)
}

func TestSm2PKCS8Dec1(t *testing.T) {
	priv, err := sm2.GenerateKey(rand.Reader) // 生成密钥对
	pub := &priv.PublicKey
	privStr := x509.WritePrivateKeyToHex(priv)
	pubStr := x509.WritePublicKeyToHex(pub)
	text := "8eadb267ecd6e860"
	//privStr := "f94410331dd9d0e1c63f5f1fadb11d85bc6c63310834cf3f00acc71ca69678bb"
	//pubStr := "04aa3599bec6aee42539f907bb752ef57c475a27c516565f976dd0f039560eb3a1c792b32af07cc139d371fa6b9298a1f03c91d698be6aa56375d4497942e668c4"
	//priv, err := x509.ReadPrivateKeyFromHex(privStr)
	//assert.Nil(t, err)
	//pub, err := x509.ReadPublicKeyFromHex(pubStr)
	//assert.Nil(t, err)
	logger.Infof("privStr: %s", privStr)
	logger.Infof("pubStr: %s", pubStr)
	logger.Infof("privStr b64: %s", dongle.Encode.FromBytes(dongle.Decode.FromString(privStr).ByHex().ToBytes()).ByBase64().ToString())
	logger.Infof("pubStr b64: %s", dongle.Encode.FromBytes(dongle.Decode.FromString(pubStr).ByHex().ToBytes()).ByBase64().ToString())

	encStr, err := sm2.Encrypt(pub, utils.StringToBytes(text), rand.Reader, sm2.C1C3C2)
	assert.Nil(t, err)
	logger.Infof("enc_str: %s", dongle.Encode.FromBytes(encStr).ByHex().ToString())
	decStr, err := sm2.Decrypt(priv, encStr, sm2.C1C3C2)
	logger.Infof("dec_str: %s", utils.BytesToString(decStr))
}

func TestCryptobinSM2_1(t *testing.T) {
	key := utils.StringToBytes(Sm2PublicKey)
	msg := utils.StringToBytes(Plaintext)
	ciphertext, _ := crypt_asym.EncryptSm2PKCS8([][]byte{key}, msg)
	logger.Infof("ciphertext: %s", dongle.Encode.FromBytes(ciphertext).ByBase64().ToString())

	//pubKeyHex := dongle.Encode.FromBytes(dongle.Decode.FromBytes(key).ByBase64().ToBytes()).ByHex().ToBytes()
	//pub, _ := x509.ReadPublicKeyFromHex(pubKeyHex)

	kk := csm2.New().GenerateKey()
	priK := kk.CreatePKCS8PrivateKey().ToKeyString()
	pubK := kk.CreatePublicKey().ToKeyString()
	logger.Infof("private pem key: %s", priK)
	logger.Infof("public pem key: %s", pubK)
	priPem, err := x509.ReadPrivateKeyFromPem(utils.StringToBytes(priK), nil)
	assert.Nil(t, err)
	assert.NotNil(t, priPem)
	if priPem == nil || err != nil {
		logger.Errorf("priPem is nill or occur err. %v", err)
		return
	}
	pubPem, err := x509.ReadPublicKeyFromPem(utils.StringToBytes(pubK))
	if err != nil {
		logger.Errorf("err: %s", err)
		return
	}

	priHex := x509.WritePrivateKeyToHex(priPem)
	pubHex := x509.WritePublicKeyToHex(pubPem)
	logger.Infof("private hex key: %s", priHex)
	logger.Infof("public hex key: %s", pubHex)

	priHexB64 := dongle.Encode.FromBytes(dongle.Decode.FromString(priHex).ByHex().ToBytes()).ByBase64().ToString()
	pubHexB64 := dongle.Encode.FromBytes(dongle.Decode.FromString(pubHex).ByHex().ToBytes()).ByBase64().ToString()
	logger.Infof("private hex b64 key: %s", priHexB64)
	logger.Infof("public hex b64 key: %s", pubHexB64)
	// p7Q/yVtwwsXunWjwopCCNA9Scnc/w2VnPupAyqbeXe8=
	// BJmscMf17hgMeaQ0+jAfq3xsZ/gdmmyCKfoCNji9yy743oU2Q+O+TFZxfySEsa9gj0WzWWSfD6om1CObbh8E6N4=

	enc := csm2.New().
		FromString(Plaintext).
		FromPublicKeyString(pubHex).
		WithMode(csm2.C1C3C2). // C1C3C2 | C1C2C3
		Encrypt().
		ToBase64String()
	logger.Infof("enc: %s", enc)
}

func TestCryptobinSM2_2(t *testing.T) {
	key := utils.StringToBytes(Sm2PublicKey)
	msg := utils.StringToBytes(Plaintext)
	ciphertext, _ := crypt_asym.EncryptSm2PKCS8([][]byte{key}, msg)
	logger.Infof("ciphertext: %s", dongle.Encode.FromBytes(ciphertext).ByBase64().ToString())

	enc := csm2.New().
		FromString(Plaintext).
		FromPublicKey(dongle.Decode.FromString(Sm2PublicKey).ByBase64().ToBytes()).
		WithMode(csm2.C1C3C2).
		Encrypt().
		ToBase64String()
	logger.Infof("enc: %s", enc)

	dec := csm2.New().
		FromBase64String(Ciphertext).
		FromPrivateKey(dongle.Decode.FromString(Sm2PrivateKey).ByBase64().ToBytes()).
		WithMode(csm2.C1C3C2).
		Decrypt().
		ToString()
	logger.Infof("dec: %s", dec)

	dec1 := csm2.New().
		FromBytes(ciphertext).
		FromPrivateKey(dongle.Decode.FromString(Sm2PrivateKey).ByBase64().ToBytes()).
		WithMode(csm2.C1C3C2).
		Decrypt().
		ToString()
	logger.Infof("dec1: %s", dec1)
}

func TestX(t *testing.T) {
	text := `{
    "algo": "hmac_sha256",
    "msg": "自 ChatGPT 火爆出圈以来，各式大模型与生成式 AI 技术喷涌而出，医疗、金融、出行、消费零售、互联网等各个行业都在寻找利用生成式 AI 技术赋能业务创新的方法。然而，从摸索到落地，企业在应用用生成式 AI 技术尚存在门槛，使用现成的技术服务又会产生安全等方面的顾虑，比如业务数据泄漏等问题。一时间，许多企业陷入进退两难的境地。本篇就来说道说道企业在落地生成生式 AI 应用过程中的那些事。从模型选择到业务安全，生成式 AI 应用诞生的曲折实际上，生成式 AI 市场经过一年多快速的发展，一方面市场上涌现出很多大模型与配套服务，另一方面，企业面临的难题也在与日俱增。首先，从模型的选择开始，眼花缭乱的大模型就已经够厂家喝上一壶。5 月 28 日，由中国科学技术信息研究所、科技部新一代人工智能发展研究中心联合相关研究机构编写的《中国人工智能大模型地图研究报告》正式发布，报告显示，我国 10 亿参数规模以上大模型已发布近 80 个。再到 10 月中国新一代人工智能发展战略研究院发布的《2023 中国新一代人工智能科技产业发展报告》显示，目前国内大模型总数达 238 个。而据北京经信局数据，截至 10 月初，北京发布大模型数量达 115 个，其中通用大模型 12 个，垂类大模型 103 个。按照这个发展趋势，“百模大战” 或许很快就会升级为 “千模大战”。而企业如何选择大模型便会难上加难。在具体应用场景中，企业需要在准确性和性能平衡间作出衡量，有效地比较模型并根据其首选指标找到最佳选择，这就需要深厚的数据科学专业知识，也会耗费大量的人力时间成本。而模型的选择只是开始，确定模型之后，还需要结合自身业务，做模型的精调、训练等工作。在这一步，公司的业务数据类型与大模型输入所要求的数据类型需要做一定适配，同时输入的数据需要具有代表性、多样性、一致性和可靠性，这样才能实现效果更佳的输出，这便要求企业需要有既懂业务，又懂大模型技术的工程师对数据进行整理。此外，大模型的精调也需要大量的算力，需要投入大量的资金和时间来购买和维护硬件设备，或者租用云服务，同时基础设施也需要长时间的维护。而大模型技术作为新兴技术，许多公司并没有相应的人才储备与经验，这对企业来说也是不小的压力。模型本身的问题解决了之后，企业还要面临安全隐患。比如很多使用方会担心，用了某个大模型，那么自己的数据会不会都被模型方看到甚至泄露？会不会导致敏感信息泄漏，或是生成违规内容等等？那么，企业就需要确保数据在传输、存储和处理的过程中不会被泄露或者滥用，以免给业务和声誉带来损失。亚马逊云科技助力企业安全构建生成式 AI 应用面对种种挑战，对于许企业来说，最佳选择可以是给自己找一个 AI 助手，全方面辅助完成 AI 能力的嵌入，最常见的便是云上大模型平台和服务。2023 年 11 月 28 日，2023 亚马逊云科技 re:Invent 在美国拉斯维加斯盛大开启，并于 12 月 2 日圆满落下帷幕。2023 年 12 月 12 日起，2023 亚马逊云科技 re:Invent 中国行城市巡展活动将在 10 大城市开启，覆盖北京、上海、广州、深圳、成都、青岛、南京、西安、杭州、长沙 10 个城市！",
    "salt": "sed"
}`
	key := utils.StringToBytes(Sm2PublicKey)
	msg := utils.StringToBytes(text)
	ciphertext, err := crypt_asym.EncryptSm2PKCS8([][]byte{key}, msg)
	assert.Nil(t, err)
	logger.Infof("加密结果: %s", dongle.Encode.FromBytes(ciphertext).ByBase64().ToString())
}
