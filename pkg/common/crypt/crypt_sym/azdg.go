package crypt_sym

import (
	"crypto/md5"
	"github.com/golang-module/dongle"
	"time"
)

func EncryptAZDG(key [][]byte, inputData []byte) ([]byte, error) {

	//取时间的md5，并且进行16进制编码
	h := md5.New()
	h.Write([]byte(time.Now().Format("2006-01-02 15:04:05")))
	noise := dongle.Encode.FromBytes(h.Sum(nil)).ByHex().ToString() //取当前时间的md5
	h.Reset()

	//明文的每个字节与当前时间的md5的十六进制编码进行异或运算
	loopCount := len(inputData)
	//最终输出的数据内容是明文的两倍，每个明文的字节与当前时间的md5的十六进制编码进行异或运算，
	outData := make([]byte, loopCount*2)
	for i, j := 0, 0; i < loopCount; i, j = i+1, j+1 {
		outData[j] = noise[i%32]
		j++
		outData[j] = inputData[i] ^ noise[i%32]
	}
	//再次对异或后的明文进行密钥的异或加密
	return cipherEncode(key[0], outData), nil
}
func DecryptAZDG(key [][]byte, inputData []byte) ([]byte, error) {
	//对输入的base64密文进行异或解密
	inputData = cipherEncode(key[0], inputData)
	loopCount := len(inputData)
	outData := make([]byte, loopCount/2)
	//
	for i, j := 0, 0; i < loopCount; i, j = i+2, j+1 {
		outData[j] = inputData[i] ^ inputData[i+1]
	}
	return outData, nil
}

func cipherEncode(key []byte, inputData []byte) []byte {
	//对key进行md5，然后转换成16进制
	h := md5.New()
	h.Write(key)
	cipherHash := dongle.Encode.FromBytes(h.Sum(nil)).ByHex().ToString()
	h.Reset()
	loopCount := len(inputData)
	outData := make([]byte, loopCount)
	//对原文中的每个字节与密钥进行异或运算
	for i := 0; i < loopCount; i++ {
		outData[i] = inputData[i] ^ cipherHash[i%32]
	}
	return outData
}
