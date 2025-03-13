package dorm

import (
	"fmt"
	"strings"
)

// PartitionColumn
// 这个组装出来的例子类似  count(aggregate_column) over (partition by column1,column2) as alias
// 类似count 这类的聚合函数可以使用 *，但是如果是max或者min这的聚合函数AggregateColumn必须要指定列名
type PartitionColumn struct {
	Columns   []string `json:"columns" binding:"required"`
	AggColumn string   `json:"agg_column" binding:"required"`
	AggFunc   string   `json:"agg_func" binding:"required;oneof=max min avg count"`
	Alias     string   `json:"alias" binding:"required"`
}

func ResolvePartitionColumn(dbType DbType, tableName string, pc []PartitionColumn) string {
	bd := strings.Builder{}
	for i, v := range pc {
		if i > 0 {
			bd.WriteByte(',')
		}

		alias := v.Alias
		if alias == "" {
			alias = "_" + v.AggFunc
		}

		if strings.ToLower(v.AggFunc) == "count" {
			bd.WriteString(` count(*) `)
		} else {
			if tableName != "" {
				bd.WriteString(fmt.Sprintf(`%s("%s"."%s") `, v.AggFunc, tableName, v.AggColumn))
			} else {
				bd.WriteString(fmt.Sprintf(`%s("%s") `, v.AggFunc, v.AggColumn))
			}
		}

		bd.WriteString(fmt.Sprintf(` over (partition by %s) as %s`, JoinColumns(dbType, tableName, v.Columns), Quote(dbType, alias)))
	}
	return bd.String()

}
