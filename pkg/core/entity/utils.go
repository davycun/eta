package entity

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/eta/constants"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"reflect"
	"strings"
)

func SliceToMap[T any](key string, data ...T) map[string]T {
	rs := make(map[string]T)
	for _, v := range data {
		if k := GetString(v, key); k != "" {
			rs[k] = v
		}
	}
	return rs
}

func GetDefaultColumns(e any) []string {
	if x, ok := e.(ColumnDefaultInterface); ok {
		return x.DefaultColumns()
	}
	return []string{"*"}
}
func GetMustColumns(e any) []string {
	if x, ok := e.(ColumnsMustInterface); ok {
		return x.MustColumns()
	}
	return []string{IdDbName, UpdatedAtDbName}
}
func GetWideTableName(e any) string {
	if x, ok := e.(WideInterface); ok {
		return x.WideTableName()
	}
	if x, ok := e.(schema.TablerWithNamer); ok {
		tbName := x.TableName(nil)
		if strings.Contains(tbName, "_") {
			tb := tbName[strings.Index(tbName, "_")+1:]
			return fmt.Sprintf("t_%s", tb)
		}
		return fmt.Sprintf("t_wide_%s", tbName)
	}
	return "t_wide"
}
func GetFullWideIndexName(db *gorm.DB, obj any) string {
	return fmt.Sprintf("%s_%s", dorm.GetDbSchema(db), GetWideTableName(obj))
}
func SetTableName(db *gorm.DB, obj any) *gorm.DB {
	//注意这里传入的schemaTableName不要用引号括起来，否则mysql会有问题
	return dorm.Table(db, GetTableName(obj))
}

// GetTableName
// 获取实体对应的表名
// 1. 如果实现了schema.TablerWithNamer或者schema.Tabler那么调用对应的TableName方法返回表名
// 2. 从实体e中查找有没有一个字段TableName，如果有并且不为空返回对应的表名
// 3.获取实体包名及实体名，下划线取代驼峰及dot
func GetTableName(e any) string {
	if e == nil {
		return ""
	}
	stb := ""
	if x, ok := e.(schema.TablerWithNamer); ok {
		stb = x.TableName(nil)
	} else if x1, ok1 := e.(schema.Tabler); ok1 {
		stb = x1.TableName()
	}
	_, tb := dorm.SplitSchemaTableName(stb)
	if tb != "" {
		return tb
	}
	tb = GetString(e, constants.TemplateTableNameField)
	if tb == "" {
		tb = utils.GetEntityName(reflect.TypeOf(e))
	}
	return tb
}

func GetRaDbFields(obj any) []string {
	fields := make([]string, 0)
	if x, ok := obj.(RaInterface); ok {
		fields = x.RaDbFields()
		return fields
	}

	fds, b := Get(obj, constants.TemplateRaDbFields)
	if b {
		fields = fds.([]string)
	}
	return fields
}

func SupportRA(obj any) bool {
	fs := GetRaDbFields(obj)
	return len(fs) > 0
}

func SupportEs(e any) bool {
	if _, ok := e.(EsInterface); ok {
		return true
	}
	return false
}

// GetEsIndexName 获取 ES 索引名
func GetEsIndexName(scm string, name string) string {
	return fmt.Sprintf("%s_%s", scm, name)
}

// GetEsIndexNameByDb 获取 ES 索引名
func GetEsIndexNameByDb(db *gorm.DB, e any) string {
	idxName := GetTableName(e)
	if db == nil {
		return idxName
	}
	return GetEsIndexName(dorm.GetDbSchema(db), idxName)
}

// GetRetrieveIndexName 获取 ES 检索索引名
func GetRetrieveIndexName(db *gorm.DB, e any) string {
	// 有宽表先用宽表
	if _, ok := e.(WideInterface); ok {
		return GetFullWideIndexName(db, e)
	}
	// 没有宽表就用RA表
	return GetEsIndexNameByDb(db, e)
}
