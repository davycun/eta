package dorm

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/errs"
	"strings"

	"gorm.io/gorm"
)

// TriggerEnable 打开触发器
func TriggerEnable(db *gorm.DB, triggerName string, tableName string) error {
	dbType := GetDbType(db)
	q := quote(dbType, []string{triggerName, tableName})
	triggerName, tableName = q[0], q[1]

	switch dbType {
	case DaMeng:
		return db.Exec(fmt.Sprintf(`ALTER TRIGGER %s ENABLE`, triggerName)).Error
	case PostgreSQL:
		if strings.Contains(triggerName, ".") {
			triggerName = strings.Split(triggerName, ".")[1]
		}
		return db.Exec(fmt.Sprintf(`ALTER TABLE %s ENABLE TRIGGER %s`, tableName, triggerName)).Error
	case Mysql:
		return db.Exec(fmt.Sprintf(`ALTER TABLE %s ENABLE TRIGGER %s`, tableName, triggerName)).Error
	default:
		return errs.NewClientError("不支持的数据库类型")
	}
}

// TriggerDisable 关闭触发器
func TriggerDisable(db *gorm.DB, triggerName string, tableName string) error {
	dbType := GetDbType(db)
	q := quote(dbType, []string{triggerName, tableName})
	triggerName, tableName = q[0], q[1]

	switch dbType {
	case DaMeng:
		return db.Exec(fmt.Sprintf(`ALTER TRIGGER %s DISABLE`, triggerName)).Error
	case PostgreSQL:
		if strings.Contains(triggerName, ".") {
			triggerName = strings.Split(triggerName, ".")[1]
		}
		return db.Exec(fmt.Sprintf(`ALTER TABLE %s DISABLE TRIGGER %s`, tableName, triggerName)).Error
	case Mysql:
		return db.Exec(fmt.Sprintf(`ALTER TABLE %s DISABLE TRIGGER %s`, tableName, triggerName)).Error
	default:
		return errs.NewClientError("不支持的数据库类型")
	}
}

// TriggerDelete 删除触发器
func TriggerDelete(db *gorm.DB, triggerName string, tableName string) error {
	dbType := GetDbType(db)
	q := quote(dbType, []string{triggerName, tableName})
	triggerName, tableName = q[0], q[1]

	switch dbType {
	case DaMeng:
		return db.Exec(fmt.Sprintf(`DROP TRIGGER IF EXISTS %s`, triggerName)).Error
	case PostgreSQL:
		if strings.Contains(triggerName, ".") {
			triggerName = strings.Split(triggerName, ".")[1]
		}
		return db.Exec(fmt.Sprintf(`DROP TRIGGER IF EXISTS %s ON %s CASCADE`, triggerName, tableName)).Error
	case Mysql:
		return db.Exec(fmt.Sprintf(`DROP TRIGGER IF EXISTS %s`, triggerName)).Error
	default:
		return errs.NewClientError("不支持的数据库类型")
	}
}

// TriggerExists 触发器是否存在
func TriggerExists(db *gorm.DB, triggerName string, tableName string) (bool, error) {
	dbType := GetDbType(db)
	//q := quote(dbType, []string{triggerName, tableName})
	//triggerName, tableName = q[0], q[1]

	switch dbType {
	case DaMeng:
		var count int
		sq := `SELECT COUNT(*) FROM SYS.ALL_TRIGGERS WHERE owner=? AND table_name=? AND trigger_name=?`
		err := db.Raw(sq, GetDbSchema(db), tableName, triggerName).Scan(&count).Error
		return count > 0, err
	case PostgreSQL:
		if strings.Contains(triggerName, ".") {
			triggerName = strings.Split(triggerName, ".")[1]
		}
		var count int
		sq := `SELECT COUNT(pg_trigger.*) FROM pg_trigger, pg_class 
             	WHERE tgrelid=pg_class.oid AND pg_class.relname=? AND pg_trigger.tgname=?`
		err := db.Raw(sq, tableName, triggerName).Scan(&count).Error
		return count > 0, err
	case Mysql:
		var count int
		sq := "SELECT COUNT(*) FROM information_schema.triggers WHERE trigger_name=?"
		err := db.Raw(sq, triggerName).Scan(&count).Error
		return count > 0, err
	default:
		return false, errs.NewClientError("不支持的数据库类型")
	}
}

func quote(dbType DbType, str []string) []string {
	strs := make([]string, 0, len(str))
	for _, v := range str {
		if strings.HasPrefix(v, ".") {
			v = strings.TrimPrefix(v, ".")
		}
		if strings.Contains(v, ".") {
			strs = append(strs, Quote(dbType, strings.Split(v, ".")...))
		} else {
			strs = append(strs, Quote(dbType, v))
		}
	}
	return strs
}
