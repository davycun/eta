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

type Integer struct {
	Data  int64
	Valid bool
}

func (d Integer) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	if !d.Valid {
		return clause.Expr{SQL: "null"}
	}
	return clause.Expr{SQL: `?`, Vars: []interface{}{d.Data}}
}
func (d *Integer) Scan(src any) error {

	switch x := src.(type) {
	case nil:
		return nil
	case int:
		d.Data = int64(x)
		d.Valid = true
	case int8:
		d.Data = int64(x)
		d.Valid = true
	case int16:
		d.Data = int64(x)
		d.Valid = true
	case int32:
		d.Data = int64(x)
		d.Valid = true
	case int64:
		d.Data = x
		d.Valid = true
	case uint:
		d.Data = int64(x)
		d.Valid = true
	case uint8:
		d.Data = int64(x)
		d.Valid = true
	case uint16:
		d.Data = int64(x)
		d.Valid = true
	case uint32:
		d.Data = int64(x)
		d.Valid = true
	case uint64:
		d.Data = int64(x)
		d.Valid = true
	case float32:
		d.Data = int64(x)
		d.Valid = true
	case float64:
		d.Data = int64(x)
		d.Valid = true
	case []byte:
		dt, err := strconv.ParseInt(string(x), 10, 64)
		if err != nil {
			return err
		}
		d.Data = dt
		d.Valid = true
	case string:
		dt, err := strconv.ParseInt(x, 10, 64)
		if err != nil {
			return err
		}
		d.Data = dt
		d.Valid = true
	}
	return nil
}

func (d Integer) Value() (driver.Value, error) {
	if d.Valid {
		return d.Data, nil
	} else {
		return nil, nil
	}
}

func (d *Integer) UnmarshalJSON(bytes []byte) error {

	if bytes == nil || len(bytes) < 1 {
		return nil
	}
	var i int64

	s := utils.BytesToString(bytes)
	if strings.ToLower(s) == "null" {
		d.Valid = false
		return nil
	}
	err := json.Unmarshal(bytes, &i)
	if err != nil {
		d.Valid = false
		return err
	}
	d.Data = i
	d.Valid = true
	return nil
}

func (d Integer) MarshalJSON() ([]byte, error) {
	if d.Valid {
		return json.Marshal(d.Data)
	}
	return []byte("null"), nil
}

func (d Integer) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	var (
		dbType = dorm.GetDbType(db)
	)
	switch dbType {
	case dorm.Doris:
		return "BIGINT"

	}
	return "integer"
}

func (d Integer) GormDataType() string {
	return "integer"
}

// MarshalBinary for go-redis
func (d *Integer) MarshalBinary() (data []byte, err error) {
	if !d.Valid {
		return make([]byte, 0, 0), nil
	}
	return utils.StringToBytes(strconv.FormatInt(d.Data, 10)), nil
}

func (d *Integer) UnmarshalText(text []byte) error {
	if text == nil || len(text) == 0 {
		d.Valid = false
		return nil
	}
	i, err := strconv.ParseInt(utils.BytesToString(text), 10, 64)
	if err != nil {
		d.Valid = false
		return nil
	}
	d.Data = i
	return nil
}
