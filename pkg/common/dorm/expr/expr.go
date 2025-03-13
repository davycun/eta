package expr

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/utils"
	"reflect"
	"strconv"
	"strings"
)

var (
	geoFunc = []string{"st_dimension", "st_coorddim", "st_geometrytype", "st_astext", "st_astext", "st_asbinary",
		"st_startpoint", "st_endpoint", "st_length", "st_perimeter", "st_numpoints", "st_centroid",
		"st_area", "st_boundary",
		"st_equals", "st_disjoint", "st_intersects", "st_touches", "st_crosses", "st_within", "st_contains", "st_overlaps", "st_distance", "st_intersection"}
)

func ExplainExprColumn(dbType dorm.DbType, column ExpColumn, tbName string) (string, error) {
	builder := strings.Builder{}
	str, err := ExplainExpr(dbType, column.Expression, tbName)
	if err != nil {
		return "", err
	}
	builder.WriteString(str)
	if column.Alias != "" {
		builder.WriteString(fmt.Sprintf(" as %s ", dorm.Quote(dbType, column.Alias)))
	}
	return builder.String(), nil
}

// ExplainExpr
// TODO 对表达式的处理，还需要处理不同数据库类型的不同函数，contains、contains_arr、json_size
func ExplainExpr(dbType dorm.DbType, exp Expression, tbName string) (string, error) {
	if exp.Vars == nil || len(exp.Vars) < 1 {
		//如果不传入tbName并且不传入参数，那么就返回原本的值...
		col := exp.Expr
		if tbName != "" {
			if col != "*" {
				return dorm.Quote(dbType, tbName, col), nil
			}
			return fmt.Sprintf("%s.*", dorm.Quote(dbType, tbName)), nil
		}
		return col, nil
	}
	bd := strings.Builder{}

	if strings.Count(exp.Expr, "?") != len(exp.Vars) {
		return "", errors.New("表达式参数与实际参数数量不匹配")
	}

	if !checkExpr(exp.Expr) {
		return "", errors.New("表达式不合法")
	}

	exp = dialectExpr(dbType, exp)

	idx := 0
	for _, v := range []byte(exp.Expr) {
		switch v {
		case '?':
			ev := exp.Vars[idx]
			switch ev.Type {
			case VarTypeColumn:
				col := fmt.Sprintf("%v", ev.Value)
				if col == "*" {
					bd.WriteString(col)
				} else {
					bd.WriteString(dorm.Quote(dbType, tbName, col))
				}
			case VarTypeValue:
				bd.WriteString(ExplainExprValue(dbType, ev.Value))
			}
			idx++
		default:
			bd.WriteByte(v)
		}
	}

	return bd.String(), nil
}

func checkExpr(exp string) bool {
	exp = strings.ToLower(exp)
	if (strings.Contains(exp, " and ") || strings.Contains(exp, " or ")) && !strings.Contains(exp, "between") {
		return false
	}
	return true
}

// 针对不同的数据库的表达式进行适配
func dialectExpr(dbType dorm.DbType, exp Expression) Expression {
	epr := strings.ToLower(strings.TrimSpace(exp.Expr))

	for _, v := range geoFunc {
		if strings.Contains(epr, v) && dbType == dorm.DaMeng {
			switch v {
			case "st_centroid", "st_boundary":
				exp.Expr = `dmgeo2.ST_AsBinary(dmgeo2.` + v + `(dmgeo2.ST_GeomFromWKB(dmgeo.ST_Asbinary(?))))`
			case "st_area":
				exp.Expr = `dmgeo2.` + v + `(dmgeo2.ST_GeomFromWKB(dmgeo.ST_Asbinary(?)))`
			default:
				exp.Expr = strings.ReplaceAll(epr, v, "dmgeo."+v)
			}
			return exp
		}
	}

	//对全文检索的支持，pg默认只支持英文的，如果需要支持中文，需要zhparser或者jieba插件，所以暂时适配成为like
	//CREATE EXTENSION zhparser;
	//CREATE INDEX idx_table_field ON table USING gin(to_tsvector('zhparser', field));
	//select * from eta_1.t_people where to_tsvector('zhparser',name) @@ to_tsquery('zhparser','李');
	if strings.HasPrefix(epr, "contains(") {
		if dbType != dorm.DaMeng {
			exp.Expr = "? like ?"
			exp.Vars[1].Value = "%" + fmt.Sprintf("%s", exp.Vars[1].Value) + "%"
			return exp
		}
	}

	//对数组检索带索引的支持
	if strings.HasPrefix(epr, "contains_arr(") {
		switch dbType {
		case dorm.PostgreSQL:
			count := strings.Count(epr, "?")
			bd := strings.Builder{}
			bd.WriteString(`? && '{`)
			for i := 0; i < count-1; i++ {
				if i > 0 {
					bd.WriteByte(',')
				}
				bd.WriteByte('?')
			}
			bd.WriteString(`}'`)
			exp.Expr = bd.String()
		case dorm.Doris:
			exp.Expr = strings.ReplaceAll(epr, "contains_arr", "array_contains")
		case dorm.DaMeng:
			exp.Expr = strings.ReplaceAll(epr, "contains_arr", "contains")
		}
	}

	//处理json_size()函数
	if strings.HasPrefix(epr, "json_size") {
		switch dbType {
		case dorm.DaMeng:
			/// 已经有函数支持 `json_size(?)`
		case dorm.PostgreSQL:
			//需要把 json_size(?) 替换为 jsonb_array_length(?)
			exp.Expr = strings.ReplaceAll(exp.Expr, "json_size", "jsonb_array_length")
		case dorm.Mysql:
			//not support yet
		}
	}

	return exp
}

func ExplainExprValue(dbType dorm.DbType, val any) string {

	var (
		builder  = strings.Builder{}
		join     = ","
		lQuote   = "("
		rQuote   = ")"
		strQuote = `'`
	)
	//TODO 其实还有一些类型没有处理完，但是按照常理说filter 条件过滤不太会出现下面之外的类型
	//不一次性处理的原因是%v采用默认的格式，会生成科学计数的方式，这种方式在数据库查询比较慢
	//case int, int8, int16, int32, int64,uint,uint8,uint16,uint32,uint64,bool,float32, float64:
	//	builder.WriteString(fmt.Sprintf("%v", val))
	switch x := val.(type) {
	case string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool, json.Number:
		return ExplainToString(dbType, "", "", "", strQuote, x)
	case []json.Number:
		//return ExplainToString(dbType, join, true, x...)
		return ExplainToString(dbType, join, lQuote, rQuote, strQuote, x...)
	case []int:
		return ExplainToString(dbType, join, lQuote, rQuote, strQuote, x...)
	case []int64:
		return ExplainToString(dbType, join, lQuote, rQuote, strQuote, x...)
	case []int32:
		return ExplainToString(dbType, join, lQuote, rQuote, strQuote, x...)
	case []int16:
		return ExplainToString(dbType, join, lQuote, rQuote, strQuote, x...)
	case []int8:
		return ExplainToString(dbType, join, lQuote, rQuote, strQuote, x...)
	case []float64:
		return ExplainToString(dbType, join, lQuote, rQuote, strQuote, x...)
	case []float32:
		return ExplainToString(dbType, join, lQuote, rQuote, strQuote, x...)
	case []string:
		return ExplainToString(dbType, join, lQuote, rQuote, strQuote, x...)
	case []interface{}:
		return ExplainToString(dbType, join, lQuote, rQuote, strQuote, x...)
	default:
		s := fmt.Sprint(val)
		WriteExprStringValue(&builder, s)
	}
	return builder.String()
}

// ExplainToString
// join 表示需要用什么字符进行拼接
// outSideQuote表示外围用什么包裹，比如要构造in条件就是括号，比如要构造Postgresql的数组就用大括号
// strQuote 表示针对t是字符串的时候用什么字符包括，如果是in条件那么就是单引号，如果是构造一个数组的表达式，那么就是双引号
// 比如：
// 'a' or 'b' or 'c'  -> 用在达梦的contains(?,?) 函数中
// {"a","b","c"} -> 用在pg的数组
func ExplainToString[T any](dbType dorm.DbType, join string, lOutSideQuote, rOutSideQuote string, strQuote string, t ...T) string {
	if len(t) < 1 {
		return ""
	}
	bd := strings.Builder{}
	bd.WriteString(lOutSideQuote)
	for i, v := range t {
		if i > 0 {
			bd.WriteString(join)
		}
		switch x := any(v).(type) {
		case int8, int16, int, int32, int64, uint8, uint16, uint, uint32, uint64:
			bd.WriteString(fmt.Sprintf("%d", x))
		case float32:
			////prec 25的原因是，float64 最大是 接近的整数是1<<63-1，有19位，在加上默认的6位小数，所以总长度是25为，采用fmt是g表示浮点prec参数是总长度
			//不一次性处理的原因是%v采用默认的格式，会生成科学计数的方式，这种方式在数据库查询比较慢
			//	builder.WriteString(fmt.Sprintf("%v", x))
			bd.WriteString(strconv.FormatFloat(float64(x), 'f', -1, 64))
		case float64:
			bd.WriteString(strconv.FormatFloat(x, 'f', -1, 64))
			//bd.WriteString(strconv.FormatFloat(x, 'g', 25, 64))
		case string:
			//is not null 或者 is null
			if strings.ToLower(x) == "null" {
				bd.WriteString(fmt.Sprintf(`%s`, x))
			} else {
				//bd.WriteString(fmt.Sprintf(`'%s'`, x))
				//bd.WriteString(fmt.Sprintf("%s%s%s", strQuote, x, strQuote))
				WriteExprStringValue(&bd, x)
			}
		case bool:
			switch dbType {
			case dorm.DaMeng:
				if x {
					bd.WriteString("1")
				} else {
					bd.WriteString("0")
				}
			default:
				bd.WriteString(strconv.FormatBool(x))
			}
		case json.Number:
			bd.WriteString(x.String())
		default:
			bd.WriteString(fmt.Sprintf("%v", x))

		}
	}
	bd.WriteString(rOutSideQuote)
	return bd.String()
}

func WriteExprStringValue(builder *strings.Builder, str string) {

	bs := utils.StringToBytes(str)
	builder.WriteString("'")
	for _, v := range bs {
		switch v {
		//防止sql注入
		case '\'':
			builder.WriteByte('\'')
		}
		builder.WriteByte(v)
	}
	builder.WriteString("'")
}

func ExplainColumnType(cols ...ExpColumn) map[string]reflect.Type {

	tp := make(map[string]reflect.Type)
	for _, v := range cols {
		nm := v.Expr
		if v.Alias != "" {
			nm = v.Alias
		}
		if nm == "" || v.Type == "" {
			continue
		}

		tmp, ok := ctype.GetType(v.Type)
		if ok {
			tp[nm] = tmp
		}
	}
	return tp
}
