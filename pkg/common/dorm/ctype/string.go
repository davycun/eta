package ctype

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"strings"
)

type String struct {
	Data  string
	Valid bool
}

func (s String) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	if !s.Valid {
		return clause.Expr{SQL: "null"}
	}
	return clause.Expr{SQL: `?`, Vars: []interface{}{s.Data}}
}

func (s *String) Scan(src any) error {

	if sc, ok := GetScanner(TypeStringName); ok {
		val, err := sc(src)
		s.Data = ToString(val)
		s.Valid = true
		return err
	}
	switch src.(type) {
	case nil:
		return nil
	case string:
		s.Data = src.(string)
		s.Valid = true
	case *string:
		s.Data = *src.(*string)
		s.Valid = true
	case []byte:
		s.Data = string(src.([]byte))
		s.Valid = true
	}

	return nil
}

func (s String) Value() (driver.Value, error) {
	if s.Valid {
		return s.Data, nil
	} else {
		return nil, nil
	}
}

func (s *String) UnmarshalJSON(bytes []byte) error {

	if bytes == nil || len(bytes) < 1 {
		return nil
	}
	var str string

	st := utils.BytesToString(bytes)
	if strings.ToLower(st) == "null" {
		s.Valid = false
		return nil
	}
	err := json.Unmarshal(bytes, &str)
	if err != nil {
		s.Valid = false
		return err
	}
	s.Data = str
	s.Valid = true
	return nil
}

func (s String) MarshalJSON() ([]byte, error) {
	if s.Valid {
		return json.Marshal(s.Data)
	}
	return nullValue, nil
}

func (s String) GormDBDataType(db *gorm.DB, field *schema.Field) string {

	tp, err := GetDbTypeName(db, TypeStringName)
	if err != nil {
		logger.Error(err)
	}
	return tp
}

func (s String) GormDataType() string {
	return TypeStringName
}

// MarshalBinary for go-redis
func (s *String) MarshalBinary() (data []byte, err error) {
	if !s.Valid {
		return make([]byte, 0, 0), nil
	}
	return utils.StringToBytes(s.Data), nil
}

func (s *String) UnmarshalText(text []byte) error {

	if text == nil || len(text) < 1 {
		return nil
	}
	st := utils.BytesToString(text)
	if strings.TrimSpace(strings.ToLower(st)) == "null" {
		return nil
	}
	s.Data = st
	s.Valid = true
	return nil
}
