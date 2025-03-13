package dorm

import (
	"fmt"
	"strings"
)

const (
	AggFuncCount       = "count"
	AggFuncMax         = "max"
	AggFuncMin         = "min"
	AggFuncAvg         = "avg"
	AggFuncCardinality = "cardinality"
	AggFuncValueCount  = "value_count"
)

// AggregateColumn
// 组装出来的内容类似如下，max(column) as alias ，其中max是aggregate_func
type AggregateColumn struct {
	Column   string `json:"column" binding:"required"`
	AggFunc  string `json:"agg_func" binding:"required,oneof=count max min avg cardinality"` //cardinality只为给ES用
	Distinct bool   `json:"distinct,omitempty"`
	Alias    string `json:"alias" binding:"required"`
}

func ResolveAggregateColumn(dbType DbType, tableName string, ac ...AggregateColumn) string {

	bd := strings.Builder{}
	for i, v := range ac {
		if i > 0 {
			bd.WriteByte(',')
		}

		alias := v.Alias
		if alias == "" {
			alias = "_" + v.AggFunc
		}

		if strings.ToLower(v.AggFunc) == "count" {
			bd.WriteString(`count(*) `)
		} else {
			distinct := ""
			if v.Distinct {
				distinct = "distinct"
			}
			bd.WriteString(fmt.Sprintf(`%s(%s %s) `, v.AggFunc, distinct, Quote(dbType, tableName, v.Column)))
		}
		bd.WriteString(fmt.Sprintf(` as %s`, Quote(dbType, alias)))
	}
	return bd.String()
}
