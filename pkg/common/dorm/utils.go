package dorm

import (
	"fmt"
	"gorm.io/gorm"
	"strings"
)

func CopyGormSetting(src *gorm.DB, target *gorm.DB) {
	src.Statement.Settings.Range(func(key, value any) bool {
		target.Statement.Settings.Store(key, value)
		return true
	})
}
func JoinColumns(dbType DbType, tbName string, columns []string) string {
	if columns == nil || len(columns) < 1 {
		return ""
	}
	cols := make([]string, 0, len(columns))
	for _, v := range columns {
		qtCol := Quote(dbType, tbName, v)
		if qtCol != "" {
			cols = append(cols, qtCol)
		}
	}
	return strings.Join(cols, ",")
}
func JoinColumnsWithPrefixAlias(dbType DbType, tbName, prefix string, columns []string) string {

	if columns == nil || len(columns) < 1 {
		return ""
	}
	if strings.Contains(tbName, ".") {
		sp := strings.Split(tbName, ".")
		if len(sp) == 2 {
			tbName = sp[1]
		}
	}

	sqs := make([]string, 0, len(columns))

	for _, v := range columns {
		if v == "*" {
			continue
		}
		sqs = append(sqs, fmt.Sprintf("%s as %s", Quote(dbType, tbName, v), Quote(dbType, fmt.Sprintf("%s%s", prefix, v))))
	}
	if len(sqs) < 1 {
		return ""
	}
	return strings.Join(sqs, ",")

}

func CurrentDateTimeFunc(dbType DbType) string {
	switch dbType {
	case PostgreSQL, Mysql:
		return "now()"
	case DaMeng:
		return "CURRENT_TIMESTAMP()"
	}
	return ""
}

func CurrentIDFunc(dbType DbType, scm string) string {
	switch dbType {
	case PostgreSQL, Mysql:
		return fmt.Sprintf(`nextval('%s.%s')`, scm, SequenceIdName)
	case DaMeng:
		return fmt.Sprintf(`TO_CHAR("` + scm + `"."` + SequenceIdName + `".NEXTVAL)`)
	}
	return ""
}
func SplitSchemaTableName(scmTbName string) (scm string, tbName string) {
	if scmTbName == "" {
		return "", ""
	}
	if !strings.Contains(scmTbName, ".") {
		return "", scmTbName
	}
	st := strings.Split(scmTbName, ".")
	switch len(st) {
	case 0:
		return "", ""
	case 1:
		return "", st[0]
	case 2:
		return st[0], st[1]
	default:
		return st[0], st[len(st)-1]
	}
}
