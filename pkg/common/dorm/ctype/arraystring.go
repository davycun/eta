package ctype

import (
	"context"
	"fmt"
	dm "github.com/davycun/dm8-go-driver"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/utils"
	jsoniter "github.com/json-iterator/go"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"strings"
)

type StringArray struct {
	Data  []string
	Valid bool
}

func (s StringArray) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	if !s.Valid || len(s.Data) < 1 {
		return clause.Expr{SQL: "null"}
	}
	vs := make([]interface{}, len(s.Data))
	for i, v := range s.Data {
		vs[i] = v
	}
	switch dorm.GetDbType(db) {
	case dorm.PostgreSQL:
		return clause.Expr{SQL: `?`, Vars: []interface{}{s.valueString(true)}}
	case dorm.DaMeng:
		return clause.Expr{SQL: s.valueStringDm(db, true), Vars: vs}
	case dorm.Doris, dorm.Mysql:
		return clause.Expr{SQL: `?`, Vars: []interface{}{s.valueStringDoris(true)}}
	}
	return clause.Expr{SQL: "null"}
}

func (s StringArray) valueString(quote bool) string {
	bd := strings.Builder{}
	bd.WriteByte('{')
	for i, v := range s.Data {
		if i > 0 {
			bd.WriteByte(',')
		}
		if quote {
			bd.WriteByte('"')
		}
		bd.WriteString(v)
		if quote {
			bd.WriteByte('"')
		}
	}
	bd.WriteByte('}')
	return bd.String()
}
func (s StringArray) valueStringDoris(quote bool) string {
	bd := strings.Builder{}
	bd.WriteByte('[')
	for i, v := range s.Data {
		if i > 0 {
			bd.WriteByte(',')
		}
		if quote {
			bd.WriteString(`"`)
		}
		bd.WriteString(v)
		if quote {
			bd.WriteString(`"`)
		}
	}
	bd.WriteByte(']')
	return bd.String()
}
func (s StringArray) valueStringDm(db *gorm.DB, prepare bool) string {
	//scm := dorm.GetDbUser(db)
	var (
		bd     = strings.Builder{}
		dbUser = dorm.GetDbUser(db)
	)
	bd.WriteString(dbUser + ".ARR_STR_CLS(" + dbUser + ".ARR_STR(")
	for i, v := range s.Data {
		if i > 0 {
			bd.WriteByte(',')
		}
		if prepare {
			bd.WriteByte('?')
		} else {
			bd.WriteString(v)
		}
	}
	bd.WriteString("))")
	return bd.String()
}

func (s *StringArray) Scan(src any) error {
	switch x := src.(type) {
	case nil:
		return nil
	case string:
		str := src.(string)
		return s.scan(str)
	case *string:
		str := src.(*string)
		return s.scan(*str)
	case *dm.DmArray:
		arr := src.(*dm.DmArray)
		return s.scanDmArray(arr)
	case *dm.DmStruct:
		arr := src.(*dm.DmStruct)
		return s.scanDmStruct(arr)
	case []byte:
		str := string(x)
		return s.scan(str)
	default:
		logger.Errorf("not yet support for type %T scan to StringArray", src)
	}
	return nil
}

func (s *StringArray) scan(arr string) error {

	var (
		pre        int32
		quoteStart = false
		rs         = make([]string, 0, 20)
	)
	bd := strings.Builder{}
	for i, v := range arr {
		// "==34  \==92 ,==44  {==123  }==125 [==91  ]==93
		if v == 34 {
			if !quoteStart {
				quoteStart = true
				pre = v
				continue
			} else {
				if pre != 92 {
					quoteStart = false
					if bd.Len() > 0 {
						rs = append(rs, strings.TrimSpace(bd.String()))
						bd.Reset()
					}
					pre = v
					continue
				}
			}
		}
		//第一个{
		if (v == 123 || v == 91) && i == 0 {
			pre = v
			continue
		}
		if v == 125 && i == (len(arr)-1) {
			if bd.Len() > 0 {
				rs = append(rs, strings.TrimSpace(bd.String()))
				bd.Reset()
			}
			pre = v
			continue
		}
		if v == 44 && !quoteStart {
			if bd.Len() > 0 {
				rs = append(rs, strings.TrimSpace(bd.String()))
				bd.Reset()
			}
			pre = v
			continue
		}
		bd.WriteRune(v)
		pre = v
	}
	s.Valid = true
	s.Data = rs
	return nil
}

func (s *StringArray) scanDmStruct(st *dm.DmStruct) error {

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

func (s *StringArray) scanDmArray(arr *dm.DmArray) error {
	//err := arr.Scan(&s.Data)
	array, err := arr.GetArray()
	ia, ok := array.([]interface{})
	if ok {
		s.Data = make([]string, len(ia))
		for i, v := range ia {
			s.Data[i] = strings.TrimSpace(fmt.Sprintf("%v", v))
		}
		s.Valid = true
	}
	return err
}

func (s *StringArray) UnmarshalJSON(bs []byte) error {

	var (
		dt  []string
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

func (s StringArray) MarshalJSON() ([]byte, error) {
	if s.Valid {
		dt, err := jsoniter.Marshal(s.Data)
		return dt, err
	}
	return nullValue, nil
}

func (s StringArray) GormDBDataType(db *gorm.DB, field *schema.Field) string {

	tp, err := GetDbTypeName(db, TypeArrayStringName)
	if err != nil {
		logger.Error(err)
	}
	return tp
}

func (s StringArray) GormDataType() string {
	return TypeArrayStringName
}

// DataType For nebulaGraph
func (s StringArray) DataType() string {
	return "STRING"
}
