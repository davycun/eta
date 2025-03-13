package validity

import (
	"database/sql/driver"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/go-playground/validator/v10"
	"reflect"
	"regexp"
	"strconv"
)

const (
	ValidateTagName = "binding"
	MobileTagName   = "mobile"
	IdCardTagName   = "id_card"
	IgnoreTagName   = "ignore"
)

var (
	// 身份证号码 权重
	idNoWeightArray = []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	// 身份证号码 校验码
	idNoCheckCode = "10X98765432"
)

func ValidateValuer(field reflect.Value) interface{} {
	if valuer, ok := field.Interface().(driver.Valuer); ok {
		val, err := valuer.Value()
		if err == nil {
			return val
		}
		// handle the error how you want
	}
	return nil
}

// MobileV
// example `validate:"required,mobile"`
func MobileV(fl validator.FieldLevel) bool {
	value := fl.Field().Interface().(string)
	if value == "" {
		return true
	}
	matched, err := regexp.MatchString(`^1[3456789][0-9]{9}$`, value)
	if err != nil {
		logger.Errorf("手机号校验出错: %s", err.Error())
	}
	return matched
}

// IdCardV
// example `validate:"required,id_card"`
func IdCardV(fl validator.FieldLevel) bool {
	value := fl.Field().Interface().(string)
	return IdCard(value)
}

// IdCard
func IdCard(value string) bool {
	if value == "" {
		return true
	}
	//match, _ := regexp.MatchString(`^[1-9]\d{5}(19|20)\d{2}((0[1-9])|(10|11|12))((0[1-9])|([1-2][0-9])|30|31)\d{3}([0-9Xx])$`, value)
	matched, err := regexp.MatchString(`(^\d{15}$)|(^\d{18}$)|(^\d{17}(\d|X|x)$)`, value)
	if err != nil {
		logger.Errorf("身份证号码校验出错: %s", err.Error())
		return false
	}
	// 18位身份证，校验码
	if matched && len(value) == 18 {
		data := value[0:17]
		s := 0
		for i, _ := range data {
			n, _ := strconv.Atoi(string(data[i]))
			s += n * idNoWeightArray[i]
		}
		y := s % 11
		matched = idNoCheckCode[y:y+1] == value[17:18]
	}

	return matched
}

func Ignore(fl validator.FieldLevel) bool {

	return true
}
