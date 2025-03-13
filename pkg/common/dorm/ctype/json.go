package ctype

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/utils"
	jsoniter "github.com/json-iterator/go"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Json struct {
	Data  interface{}
	Valid bool
}

func (j Json) Value() (driver.Value, error) {
	if !j.Valid {
		return nil, nil
	}
	bytes, err := jsoniter.Marshal(j.Data)
	if err != nil {
		return nil, err
	}
	return utils.BytesToString(bytes), nil
}

func (j *Json) Scan(src any) error {
	if src == nil {
		j.Valid = false
		return nil
	}
	switch src.(type) {
	case string:
		j.Valid = true
		return j.scanString(src.(string))
	case *string:
		j.Valid = true
		return j.scanString(*src.(*string))
	case []byte:
		j.Valid = true
		return j.scanBytes(src.([]byte))
	case *[]byte:
		j.Valid = true
		return j.scanBytes(*src.(*[]byte))
	default:
		j.Valid = false
		return errors.New(fmt.Sprintf("not support scan type for Json %T", src))
	}
}

func (j *Json) scanBytes(src []byte) error {
	initJsonData(j, src)
	return jsoniter.Unmarshal(src, &j.Data)
}
func (j *Json) scanString(src string) error {
	initJsonData(j, src)
	return jsoniter.UnmarshalFromString(src, &j.Data)

}

type sliceOrString interface {
	string | []byte
}

func initJsonData[V sliceOrString](j *Json, src V) {
	var (
		a = byte('[')
		b = byte('{')
	)
	if src[0] == a {
		j.Data = []map[string]any{}
	}
	if src[0] == b {
		j.Data = map[string]any{}
	}
}

func (j Json) MarshalJSON() ([]byte, error) {
	if !j.Valid {
		return []byte("null"), nil
	}
	return jsoniter.Marshal(j.Data)
}

func (j *Json) UnmarshalJSON(bytes []byte) error {
	initJsonData(j, bytes)
	err := jsoniter.Unmarshal(bytes, &j.Data)
	if err != nil {
		j.Valid = false
		return err
	}
	j.Valid = true
	return nil
}

func (j Json) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return dorm.JsonGormDBDataType(db, field)
}

func (j Json) GormDataType() string {
	return dorm.JsonGormDataType()
}
