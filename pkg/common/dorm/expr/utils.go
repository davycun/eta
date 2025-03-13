package expr

import (
	"github.com/davycun/eta/pkg/common/dorm"
	"strings"
)

func JoinExprColumn(dbType dorm.DbType, tbName string, columns ...ExpColumn) (string, error) {
	bd := strings.Builder{}
	for i, v := range columns {
		column, err := ExplainExprColumn(dbType, v, tbName)
		if err != nil {
			return "", err
		}
		if i > 0 {
			bd.WriteByte(',')
		}

		bd.WriteString(column)
	}
	return bd.String(), nil
}

func NewAliasColumn(col, alias string) ExpColumn {
	ec := ExpColumn{}
	ec.Expr = col
	ec.Alias = alias
	return ec
}
