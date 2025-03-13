package crypt_test

import (
	"github.com/davycun/eta/pkg/common/crypt"
	"github.com/davycun/eta/pkg/common/logger"
	"testing"
)

func TestStringToSlice(t *testing.T) {
	sliceSize := 3
	msg := "0123456789abcd啊哦俄一唔于"
	encStrList := crypt.StringToSlice(msg, sliceSize)
	logger.Infof("%v", encStrList)
	msg = "012"
	encStrList = crypt.StringToSlice(msg, sliceSize)
	logger.Infof("%v", encStrList)
	msg = "01"
	encStrList = crypt.StringToSlice(msg, sliceSize)
	logger.Infof("%v", encStrList)
	msg = ""
	encStrList = crypt.StringToSlice(msg, sliceSize)
	logger.Infof("%v", encStrList)
}

func TestSlice(t *testing.T) {
	msg := []string{"012", "123", "234", "345", "456", "567", "678", "789", "89a", "9ab", "abc", "bcd", "cd啊", "d啊哦", "啊哦俄", "哦俄一", "俄一唔", "一唔于"}
	runeList := crypt.SliceToString(msg)
	logger.Infof("%v", string(runeList))
	msg = []string{"012"}
	runeList = crypt.SliceToString(msg)
	logger.Infof("%v", string(runeList))
	msg = []string{"01"}
	runeList = crypt.SliceToString(msg)
	logger.Infof("%v", string(runeList))
	msg = []string{""}
	runeList = crypt.SliceToString(msg)
	logger.Infof("%v", string(runeList))
}
