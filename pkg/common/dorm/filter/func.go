package filter

import (
	"encoding/json"
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/expr"
	"github.com/davycun/eta/pkg/common/utils"
	"strings"
)

// ResolveContainsWhere TODO 兼容PG部分等待测试
func ResolveContainsWhere(dbType dorm.DbType, not bool, cds ...Filter) string {
	return ResolveContainsWhereTable(dbType, "", not, cds...)
}
func ResolveContainsWhereTable(dbType dorm.DbType, tableName string, not bool, cds ...Filter) string {
	if len(cds) < 1 {
		return ""
	}
	var (
		bd     = strings.Builder{}
		hasPre = false
	)

	for _, v := range cds {
		opt := v.LogicalOperator
		if opt == "" {
			opt = And
		}
		str := ExplainContainsValue(dbType, tableName, v.Column, not, v.Value)
		if str != "" {
			if hasPre {
				bd.WriteString(fmt.Sprintf(` %s `, opt))
			}
			bd.WriteString(str)
			hasPre = true
		}
		if len(v.Filters) > 0 {
			s := ResolveContainsWhereTable(dbType, tableName, not, v.Filters...)
			if s != "" {
				bd.WriteString(fmt.Sprintf(` %s (%s) `, opt, s))
			}
		}
	}
	return bd.String()
}

// ExplainContainsValue
// 主要是处理达梦的Contains函数
func ExplainContainsValue(dbType dorm.DbType, tableName string, column string, not bool, val any) string {

	switch x := val.(type) {
	case string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool, json.Number:
		return ExplainContainsAnyElement(dbType, tableName, column, not, val)
	case []json.Number:
		return ExplainContainsAnyElement(dbType, tableName, column, not, x...)
	case []int:
		return ExplainContainsAnyElement(dbType, tableName, column, not, x...)
	case []int64:
		return ExplainContainsAnyElement(dbType, tableName, column, not, x...)
	case []int32:
		return ExplainContainsAnyElement(dbType, tableName, column, not, x...)
	case []int16:
		return ExplainContainsAnyElement(dbType, tableName, column, not, x...)
	case []int8:
		return ExplainContainsAnyElement(dbType, tableName, column, not, x...)
	case []float64:
		return ExplainContainsAnyElement(dbType, tableName, column, not, x...)
	case []float32:
		return ExplainContainsAnyElement(dbType, tableName, column, not, x...)
	case []string:
		return ExplainContainsAnyElement(dbType, tableName, column, not, x...)
	case []bool:
		return ExplainContainsAnyElement(dbType, tableName, column, not, x...)
	case []interface{}:
		return ExplainContainsAnyElement(dbType, tableName, column, not, x...)
	default:
		bd := strings.Builder{}
		s := fmt.Sprint(val)
		if s != "" {
			expr.WriteExprStringValue(&bd, s)
		}
		return bd.String()
	}
}

func ExplainContainsAnyElement[T any](dbType dorm.DbType, tableName string, column string, not bool, val ...T) string {

	if len(val) < 1 {
		return ""
	}

	var (
		no            = ""
		bd            = strings.Builder{}
		tableColQuote = dorm.Quote(dbType, column)
		join          = "or"
		lQuote        = "("
		rQuote        = ")"
		strQuote      = `'`
	)
	if tableName != "" {
		tableColQuote = dorm.Quote(dbType, tableName) + "." + tableColQuote
	}
	if not {
		no = "not"
	}

	switch dbType {
	case dorm.PostgreSQL:
		join = ","
		lQuote = "{"
		rQuote = "}"
		strQuote = `"`
	}

	switch dbType {
	case dorm.DaMeng:
		// where not contains("live_ids",'a' or 'b')
		valStr := expr.ExplainToString(dbType, join, lQuote, rQuote, strQuote, val...)
		if valStr != "" {
			bd.WriteString(fmt.Sprintf(`%s contains(%s,%s) `, no, tableColQuote, valStr))
		}
	case dorm.PostgreSQL:
		// where not ("live_ids" && '{"a","b"}')
		valStr := expr.ExplainToString(dbType, join, lQuote, rQuote, strQuote, val...)
		if valStr != "" {
			bd.WriteString(fmt.Sprintf(`%s (%s && '%s')`, no, tableColQuote, valStr))
		}
	case dorm.Doris:

		//where ( not array_contains(`live_ids`,'a') or not array_contains(`live_ids`,'b'))
		if len(val) > 1 {
			bd.WriteString("(")
		}
		cs := make([]string, 0, len(val))
		for _, v := range val {
			tmp := expr.ExplainToString(dbType, "", "", "", `'`, v)
			if tmp != "" {
				cs = append(cs, fmt.Sprintf(`array_contains(%s,%s)`, tableColQuote, tmp))
			}
		}
		bd.WriteString(fmt.Sprintf(`%s %s`, no, strings.Join(cs, " or ")))
		if len(val) > 1 {
			bd.WriteString(")")
		}
	case dorm.Mysql:
		//TODO not yet support

	}

	return bd.String()
}
func BuildOrContains[T string | int64 | int32 | int](dbType dorm.DbType, fieldName string, dt ...T) string {

	if len(dt) < 1 {
		return ""
	}

	return ExplainContainsAnyElement(dbType, "", fieldName, false, dt...)
}

func KeywordToFilter(col string, keyword string) []Filter {
	rs := make([]Filter, 0, 1)
	if keyword != "" {
		scs := utils.Split(keyword, ",", "，", " ")
		for _, v := range scs {
			flt := Filter{}
			flt.Expr.Expr = "contains(?,?)"
			flt.Expr.Vars = []expr.ExpVar{
				{Type: expr.VarTypeColumn, Value: col},
				{Type: expr.VarTypeValue, Value: v},
			}
			rs = append(rs, flt)
		}
	}

	return rs
}
