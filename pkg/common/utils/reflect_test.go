package utils_test

import (
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestPrtToReal(t *testing.T) {

	type My struct {
		A int
		B string
		C bool
	}

	var (
		a  = 1
		b  = &a
		c  = "5"
		d  = &c
		m1 = My{A: 1, B: "5", C: true}
		m2 = &m1
	)
	assert.Equal(t, 1, utils.PrtToReal(a))
	assert.Equal(t, 1, utils.PrtToReal(b))
	assert.Equal(t, "5", utils.PrtToReal(c))
	assert.Equal(t, "5", utils.PrtToReal(d))

	m11 := utils.PrtToReal(m1).(My)
	m21 := utils.PrtToReal(m2).(My)

	assert.True(t, m11.C)
	assert.True(t, m21.C)
}

func TestNewPointer(t *testing.T) {

	type My struct {
		A int
		B string
		C bool
	}
	tp := reflect.TypeOf(My{})

	x := utils.NewPointer(tp, false)
	_, ok := x.(*My)
	assert.True(t, ok)

	x = utils.NewPointer(tp, true)
	_, ok = x.(*[]My)
	assert.True(t, ok)

	x = utils.NewPointer(reflect.TypeOf(int64(1)), false)
	_, ok = x.(*int64)
	assert.True(t, ok)

	x = utils.NewPointer(reflect.TypeOf("abc"), true)
	_, ok = x.(*[]string)
	assert.True(t, ok)

}
