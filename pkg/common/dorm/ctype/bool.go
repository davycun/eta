package ctype

import (
	"context"
	"encoding/json"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"strings"
)

type Boolean struct {
	Data  bool
	Valid bool
}

func (b Boolean) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	if !b.Valid {
		return clause.Expr{SQL: "null"}
	}

	switch dorm.GetDbType(db) {
	case dorm.PostgreSQL, dorm.Mysql, dorm.Doris:
		return clause.Expr{SQL: `?`, Vars: []interface{}{b.Data}}
	case dorm.DaMeng:
		i := 0
		if b.Data {
			i = 1
		}
		return clause.Expr{SQL: `?`, Vars: []interface{}{i}}
	}
	return clause.Expr{SQL: "null"}
}

func (b *Boolean) Scan(src any) error {
	if src == nil {
		b.Valid = false
		return nil
	}
	switch src.(type) {
	case bool:
		v := src.(bool)
		b.Valid = true
		b.Data = v
	case int8: //for dameng
		v := src.(int8)
		b.Valid = true
		if v == 1 {
			b.Data = true
		} else {
			b.Data = false
		}
	case int64: //for mysql
		v := src.(int64)
		b.Valid = true
		if v == 1 {
			b.Data = true
		} else {
			b.Data = false
		}
	}
	return nil
}

func (b Boolean) MarshalJSON() ([]byte, error) {
	if b.Valid {
		return json.Marshal(b.Data)
	}
	return nullValue, nil
}

func (b *Boolean) UnmarshalJSON(bytes []byte) error {

	if bytes == nil || len(bytes) < 1 {
		return nil
	}
	var f bool

	s := utils.BytesToString(bytes)
	if strings.ToLower(s) == "null" {
		b.Valid = false
		return nil
	}
	err := json.Unmarshal(bytes, &f)
	if err != nil {
		b.Valid = false
		return err
	}
	b.Data = f
	b.Valid = true
	return nil
}

func (b Boolean) GormDBDataType(db *gorm.DB, field *schema.Field) string {

	tp, err := GetDbTypeName(db, TypeBoolName)
	if err != nil {
		logger.Error(err)
	}
	return tp
}

func (b Boolean) GormDataType() string {
	return TypeBoolName
}

func (b Boolean) Value() bool {
	return b.Valid && b.Data
}
