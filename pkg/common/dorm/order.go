package dorm

import (
	"fmt"
	"strconv"
	"strings"
)

type OrderBy struct {
	Column string `json:"column"`
	Asc    bool   `json:"asc"`
}

func ResolveOrderByString(orderBy []OrderBy, tableName, defaultOrderBy string, addKeyWords bool) string {
	bd := strings.Builder{}

	for i, _ := range orderBy {
		if i == 0 {
			if addKeyWords {
				bd.WriteString(" order by ")
			}
		} else {
			bd.WriteString(",")
		}

		if tableName != "" {
			bd.WriteString(fmt.Sprintf(`"%s"."%s"`, tableName, orderBy[i].Column))
		} else {
			bd.WriteString(fmt.Sprintf(`"%s"`, orderBy[i].Column))
		}

		if orderBy[i].Asc {
			bd.WriteString(" asc ")
		} else {
			bd.WriteString(" desc ")
		}
	}
	od := bd.String()
	if od == "" && defaultOrderBy != "" {
		return " order by " + defaultOrderBy
	}
	return od
}
func ResolveLimitString(offset, limit int) string {
	if limit < 1 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	return ` limit ` + strconv.Itoa(limit) + ` offset ` + strconv.Itoa(offset)
}
func ResolveOrderDesc(asc bool) string {
	if asc {
		return "asc"
	}
	return "desc"
}

func ResolveEsOrderBy(orderByList ...OrderBy) []map[string]interface{} {

	esOrderBy := make([]map[string]interface{}, 0, len(orderByList))
	if len(orderByList) < 1 {
		return esOrderBy
	}
	for _, orderBy := range orderByList {
		col := orderBy.Column + ".keyword"
		esOrderBy = append(esOrderBy, map[string]interface{}{
			col: map[string]interface{}{
				"order": ResolveOrderDesc(orderBy.Asc),
			},
		})
	}
	return esOrderBy
}
