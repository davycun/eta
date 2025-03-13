package dorm

import (
	"fmt"
	"gorm.io/gorm"
	"strings"
)

var (
	DaMeng     DbType = "dm"
	PostgreSQL DbType = "postgres"
	Mysql      DbType = "mysql"
	Nebula     DbType = "nebula"
	Doris      DbType = "doris"
	ES         DbType = "es"
)

type DbType string

func (d DbType) String() string {
	return string(d)
}

func GetDbType(db *gorm.DB) DbType {
	tpStr := ""
	//理论上可能出现不支持的情况，或者一些驱动改动名字之后导致这个变动
	if strategy, ok := db.NamingStrategy.(NamingStrategy); ok {
		tpStr = strategy.Config.Type
	} else if s, ok1 := db.NamingStrategy.(*NamingStrategy); ok1 {
		tpStr = s.Config.Type
	}
	if tpStr == "" {
		tpStr = db.Name()
	}

	return DbType(tpStr)
}
func GetDbSchema(db *gorm.DB) string {
	scm := ""
	if strategy, ok := db.NamingStrategy.(NamingStrategy); ok {
		scm = strategy.Config.Schema
	} else if s, ok1 := db.NamingStrategy.(*NamingStrategy); ok1 {
		scm = s.Config.Schema
	}
	return scm
}
func GetDbTable(db *gorm.DB, tableName string) string {
	strategy, ok := db.NamingStrategy.(NamingStrategy)
	if ok {
		sch := strategy.Config.Schema
		dbType := GetDbType(db)
		switch dbType {
		case Mysql:
			return fmt.Sprintf("`%s`.`%s`", sch, tableName)
		default:
			return fmt.Sprintf(`"%s"."%s"`, sch, tableName)
		}
	}
	return tableName
}
func GetDbColumn(db *gorm.DB, tableName, column string) string {
	strategy, ok := db.NamingStrategy.(NamingStrategy)
	if ok {
		sch := strategy.Config.Schema
		dbType := GetDbType(db)
		switch dbType {
		case Mysql:
			if tableName == "" {
				return fmt.Sprintf("`%s`", column)
			}
			return fmt.Sprintf("`%s`.`%s`.`%s`", sch, tableName, column)
		default:
			if tableName == "" {
				return fmt.Sprintf(`"%s"`, column)
			}
			return fmt.Sprintf(`"%s"."%s"."%s"`, sch, tableName, column)
		}
	}
	return tableName
}
func GetDbUser(db *gorm.DB) string {
	strategy, ok := db.NamingStrategy.(NamingStrategy)
	if ok {
		return strings.ToUpper(strategy.Config.User)
	}
	return ""
}
func GetDbHost(db *gorm.DB) string {
	strategy, ok := db.NamingStrategy.(NamingStrategy)
	if ok {
		return strings.ToUpper(strategy.Config.Host)
	}
	return ""
}
func GetDbPort(db *gorm.DB) int {
	strategy, ok := db.NamingStrategy.(NamingStrategy)
	if ok {
		return strategy.Config.Port
	}
	return 0
}
func GetDbConfig(db *gorm.DB) Database {
	strategy, ok := db.NamingStrategy.(NamingStrategy)
	if ok {
		return strategy.Config
	}
	return Database{}
}

func EqualsDatabase(src *gorm.DB, target *gorm.DB) bool {
	srcHost := GetDbHost(src)
	srcPort := GetDbPort(src)
	srcSchema := GetDbSchema(src)
	targetHost := GetDbHost(target)
	targetPort := GetDbPort(target)
	targetSchema := GetDbSchema(target)

	return srcHost == targetHost && srcPort == targetPort && srcSchema == targetSchema
}
