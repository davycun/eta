package validity_test

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/common/validity"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

var (
	mobiels = map[string]bool{
		"13917182633":   true,
		"14917182633":   true,
		"03917182633":   false,
		"23917182633":   false,
		"12917182633":   false,
		"159171823":     false,
		"1391712554633": false,
		"11917182633":   false,
	}
)

func TestTagParam(t *testing.T) {

	println(string(rune(0x2C)))
	println(string(rune(0x7C)))
	compile := regexp.MustCompile(`'[^']*'|\S+`)

	as := compile.FindAllString("abc ps cs", -1)

	for _, v := range as {
		println(v)
	}

}

func TestMobile(t *testing.T) {
	for k, v := range mobiels {
		matched, err := regexp.MatchString(`^1[345678][0-9]{9}$`, k)
		assert.Nil(t, err)
		assert.Equal(t, matched, v)
	}
	//validate := validator.New()
	type TMobile struct {
		Test1 string  `validate:"required,mobile"`
		Test2 string  `validate:"mobile"`
		Test3 *string `validate:"mobile"`
	}
	v := "23r43243"
	tm := TMobile{
		Test1: "13234324223",
		Test2: "",
		Test3: &v,
	}
	val := validator.New()
	val.RegisterValidation(validity.MobileTagName, validity.MobileV)
	errs := val.StructExcept(tm)
	if errs != nil {
		fmt.Println(errs.Error())
	}
	println("23r43243" == *tm.Test3)
}

func TestStructValidation_SkipsNotGivenNestedStructs(t *testing.T) {
	type Inner struct {
		Test string `validate:"required"`
	}
	type TestStruct struct {
		//In Inner `validate:"-"`
		In Inner `validate:"nostructlevel"`
		//In Inner `validate:"structonly"`
	}

	ts := TestStruct{}
	val := validator.New()

	fn1 := func(fl validator.FieldLevel) bool {
		println("asdfasdfasdfasdfsadf")
		return fl.Field().FieldByName("Test").String() == "some value"
	}
	fn2 := func(fl validator.FieldLevel) bool {
		return fl.Field().FieldByName("Test").String() == "another value"
	}

	val.RegisterValidation("one", fn1)
	val.RegisterValidation("two", fn2)
	errs := val.StructExcept(ts)
	logger.Infof("errors:%v", errs)
	//assert.Equal(t, len(errs.(validator.ValidationErrors)), 1)
	//assert.AssertError(t, errs, "TestStruct.Test", "TestStruct.Test", "Test", "Test", "one")
}

func TestIdCard(t *testing.T) {
	match := validity.IdCard("469002199908011189")
	assert.True(t, match)
	match = validity.IdCard("469002199908011182")
	assert.False(t, match)

	birthDayStr := "469002199908011189"[6:14]
	logger.Infof("birthDayStr:%s", birthDayStr)
	time, err := utils.FormatStrToTime(birthDayStr)
	assert.NoError(t, err)
	logger.Infof("birthDay:%v", time)
}
