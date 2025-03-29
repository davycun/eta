package ctype_test

import (
	"database/sql"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/modern-go/reflect2"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type One struct {
	MyInt   int64
	MyBool  bool
	MyStr   string
	MyFloat float64
}

type St struct {
	MyIntArr      ctype.Int64Array
	MyStrArray    ctype.StringArray
	MyBool        ctype.Boolean
	MyFloat       ctype.Float
	MyInt         ctype.Integer
	MyJson        ctype.Json
	MyStr         ctype.String
	MyTime        ctype.LocalTime
	MyIntArrPrt   *ctype.Int64Array
	MyStrArrayPrt *ctype.StringArray
	MyBoolPrt     *ctype.Boolean
	MyFloatPrt    *ctype.Float
	MyIntPrt      *ctype.Integer
	MyJsonPrt     *ctype.Json
	MyStrPrt      *ctype.String
	MyTimePrt     *ctype.LocalTime

	MyOne  One `json:"my_one,omitempty" gorm:"column:my_one;serializer:json"`
	MyOne2 One `json:"my_one2,omitempty" gorm:"embedded;embeddedPrefix:h_"`
}

func TestColumnType(t *testing.T) {

	obj := St{}
	fieldType := ctype.GetColType(obj)

	tp, _ := ctype.GetFieldType(ctype.TypeArrayIntName)
	assert.Equal(t, tp, fieldType["my_int_arr"])
	tp, _ = ctype.GetFieldType(ctype.TypeArrayStringName)
	assert.Equal(t, tp, fieldType["my_str_array"])
	tp, _ = ctype.GetFieldType(ctype.TypeBoolName)
	assert.Equal(t, tp, fieldType["my_bool"])
	tp, _ = ctype.GetFieldType(ctype.TypeNumericName)
	assert.Equal(t, tp, fieldType["my_float"])
	tp, _ = ctype.GetFieldType(ctype.TypeIntegerName)
	assert.Equal(t, tp, fieldType["my_int"])
	tp, _ = ctype.GetFieldType(ctype.TypeJsonName)
	assert.Equal(t, tp, fieldType["my_json"])
	tp, _ = ctype.GetFieldType(ctype.TypeStringName)
	assert.Equal(t, tp, fieldType["my_str"])
	tp, _ = ctype.GetFieldType(ctype.TypeTimeName)
	assert.Equal(t, tp, fieldType["my_time"])

	tp, _ = ctype.GetFieldType(ctype.TypeArrayIntName)
	assert.Equal(t, tp, fieldType["my_int_arr_prt"])
	tp, _ = ctype.GetFieldType(ctype.TypeArrayStringName)
	assert.Equal(t, tp, fieldType["my_str_array_prt"])
	tp, _ = ctype.GetFieldType(ctype.TypeBoolName)
	assert.Equal(t, tp, fieldType["my_bool_prt"])
	tp, _ = ctype.GetFieldType(ctype.TypeNumericName)
	assert.Equal(t, tp, fieldType["my_float_prt"])
	tp, _ = ctype.GetFieldType(ctype.TypeIntegerName)
	assert.Equal(t, tp, fieldType["my_int_prt"])
	tp, _ = ctype.GetFieldType(ctype.TypeJsonName)
	assert.Equal(t, tp, fieldType["my_json_prt"])
	tp, _ = ctype.GetFieldType(ctype.TypeStringName)
	assert.Equal(t, tp, fieldType["my_str_prt"])
	tp, _ = ctype.GetFieldType(ctype.TypeTimeName)
	assert.Equal(t, tp, fieldType["my_time_prt"])
	tp, _ = ctype.GetFieldType(ctype.TypeIntegerName)
	assert.Equal(t, tp, fieldType["h_my_int"])
	tp, _ = ctype.GetFieldType(ctype.TypeBoolName)
	assert.Equal(t, tp, fieldType["h_my_bool"])
	tp, _ = ctype.GetFieldType(ctype.TypeStringName)
	assert.Equal(t, tp, fieldType["h_my_str"])
	tp, _ = ctype.GetFieldType(ctype.TypeNumericName)
	assert.Equal(t, tp, fieldType["h_my_float"])

}

func TestImplements(t *testing.T) {
	s1 := reflect.TypeFor[sql.Scanner]()
	s2 := reflect.TypeOf((*sql.Scanner)(nil)).Elem()

	x := &entity.BaseEntity{}
	tp := reflect2.TypeOf(x).Type1()

	if tp.Kind() == reflect.Pointer {
		tp = tp.Elem()
	}

	fd, b := tp.FieldByName("FieldUpdaterIds")
	assert.True(t, b)
	assert.True(t, fd.Type.Implements(s1))
	assert.True(t, fd.Type.Implements(s2))

	assert.False(t, fd.Type.Elem().Implements(s1))
	assert.False(t, fd.Type.Elem().Implements(s2))
}
