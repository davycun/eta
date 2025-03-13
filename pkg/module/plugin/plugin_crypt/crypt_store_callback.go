package plugin_crypt

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/crypt"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/duke-git/lancet/v2/slice"
	"reflect"
	"strings"
	"sync"
)

var (
	//TODO 需要定时清理
	encryptCache = sync.Map{}
)

func CleanEncryptCache() {
	encryptCache = sync.Map{}
}

// StoreEncrypt 在数据插入和修改前，加密
func StoreEncrypt(cfg *hook.SrvConfig, pos hook.CallbackPosition) (err error) {
	if pos != hook.CallbackBefore || (cfg.Method != iface.MethodCreate && cfg.Method != iface.MethodUpdate && cfg.Method != iface.MethodUpdateByFilters) {
		return
	}
	table := entity.GetContextTable(cfg.Ctx)
	if len(table.CryptFields) < 1 {
		return
	}
	for _, v := range cfg.Values {
		for _, cryptInfo := range table.CryptFields {
			if !cryptInfo.Enable || cryptInfo.Field == "" {
				continue
			}
			err = encryptValue(v, cryptInfo)
			if err != nil {
				return err
			}
		}
	}
	return
}
func encryptValue(val reflect.Value, cryptFieldInfo entity.CryptFieldInfo) error {
	if !val.IsValid() || val.IsZero() || !val.CanInterface() {
		return nil
	}
	valInter := val.Interface()
	fieldStr := entity.GetString(valInter, cryptFieldInfo.Field)
	if fieldStr == "" {
		return nil
	}
	rs, ok := encryptData(cryptFieldInfo, []rune(fieldStr))
	if rs == "" || !ok {
		return nil
	}
	return entity.SetValue(val, cryptFieldInfo.Field, rs)
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
			key := fmt.Sprintf("%s@%s", cryptInfo.SecretKey[0], plaintext)
			if rsData, ok := encryptCache.Load(key); ok {
				return rsData.(string)
			}
			rs, err := crypt.EncryptBase64(cryptInfo.Algo, cryptInfo.SecretKey[0], plaintext)
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

// StoreDecrypt 在数据查询后，解密
func StoreDecrypt(cfg *hook.SrvConfig, pos hook.CallbackPosition) error {
	return hook.AfterRetrieve(cfg, pos, func(cfg *hook.SrvConfig, rs []reflect.Value) error {
		table := cfg.GetTable()
		if len(table.CryptFields) < 1 {
			return nil
		}
		for _, val := range rs {
			for _, cryptInfo := range table.CryptFields {
				if !cryptInfo.Enable || cryptInfo.Field == "" {
					continue
				}
				err := decryptValue(val, cryptInfo)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
}
func decryptValue(val reflect.Value, cryptFieldInfo entity.CryptFieldInfo) error {

	if !val.IsValid() || val.IsZero() || !val.CanInterface() {
		return nil
	}
	valInter := val.Interface()
	fieldStr := entity.GetString(valInter, cryptFieldInfo.Field)
	rsStr, ok := decryptData(cryptFieldInfo, []rune(fieldStr))
	if rsStr == "" || !ok {
		return nil
	}
	return entity.SetValue(val, cryptFieldInfo.Field, rsStr)
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
			rs, err := crypt.DecryptBase64(cryptInfo.Algo, cryptInfo.SecretKey[0], plaintext)
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
