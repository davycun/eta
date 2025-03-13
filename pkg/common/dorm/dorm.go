package dorm

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/logger"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

const (
	GormInTransactionKey = "GormInTransactionKey"
	SequenceIdName       = "seq_eta_id"
)

func InTransaction(db *gorm.DB) bool {
	if db == nil {
		return false
	}
	_, b := db.Get(GormInTransactionKey)
	return b
}
func SetInTransaction(db *gorm.DB) {
	if db == nil {
		return
	}
	db1 := db.Set(GormInTransactionKey, true)
	CopyGormSetting(db1, db)
}

func CommitOrRollback(txDb *gorm.DB, err error) {
	var (
		err1 error
	)
	if !InTransaction(txDb) {
		return
	}
	if err != nil {
		err1 = txDb.Rollback().Error
	} else {
		err1 = txDb.Commit().Error
	}
	if err1 != nil {
		logger.Errorf("DB Rollback Or Commit failed. %v", err)
	}
}
func Transaction(db *gorm.DB) *gorm.DB {
	if db == nil {
		return db
	}
	if InTransaction(db) {
		return db
	}
	tx := db.Begin()
	SetInTransaction(tx)
	return tx
}

func Quote(dbType DbType, str ...string) string {

	if len(str) < 1 {
		return ""
	}

	cols := make([]string, 0, len(str))

	for _, v := range str {
		if v == "" {
			continue
		}
		if v == "*" {
			cols = append(cols, "*")
			continue
		}
		switch dbType {
		case PostgreSQL, DaMeng:
			cols = append(cols, fmt.Sprintf(`"%s"`, v))
		case Mysql, Doris:
			cols = append(cols, fmt.Sprintf("`%s`", v))
		default:
			cols = append(cols, fmt.Sprintf(`"%s"`, v))
		}
	}
	return strings.Join(cols, ".")
}
func QuotePlaceholder(dbType DbType, str string, args ...string) string {

	bd := strings.Builder{}
	i := 0
	for _, v := range []byte(str) {
		switch v {
		case '?':
			bd.WriteString(Quote(dbType, args[i]))
			i++
		default:
			bd.WriteByte(v)
		}
	}
	return bd.String()
}

func QuoteSchemaTableName(db *gorm.DB, tableName string) string {
	return Quote(GetDbType(db), GetDbSchema(db), tableName)
}
func BoolValue(dbType DbType, val bool) string {
	switch dbType {
	case DaMeng:
		if val {
			return "1"
		}
		return "0"
	default:
		return strconv.FormatBool(val)
	}
}

// Table
// 如果想自己控制schema，那么tableName就直接携带着schema即可，比如SYS.OBJECTS
func Table(db *gorm.DB, tableName string) *gorm.DB {
	var (
		//注意这里传入的schemaTableName不要用引号括起来，否则mysql会有问题
		scmTbName = fmt.Sprintf("%s.%s", GetDbSchema(db), tableName)
	)
	if strings.Contains(tableName, ".") {
		return db.Table(tableName)
	}
	return db.Table(scmTbName)
}
func TableWithContext(db *gorm.DB, c *ctx.Context, tableName string) *gorm.DB {
	return Table(WithContext(db, c), tableName)
}

func LoadAndDelete(db *gorm.DB, key string) (interface{}, bool) {
	return db.Statement.Settings.LoadAndDelete(key)
}
func Store(db *gorm.DB, key string, value any) {
	db.Statement.Settings.Store(key, value)
}
