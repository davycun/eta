package dorm

import (
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"strconv"
	"strings"
	"sync"
)

const (
	extraAppId = "appId"
)

type NamingStrategy struct {
	config Database
	schema.NamingStrategy
	extra sync.Map
}

func NewNamingStrategy(db Database) *NamingStrategy {
	return &NamingStrategy{config: db}
}

func (ns NamingStrategy) TableName(str string) string {
	return ns.config.Schema + "." + str
}
func (ns *NamingStrategy) SetAppId(appId string) {
	ns.extra.Store(extraAppId, appId)
}
func (ns *NamingStrategy) GetAppId() string {
	if x, ok := ns.extra.Load(extraAppId); ok {
		return x.(string)
	}
	return ""
}
func (ns *NamingStrategy) GetDatabase() Database {
	return ns.config
}

type Database struct {
	JsonType
	Host           string `json:"host" binding:"required" yaml:"host"`
	Port           int    `json:"port" binding:"required,lte=65535" yaml:"port"`
	User           string `json:"user" binding:"required" yaml:"user"`
	Password       string `json:"password" binding:"required" yaml:"password"`
	DBName         string `json:"db_name" binding:"required" yaml:"dbname"`
	Schema         string `json:"schema" yaml:"schema"`
	Type           string `json:"type" binding:"required,oneof=mysql postgres dm doris" yaml:"type"` //数据库的类型：mysql、postgres、dm
	MaxOpenCons    int    `json:"max_open_cons" yaml:"max_open_cons"`
	MaxIdleCons    int    `json:"max_idle_cons" yaml:"max_idle_cons"`
	ConMaxLifetime int    `json:"conn_max_lifetime" yaml:"conn_max_lifetime"`   //单位是秒
	ConMaxIdleTime int    `json:"conn_max_idle_time" yaml:"conn_max_idle_time"` //单位是秒
	Key            string `json:"key" yaml:"key"`
	LogLevel       int    `json:"log_level" yaml:"log_level"`
	SlowThreshold  int    `json:"slow_threshold" yaml:"slow_threshold"`
	SchemaPrefix   string `json:"schema_prefix" yaml:"schema_prefix"` //新建app的时候创建数据库schema的前缀
}

func (d *Database) GetKey() string {
	if d.Key != "" {
		return d.Key
	}
	d.Key = strings.Join([]string{d.Host, strconv.Itoa(d.Port), d.DBName, d.Schema, d.Type}, ":")
	return d.Key
}
func (d *Database) IsEmpty() bool {
	if d.Host == "" || d.Port == 0 || d.User == "" {
		return true
	}
	return false
}

type JsonType struct {
}

func (d JsonType) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return JsonGormDBDataType(db, field)
}

func (d JsonType) GormDataType() string {
	return JsonGormDataType()
}

func JsonGormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch GetDbType(db) {
	case PostgreSQL:
		return "jsonb"
	case DaMeng:
		return "CLOB"
	case Mysql, Doris:
		return "json"
	}
	return "jsonb"
}
func JsonGormDataType() string {
	return "json"
}
