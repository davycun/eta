package filter

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/expr"
	"strings"
)

type Having struct {
	LogicalOperator string   `json:"logical_operator,omitempty" binding:"oneof=and or AND OR ''"`
	Column          string   `json:"column,omitempty" binding:"required"`
	AggFunc         string   `json:"agg_func,omitempty" binding:"required"`
	Operator        string   `json:"operator,omitempty" binding:"required"`
	Value           any      `json:"value,omitempty"`
	Having          []Having `json:"having,omitempty"`
}

func ResolveHaving(dbType dorm.DbType, hvs ...Having) string {
	return ResolveHavingTable(dbType, "", hvs...)
}
func ResolveHavingTable(dbType dorm.DbType, tableName string, hvs ...Having) string {
	if len(hvs) < 1 {
		return ""
	}
	_, tableName = dorm.SplitSchemaTableName(tableName)
	sq := buildHaving(dbType, tableName, hvs)
	if sq != "" {
		return " (" + sq + ") "
	}
	return sq
}

func buildHaving(dbType dorm.DbType, tableName string, c []Having) string {
	if c == nil || len(c) < 1 {
		return ""
	}
	builder := strings.Builder{}
	for _, v := range c {
		if v.LogicalOperator == "" {
			v.LogicalOperator = And
		}
		if !ValidateOperator(v.Operator) {
			goto childFilters
		}
		writeWhere(&builder, v.LogicalOperator, buildHavingCondition(dbType, v, tableName))

	childFilters:
		if len(v.Having) > 0 {
			cds := buildHaving(dbType, tableName, v.Having)
			if cds != "" {
				writeChildWhere(&builder, v.LogicalOperator, cds)
			}
		}
	}
	return builder.String()
}

func buildHavingCondition(dbType dorm.DbType, v Having, tableName string) string {
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
	if v.Operator != Eq && (val == `''` || val == "") {
		return ""
	}
	col := v.Column
	if col != "*" {
		col = dorm.Quote(dbType, col)
	}
	if tableName != "" && col != "*" {
		builder.WriteString(fmt.Sprintf(` %s(%s.%s) `, v.AggFunc, dorm.Quote(dbType, tableName), col))
	} else {
		builder.WriteString(fmt.Sprintf(` %s(%s) `, v.AggFunc, col))
	}

	builder.WriteString(fmt.Sprintf(` %s `, v.Operator))
	builder.WriteString(val)
	return builder.String()
}
