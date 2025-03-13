package utils

import (
	"errors"
	"fmt"
	"github.com/duke-git/lancet/v2/slice"
	"strings"
)

func ParseBoolean(v string) (bool, error) {
	if v == "" {
		return false, errors.New("参数不能为空")
	}
	trueArr := []string{"是", "有", "1", "T", "TRUE", "Y", "YES"}
	falseArr := []string{"否", "无", "0", "F", "FALSE", "N", "NO"}
	v = strings.ToUpper(v)
	if slice.Contain(trueArr, v) {
		return true, nil
	}
	if slice.Contain(falseArr, v) {
		return false, nil
	}
	return false, errors.New(fmt.Sprintf("不能解析为布尔值: %s", v))
}
