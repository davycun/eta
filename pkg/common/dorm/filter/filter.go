package filter

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/expr"
	"github.com/davycun/eta/pkg/common/logger"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"strings"
)

var (
	Eq    = "="
	Neq   = "!="
	GT    = ">"
	GTE   = ">="
	LT    = "<"
	LTE   = "<="
	IN    = "in"
	NotIn = "not in"
	Like  = "like"
	IS    = "is"     // 是否为 xx
	IsNot = "is not" // 不为 xx

	And = "and"
	Or  = "or"

	operator          = []string{Eq, Neq, GT, GTE, LT, LTE, IN, NotIn, Like, IS, IsNot}
	illegalColumnChar = []string{"\"", " ", "(", ")", ">", "=", "<", ","}
	illegalExprChar   = []string{"\"", ">", "=", "<"}
)

type Filter struct {
	//and、or
	LogicalOperator string          `json:"logical_operator,omitempty" binding:"oneof=and or AND OR ''"`
	Column          string          `json:"column,omitempty"`
	Operator        string          `json:"operator,omitempty"`
	Value           any             `json:"value,omitempty"`
	Filters         []Filter        `json:"filters,omitempty"`
	Expr            expr.Expression `json:"expr,omitempty"`
}

// Filters 为了解决达梦和pg之间类型区别
type Filters []Filter

func (d Filters) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return dorm.JsonGormDBDataType(db, field)
}

func (d Filters) GormDataType() string {
	return dorm.JsonGormDataType()
}

func OperatorNegation(opt string) string {
	switch strings.TrimSpace(strings.ToLower(opt)) {
	case Eq:
		return Neq
	case Neq:
		return Eq

	case GT:
		return LT
	case GTE:
		return LTE
	case LT:
		return GT
	case LTE:
		return GTE
	case IN:
		return NotIn
	case NotIn:
		return IN
	case IS:
		return IsNot
	case IsNot:
		return IS
	}
	return Neq
}

func ValidateColumnName(name string) bool {
	if name == "" {
		return true
	}
	for _, v := range illegalColumnChar {
		if strings.Contains(name, v) {
			return false
		}
	}
	return true
}
func ValidateExpr(name string) bool {
	if name == "" {
		return true
	}
	//for _, v := range illegalExprChar {
	//	if strings.Contains(name, v) {
	//		return false
	//	}
	//}
	return true
}
func ValidateOperator(opt string) bool {
	if opt == "" {
		return true
	}
	for _, v := range operator {
		t := strings.ToLower(strings.TrimSpace(opt))
		if t == v {
			return true
		}
	}
	return false
}

func ResolveWhere(cds []Filter, dbType dorm.DbType) string {
	return ResolveWhereTable("", cds, dbType)
}
func ResolveWhereTable(tableName string, cds []Filter, dbType dorm.DbType) string {
	if cds == nil || len(cds) < 1 {
		return ""
	}
	_, tableName = dorm.SplitSchemaTableName(tableName)
	sq := buildConditions(dbType, tableName, cds)
	if sq != "" {
		return " (" + sq + ") "
	}
	return sq
}
func JoinFilters(filters ...[]Filter) []Filter {
	tmp := make([]Filter, 0, 10)
	if filters == nil {
		return tmp
	}
	for _, v := range filters {
		tmp = append(tmp, v...)
	}
	return tmp
}

func writeWhere(b *strings.Builder, opt string, wh string) {

	if wh == "" {
		return
	}
	if opt == "" {
		opt = And
	}
	if b.Len() < 1 {
		b.WriteString(wh)
	} else {
		b.WriteString(" " + opt + " ")
		b.WriteString(wh)
	}
}
func writeChildWhere(b *strings.Builder, opt string, wh string) {
	if wh == "" {
		return
	}
	if opt == "" {
		opt = And
	}
	if b.Len() < 1 {
		b.WriteString(" (" + wh + ") ")
	} else {
		b.WriteString(" " + opt + " (" + wh + ") ")
	}
}

func buildConditions(dbType dorm.DbType, tableName string, c []Filter) string {
	if c == nil || len(c) < 1 {
		return ""
	}
	builder := strings.Builder{}
	for _, v := range c {
		if v.LogicalOperator == "" {
			v.LogicalOperator = And
		}
		if v.Expr.Expr == "" && v.Column == "" {
			goto childFilters
		}
		if !ValidateOperator(v.Operator) || !ValidateColumnName(v.Column) || !ValidateExpr(v.Expr.Expr) {
			goto childFilters
		}

		if v.Expr.Expr != "" {
			writeWhere(&builder, v.LogicalOperator, buildExprCondition(dbType, v, tableName))
		} else if v.Column != "" {
			writeWhere(&builder, v.LogicalOperator, buildColumnCondition(dbType, v, tableName))
		}
	childFilters:
		if len(v.Filters) > 0 {
			cds := buildConditions(dbType, tableName, v.Filters)
			if cds != "" {
				writeChildWhere(&builder, v.LogicalOperator, cds)
			}
		}
	}
	return builder.String()
}

func buildColumnCondition(dbType dorm.DbType, v Filter, tableName string) string {
	if !ValidateColumnName(v.Column) || v.Column == "" || v.Operator == "" {
		return ""
	}
	builder := strings.Builder{}
	if v.Value == nil {
		v.Value = "null"
	}
	var (
		val = expr.ExplainExprValue(dbType, v.Value)
	)
	//if v.Operator != Eq && (val == `''` || val == "") {
	//	return ""
	//}
	if tableName != "" {
		builder.WriteString(fmt.Sprintf(` %s.%s `, dorm.Quote(dbType, tableName), dorm.Quote(dbType, v.Column)))
	} else {
		builder.WriteString(fmt.Sprintf(` %s `, dorm.Quote(dbType, v.Column)))
	}

	builder.WriteString(fmt.Sprintf(` %s `, v.Operator))
	builder.WriteString(val)
	return builder.String()
}

// TODO json_size、json_contains_any、arr_contains_str、arr_contains_int，这个四个函数是为达梦编写的。所以下面的需要进行兼容
// contains 函数是达梦针对全文检索及数组索引的函数
func buildExprCondition(dbType dorm.DbType, v Filter, tbName string) string {
	if !ValidateExpr(v.Expr.Expr) {
		return ""
	}
	fb := strings.Builder{}
	explainExpr, err := expr.ExplainExpr(dbType, v.Expr, tbName)
	if err != nil {
		logger.Errorf("构建filter出错: %s", err)
	} else {
		fb.WriteString(`(` + explainExpr + `)`)
		if v.Operator != "" && v.Value != nil {
			fb.WriteString(` ` + v.Operator + ` `)
			fb.WriteString(expr.ExplainExprValue(dbType, v.Value))
		}
	}
	return fb.String()
}
