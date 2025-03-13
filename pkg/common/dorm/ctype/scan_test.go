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

	assert.Equal(t, ctype.FieldType[ctype.TpArrayInt], fieldType["my_int_arr"])
	assert.Equal(t, ctype.FieldType[ctype.TpArrayString], fieldType["my_str_array"])
	assert.Equal(t, ctype.FieldType[ctype.TpBool], fieldType["my_bool"])
	assert.Equal(t, ctype.FieldType[ctype.TpNumeric], fieldType["my_float"])
	assert.Equal(t, ctype.FieldType[ctype.TpInteger], fieldType["my_int"])
	assert.Equal(t, ctype.FieldType[ctype.TpJson], fieldType["my_json"])
	assert.Equal(t, ctype.FieldType[ctype.TpString], fieldType["my_str"])
	assert.Equal(t, ctype.FieldType[ctype.TpTime], fieldType["my_time"])

	assert.Equal(t, ctype.FieldType[ctype.TpArrayInt], fieldType["my_int_arr_prt"])
	assert.Equal(t, ctype.FieldType[ctype.TpArrayString], fieldType["my_str_array_prt"])
	assert.Equal(t, ctype.FieldType[ctype.TpBool], fieldType["my_bool_prt"])
	assert.Equal(t, ctype.FieldType[ctype.TpNumeric], fieldType["my_float_prt"])
	assert.Equal(t, ctype.FieldType[ctype.TpInteger], fieldType["my_int_prt"])
	assert.Equal(t, ctype.FieldType[ctype.TpJson], fieldType["my_json_prt"])
	assert.Equal(t, ctype.FieldType[ctype.TpString], fieldType["my_str_prt"])
	assert.Equal(t, ctype.FieldType[ctype.TpTime], fieldType["my_time_prt"])

	assert.Equal(t, ctype.FieldType[ctype.TpInteger], fieldType["h_my_int"])
	assert.Equal(t, ctype.FieldType[ctype.TpBool], fieldType["h_my_bool"])
	assert.Equal(t, ctype.FieldType[ctype.TpString], fieldType["h_my_str"])
	assert.Equal(t, ctype.FieldType[ctype.TpNumeric], fieldType["h_my_float"])

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
