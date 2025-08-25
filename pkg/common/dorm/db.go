package dorm

import (
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
	tpStr := GetDbConfig(db).Type
	if tpStr == "" {
		tpStr = db.Name()
	}
	return DbType(tpStr)
}
func GetDbSchema(db *gorm.DB) string {
	return GetDbConfig(db).Schema
}
func GetScmTableName(db *gorm.DB, tableName string) string {
	return Quote(GetDbType(db), GetDbSchema(db), tableName)
}
func GetDbUser(db *gorm.DB) string {
	return GetDbConfig(db).User
}
func GetDbHost(db *gorm.DB) string {
	return strings.ToUpper(GetDbConfig(db).Host)
}
func GetDbPort(db *gorm.DB) int {
	return GetDbConfig(db).Port
}
func GetDbConfig(db *gorm.DB) Database {
	return getNamingStrategy(db).GetDatabase()
}
func GetAppId(db *gorm.DB) string {
	return getNamingStrategy(db).GetAppId()
}
func GetAppIdOrSchema(db *gorm.DB) string {
	id := getNamingStrategy(db).GetAppId()
	if id == "" {
		id = GetDbSchema(db)
	}
	return id
}
func SetAppId(db *gorm.DB, appId string) {
	getNamingStrategy(db).SetAppId(appId)
}

func getNamingStrategy(db *gorm.DB) *NamingStrategy {
	if db == nil {
		return NewNamingStrategy(Database{})
	}
	if x, ok := db.NamingStrategy.(NamingStrategy); ok {
		return &x
	}
	if x, ok := db.NamingStrategy.(*NamingStrategy); ok {
		return x
	}
	return NewNamingStrategy(Database{})
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
