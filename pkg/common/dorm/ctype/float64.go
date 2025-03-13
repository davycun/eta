package ctype

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"strconv"
	"strings"
)

type Float struct {
	Data      float64
	Valid     bool
	Precision int
	Scale     int
}

func (d Float) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	if !d.Valid {
		return clause.Expr{SQL: "null"}
	}
	switch dorm.GetDbType(db) {
	case dorm.Doris, dorm.Mysql:
		return clause.Expr{SQL: strconv.FormatFloat(d.Data, 'f', -1, 64)}
	}
	return clause.Expr{SQL: `?`, Vars: []interface{}{d.Data}}
}

func (d *Float) Scan(src any) error {
	if src == nil {
		d.Data = 0
		d.Valid = false
		return nil
	}
	switch x := src.(type) {
	case int8:
		d.Data = float64(x)
		d.Valid = true
	case int16:
		d.Data = float64(x)
		d.Valid = true
	case int32:
		d.Data = float64(x)
		d.Valid = true
	case int64:
		d.Data = float64(x)
		d.Valid = true
	case uint:
		d.Data = float64(x)
		d.Valid = true
	case uint8:
		d.Data = float64(x)
		d.Valid = true
	case uint16:
		d.Data = float64(x)
		d.Valid = true
	case uint32:
		d.Data = float64(x)
		d.Valid = true
	case uint64:
		d.Data = float64(x)
		d.Valid = true
	case float64:
		d.Data = x
		d.Valid = true
	case float32:
		d.Data = float64(x)
		d.Valid = true
	case string:
		float, err := strconv.ParseFloat(x, 64)
		if err != nil {
			return err
		}
		d.setScale(x)
		d.Data = float
		d.Valid = true
	case []byte:
		s := string(x)
		dt, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		d.setScale(s)
		d.Data = dt
		d.Valid = true
	default:
		d.Valid = false
	}
	return nil
}

func (d Float) Value() (driver.Value, error) {
	if d.Valid {
		return d.Data, nil
	}
	return nil, nil
}

func (d *Float) UnmarshalJSON(bytes []byte) error {
	if bytes == nil || len(bytes) < 1 {
		return nil
	}
	var f float64

	s := utils.BytesToString(bytes)
	if strings.ToLower(s) == "null" {
		d.Valid = false
		return nil
	}
	err := json.Unmarshal(bytes, &f)
	if err != nil {
		d.Valid = false
		return err
	}
	d.Data = f
	d.Valid = true
	return nil
}

func (d Float) MarshalJSON() ([]byte, error) {
	if d.Valid {
		return json.Marshal(d.Data)
	}
	return []byte("null"), nil
}

func (d Float) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	//达梦：NUMERIC[(精度 [, 标度])]，如果不指定精度和标度，缺省精度为 38，标度无限定。精度就是数据的总长度包括整数及小数。
	//PG：NUMERIC(precision, scale)，不指定precision和scale，是存储任意的precistion。
	//PG：up to 131072 digits before the decimal point; up to 16383 digits after the decimal point

	//numeric 虽然能存储非常大和非常高精度的数据，但是针对它的计算比起integer或者float来说是比较慢的。
	//但是，比如比对两个浮点数的时候，有时候不一定能得到期望的结果

	switch dorm.GetDbType(db) {
	case dorm.Doris:
		return "double"
	}

	return "numeric"
}

func (d Float) GormDataType() string {
	return "float"
}

func (d *Float) setScale(src string) {
	index := strings.Index(src, ".")
	scale := 2
	if index > -1 {
		scale = len(src[index+1:])
	}
	d.Scale = scale
}
