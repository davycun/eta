package ctype_test

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/stretchr/testify/assert"
	"reflect"
	"regexp"
	"testing"
)

func TestType(t *testing.T) {

	m := &ctype.Map{}
	v := reflect.ValueOf(m)
	assert.Equal(t, v.Kind(), reflect.Pointer)
}

func TestSet(t *testing.T) {
	assert.Equal(t, "true", fmt.Sprintf("%t", true))
	assert.Equal(t, "false", fmt.Sprintf("%t", false))
}

func TestPattern(t *testing.T) {
	pattern := `(.+)(\(([0-9]*)(,?)([0-9]*)\))`

	reg, err := regexp.Compile(pattern)
	assert.Nil(t, err)
	assert.True(t, reg.MatchString("numeric(5,3)"))
	assert.True(t, reg.MatchString("varchar(256)"))
	assert.False(t, reg.MatchString("varchar(256"))
	assert.False(t, reg.MatchString("varchar(256"))
	assert.False(t, reg.MatchString("(256)"))

	//submatch := reg.FindAllStringSubmatch("numeric(1,2", -1)
	//fmt.Println(submatch)
}
