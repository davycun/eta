package plugin_crypt

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/crypt"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/app"
	"github.com/duke-git/lancet/v2/slice"
	"gorm.io/gorm"
	"reflect"
	"regexp"
	"strings"
	"sync"
)

var (
	//TODO 需要定时清理
	encryptCache = sync.Map{}
	metaFlag     = "@@@"
	dataSplit    = "@"
	dataPattern  = "[0-9a-zA-z]+"
	metaPattern  = metaFlag + dataPattern + dataSplit + dataPattern + dataSplit + dataPattern + metaFlag
	metaCompile  = regexp.MustCompile(metaPattern)
)

func CleanEncryptCache() {
	encryptCache = sync.Map{}
}

func encryptValue(db *gorm.DB, tbName string, cryptFieldInfo entity.CryptFieldInfo, valList ...reflect.Value) error {
	if len(valList) < 1 || !cryptFieldInfo.Enable || cryptFieldInfo.Field == "" {
		return nil
	}
	for _, val := range valList {
		if !val.IsValid() || val.IsZero() || !val.CanInterface() {
			continue
		}
		valInter := val.Interface()
		fieldStr := entity.GetString(valInter, cryptFieldInfo.Field)
		if fieldStr == "" {
			continue
		}
		rs, ok := encryptData(cryptFieldInfo, []rune(fieldStr))
		if rs == "" || !ok {
			continue
		}
		err := entity.SetValue(val, cryptFieldInfo.Field, addEncryptDataPrefix(db, tbName, cryptFieldInfo.Field, rs))
		if err != nil {
			return err
		}
	}

	return nil
}
func encryptData(cryptInfo entity.CryptFieldInfo, msg []rune) (string, bool) {
	if len(msg) <= cryptInfo.KeepTxtPreCnt+cryptInfo.KeepTxtSufCnt {
		return string(msg), false
	}
	var (
		sliceSize   = cryptInfo.SliceSize
		pre         = string(msg[:cryptInfo.KeepTxtPreCnt])
		suf         = string(msg[len(msg)-cryptInfo.KeepTxtSufCnt:])
		content     = string(msg[cryptInfo.KeepTxtPreCnt : len(msg)-cryptInfo.KeepTxtSufCnt])
		encryptFunc = func(plaintext string) string {
			key := fmt.Sprintf("%s@%s", cryptInfo.GetSecretKey(), plaintext)
			if rsData, ok := encryptCache.Load(key); ok {
				return rsData.(string)
			}
			rs, err := crypt.EncryptBase64(cryptInfo.Algo, cryptInfo.GetSecretKey(), plaintext)
			if err != nil {
				logger.Errorf("encryptFunc err for %s,%s", plaintext, err)
				return ""
			}
			encryptCache.Store(key, rs)
			return rs
		}
	)
	if sliceSize <= 0 {
		return strings.Join([]string{pre, constants.CryptPrefix, encryptFunc(content), suf}, ""), true
	}
	contentSlice := crypt.StringToSlice(content, cryptInfo.SliceSize)
	encStrList := slice.Map(contentSlice, func(i int, v string) string { return encryptFunc(v) })
	return strings.Join([]string{pre, constants.CryptPrefix, strings.Join(encStrList, constants.CryptSliceSeparator), suf}, ""), true
}

func decryptValue(db *gorm.DB, tbName string, cryptFieldInfo entity.CryptFieldInfo, valList ...reflect.Value) error {

	if len(valList) < 1 || !cryptFieldInfo.Enable || cryptFieldInfo.Field == "" {
		return nil
	}
	for _, val := range valList {
		if !val.IsValid() || val.IsZero() || !val.CanInterface() {
			continue
		}
		valInter := val.Interface()
		fieldStr := entity.GetString(valInter, cryptFieldInfo.Field)
		_, encryptStr := splitEncryptData(fieldStr)

		rsStr, ok := decryptData(cryptFieldInfo, []rune(encryptStr))
		if rsStr == "" || !ok {
			continue
		}
		err := entity.SetValue(val, cryptFieldInfo.Field, rsStr)
		if err != nil {
			return err
		}
	}

	return nil
}
func decryptData(cryptInfo entity.CryptFieldInfo, msg []rune) (string, bool) {
	if len(msg) <= cryptInfo.KeepTxtPreCnt+cryptInfo.KeepTxtSufCnt {
		return string(msg), false
	}
	var (
		sliceSize   = cryptInfo.SliceSize
		pre         = string(msg[:cryptInfo.KeepTxtPreCnt])
		suf         = string(msg[len(msg)-cryptInfo.KeepTxtSufCnt:])
		content     = string(msg[cryptInfo.KeepTxtPreCnt : len(msg)-cryptInfo.KeepTxtSufCnt])
		decryptFunc = func(plaintext string) string {
			rs, err := crypt.DecryptBase64(cryptInfo.Algo, cryptInfo.GetSecretKey(), plaintext)
			if err != nil {
				logger.Errorf("decryptFunc err for %s,%s", plaintext, err)
				return ""
			}
			return string(rs)
		}
	)

	if !strings.HasPrefix(content, constants.CryptPrefix) {
		return strings.Join([]string{pre, content, suf}, ""), false
	}

	content = strings.TrimPrefix(content, constants.CryptPrefix)
	if sliceSize <= 0 {
		return strings.Join([]string{pre, decryptFunc(content), suf}, ""), true
	}

	encStrList := strings.Split(content, constants.CryptSliceSeparator)
	strList := slice.Map(encStrList, func(i int, v string) string { return decryptFunc(v) })
	return strings.Join([]string{pre, crypt.SliceToString(strList), suf}, ""), true
}

func addEncryptDataPrefix(db *gorm.DB, tbName, field, src string) string {
	id := app.LoadAppIdBySchema(dorm.GetDbSchema(db))
	//fmt性能能最好，builder第二
	return fmt.Sprintf("%s%s%s%s%s%s%s%s", metaFlag, id, dataSplit, tbName, dataSplit, field, metaFlag, src)
}

// 返回值第一个参数是存储的前缀信息，第二个加密的信息
func splitEncryptData(src string) (string, string) {
	str := metaCompile.FindString(src)
	return str, strings.TrimLeft(src, str)

}
