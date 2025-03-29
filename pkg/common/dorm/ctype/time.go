package ctype

import (
	"database/sql/driver"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"reflect"
	"strings"
	"time"
)

// LocalTime
// 反序列化的时候采用的是RFC3339格式，所以必须带上时区
// 比如东八区时间：2021-01-01T00:00:00+08:00
// RFC3339     = "2006-01-02T15:04:05Z07:00"
type LocalTime struct {
	Data  time.Time
	Valid bool
}

func (t *LocalTime) UnmarshalJSON(data []byte) error {

	//return t.Data.UnmarshalJSON(data)
	if data == nil || len(data) < 1 {
		t.Valid = false
		return nil
	}

	s := strings.Trim(utils.BytesToString(data), "\"")
	s1 := strings.ToLower(s)
	if s1 == "" || s1 == "null" || s1 == "\"null\"" {
		t.Valid = false
		return nil
	}
	err := t.Data.UnmarshalJSON(data)
	if err != nil {
		t.Valid = false
		return err
	}
	t.Valid = true
	return nil
}

func (t LocalTime) MarshalJSON() ([]byte, error) {
	if t.Valid {
		return t.Data.MarshalJSON()
	}
	return nullValue, nil
}

func (t *LocalTime) Scan(value any) error {
	if value == nil {
		return nil
	}

	if v, ok := value.(time.Time); ok {
		t.Data = v
		t.Valid = true
	}
	return nil
}

// Value implements the driver Valuer interface.
func (t LocalTime) Value() (driver.Value, error) {

	if t.Valid {
		return t.Data, nil
	}
	return nil, nil
}

func (t LocalTime) GormDBDataType(db *gorm.DB, field *schema.Field) string {

	tp, err := GetDbTypeName(db, TypeTimestampTzName)
	if err != nil {
		logger.Error(err)
	}
	return tp
}

func (t *LocalTime) GormDataType() string {
	return TypeTimestampTzName
}

func (t *LocalTime) Copy(src reflect.Value) error {
	tm, ok := src.Interface().(LocalTime)
	if ok {
		t.Data = tm.Data
		t.Valid = tm.Valid
	}
	return nil
}

func (t LocalTime) IsValid() bool {
	return t.Valid
}

func (t *LocalTime) MarshalBinary() ([]byte, error) {
	if !t.Valid {
		return []byte{}, nil
	}
	return t.Data.MarshalBinary()
}

func (t *LocalTime) UnmarshalText(text []byte) error {

	if text == nil || len(text) == 0 {
		t.Valid = false
		return nil
	}
	return t.Data.UnmarshalText(text)
}
