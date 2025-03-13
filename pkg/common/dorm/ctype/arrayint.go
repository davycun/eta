package ctype

import (
	"context"
	dm "github.com/davycun/dm8-go-driver"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/utils"
	jsoniter "github.com/json-iterator/go"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"strconv"
	"strings"
)

type Int64Array struct {
	Data  []int64
	Valid bool
}

func (s Int64Array) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	if !s.Valid || len(s.Data) < 1 {
		return clause.Expr{SQL: "null"}
	}
	vs := make([]interface{}, len(s.Data))
	for i, v := range s.Data {
		vs[i] = v
	}

	switch dorm.GetDbType(db) {
	case dorm.PostgreSQL:
		return clause.Expr{SQL: `?`, Vars: []interface{}{s.valueStringPg()}}
	case dorm.DaMeng:
		return clause.Expr{SQL: s.valueStringDm(db, true), Vars: vs}
	case dorm.Doris, dorm.Mysql:
		return clause.Expr{SQL: `?`, Vars: []interface{}{s.valueStringDoris()}}
	}
	//if dbType == dorm.DaMeng {
	//	return clause.Expr{SQL: zero}
	//}
	return clause.Expr{SQL: "null"}
}

func (s Int64Array) valueStringPg() string {
	bd := strings.Builder{}
	bd.WriteByte('{')
	for i, v := range s.Data {
		if i > 0 {
			bd.WriteByte(',')
		}
		bd.WriteString(strconv.FormatInt(v, 10))

	}
	bd.WriteByte('}')
	return bd.String()
}
func (s Int64Array) valueStringDoris() string {
	bd := strings.Builder{}
	bd.WriteByte('[')
	for i, v := range s.Data {
		if i > 0 {
			bd.WriteByte(',')
		}
		bd.WriteString(strconv.FormatInt(v, 10))

	}
	bd.WriteByte(']')
	return bd.String()
}
func (s Int64Array) valueStringDm(db *gorm.DB, prepared bool) string {
	var (
		bd     = strings.Builder{}
		dbUser = dorm.GetDbUser(db)
	)

	bd.WriteString(dbUser + ".ARR_INT_CLS(" + dbUser + ".ARR_INT(")
	for i, v := range s.Data {
		if i > 0 {
			bd.WriteByte(',')
		}
		if prepared {
			bd.WriteByte('?')
		} else {
			bd.WriteString(strconv.FormatInt(v, 10))
		}
	}
	bd.WriteString("))")
	return bd.String()
}

func (s *Int64Array) Scan(src any) error {
	switch src.(type) {
	case nil:
		return nil
	case string:
		str := src.(string)
		return s.scanString(str)
	case *string:
		str := src.(*string)
		return s.scanString(*str)
	case *dm.DmArray:
		arr := src.(*dm.DmArray)
		return s.scanDmArray(arr)
	case *dm.DmStruct:
		st := src.(*dm.DmStruct)
		return s.scanDmStruct(st)
	}
	return nil
}

func (s *Int64Array) scanString(arr string) error {
	rs := make([]int64, 0, 20)
	src := strings.ReplaceAll(arr, "{", "")
	src = strings.ReplaceAll(src, "}", "")
	ia := strings.Split(src, ",")
	for _, v := range ia {
		i, err := strconv.ParseInt(strings.Trim(v, " "), 10, 64)
		if err != nil {
			return err
		}
		rs = append(rs, i)
	}
	s.Valid = true
	s.Data = rs
	return nil
}

func (s *Int64Array) scanDmStruct(st *dm.DmStruct) error {

	attributes, err := st.GetAttributes()
	if err != nil {
		return err
	}
	for i, v := range attributes {
		switch v.(type) {
		case *dm.DmArray:
			return s.scanDmArray(attributes[i].(*dm.DmArray))
		}
	}

	return nil
}

func (s *Int64Array) scanDmArray(arr *dm.DmArray) error {
	err := arr.Scan(&s.Data)
	array, err := arr.GetArray()
	ia, ok := array.([]interface{})
	if ok {
		s.Data = make([]int64, len(ia))
		for i, v := range ia {
			switch v.(type) {
			case int:
				s.Data[i] = int64(v.(int))
			case int32:
				s.Data[i] = int64(v.(int32))
			case int64:
				s.Data[i] = v.(int64)
			}
		}
		s.Valid = true
	}
	return err
}

func (s *Int64Array) UnmarshalJSON(bs []byte) error {

	var (
		dt  []int64
		err error
	)
	str := utils.BytesToString(bs)
	if strings.ToLower(str) == "null" {
		s.Valid = false
		s.Data = nil
		return nil
	}
	//decoder := json.NewDecoder()
	err = jsoniter.Unmarshal(bs, &dt)
	if err != nil {
		s.Valid = false
		return err
	}
	s.Valid = true
	s.Data = dt
	return nil
}

func (s Int64Array) MarshalJSON() ([]byte, error) {
	if s.Valid {
		dt, err := jsoniter.Marshal(s.Data)
		return dt, err
	}
	return []byte("null"), nil
}

func (s Int64Array) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	var (
		dbType = dorm.GetDbType(db)
		dbUser = dorm.GetDbUser(db)
	)
	switch dbType {
	case dorm.PostgreSQL:
		return "integer[]"
	case dorm.DaMeng:
		return dbUser + ".ARR_INT_CLS"
	case dorm.Mysql:
		//TODO  not yet support
		return "json"
	case dorm.Doris:
		return "ARRAY<BIGINT>"
	}
	return "integer[]"
}

func (s Int64Array) GormDataType() string {
	return TpArrayInt
}

// DataType For nebulaGraph
func (s Int64Array) DataType() string {
	return "STRING"
}
