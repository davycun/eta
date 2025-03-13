package ctype

import (
	"database/sql"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/modern-go/reflect2"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"maps"
	"reflect"
	"strings"
	"sync"
)

var (
	//colTypeCache = make(map[reflect.Type]map[string]reflect.Type)
	//存储的结构是 reflect.Type -> map[string]reflect.Type
	colTypeCache = sync.Map{}
	scanTp       = reflect.TypeFor[sql.Scanner]()
)

func ScanRows(db *gorm.DB, dest map[string]reflect.Type) (result []Map, err error) {
	var (
		columnTypes []*sql.ColumnType
		rows        *sql.Rows
	)

	rows, err = db.Rows()
	if err != nil {
		return nil, err
	}
	if rows != nil {
		defer rows.Close()
	}

	columnTypes, err = rows.ColumnTypes()
	columns, err := rows.Columns()
	if err != nil {
		return
	}
	rs := make([]Map, 0, 100)
	for rows.Next() {
		row := make([]interface{}, len(columnTypes))
		for i, v := range columnTypes {
			var (
				tp reflect.Type
				ok bool
			)
			if dest != nil {
				tp, ok = dest[columns[i]]
			}
			if ok {
				row[i] = reflect.New(tp).Interface()
			} else {
				row[i] = NewType(v, true)
			}
		}
		err = rows.Scan(row...)
		h := Map{}
		for i, v := range columnTypes {
			h[v.Name()] = row[i]
		}
		rs = append(rs, h)
	}
	return rs, err
}

func GetType(name string) (reflect.Type, bool) {
	tp, ok := FieldType[name]
	return tp, ok
}

func GetColType(obj any) map[string]reflect.Type {

	var (
		tp = reflect2.TypeOf(obj).Type1()
	)
	m, ok := colTypeCache.Load(tp)
	if !ok {
		if tp.Kind() == reflect.Pointer {
			m, ok = colTypeCache.Load(tp.Elem())
		}
	}
	if ok {
		//遍历m,返回一个新map，避免使用者直接操作map会导致并发写的问题
		m1 := make(map[string]reflect.Type)
		m2 := m.(map[string]reflect.Type)
		maps.Copy(m1, m2)
		return m1
	}
	m1 := structFieldType(tp)
	colTypeCache.Store(tp, m1)
	return m1
}

func structFieldType(tp reflect.Type) map[string]reflect.Type {

	var (
		rs = make(map[string]reflect.Type)
	)

	switch tp.Kind() {
	case reflect.Pointer:
		return structFieldType(tp.Elem())
	case reflect.Struct:
		//默认往后执行
	default:
		return rs
	}

	for i := 0; i < tp.NumField(); i++ {
		fd := tp.Field(i)
		if !fd.IsExported() {
			continue
		}

		var (
			colName        = ""
			jsName         = getJsonName(fd.Tag) //json tag定义的名字
			gormTag        = getGormTag(fd.Tag)  //gorm的tag
			gormColName    = gormTag["COLUMN"]   //gorm定义的column名字
			serializerType = gormTag["SERIALIZER"]
			embedded       = gormTag["EMBEDDED"]
			gormIgnore     = gormTag["-"]
			embeddedPrefix = gormTag[strings.ToUpper("embeddedPrefix")]
			fdType         = fd.Type
			fdPrtType      = fd.Type
		)
		if gormIgnore != "" {
			continue
		}
		//处理字段名称
		if gormColName != "" {
			colName = gormColName
		} else if jsName != "" {
			colName = jsName
		} else {
			colName = utils.HumpToUnderline(fd.Name)
		}

		switch fdType.Kind() {
		case reflect.Pointer:
			fdType = fd.Type.Elem()
		case reflect.Struct:
			fdPrtType = reflect.New(fdType).Type()
		case reflect.Bool:
			rs[colName] = FieldType[TpBool]
			continue
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			rs[colName] = FieldType[TpInteger]
			continue
		case reflect.Float32, reflect.Float64:
			rs[colName] = FieldType[TpNumeric]
			continue
		case reflect.String:
			rs[colName] = FieldType[TpString]
			continue

		default:
			continue
		}

		//处理字段类型
		//if embedded != "" && fdType.Kind() == reflect.Struct {
		//如果实现sql.Scanner接口表示这个结构体就是一个包装类型，不需要再对当前结构体内部进行解析
		if fdType.Kind() == reflect.Struct && !fdPrtType.Implements(scanTp) {
			tmp := GetColType(reflect.New(fd.Type).Interface())
			//处理内嵌字段
			//tmp := structFieldType(fdType)
			if embedded != "" {
				for k, v := range tmp {
					if embeddedPrefix != "" {
						k = embeddedPrefix + k
					}
					rs[k] = v
				}
			} else {
				for k, v := range tmp {
					if _, ok := rs[k]; !ok {
						rs[k] = v
					}
				}
			}
			continue
		}

		if serializerType == "json" || serializerType == "jsonb" {
			//处理json序列化
			fdType = FieldType[TpJson]

		}
		rs[colName] = fdType
	}
	return rs
}

func getGormTag(tag reflect.StructTag) map[string]string {
	value := tag.Get("gorm")
	return schema.ParseTagSetting(value, ";")
}

func getJsonName(tag reflect.StructTag) string {
	var (
		j = tag.Get("json")
	)
	j = strings.TrimSpace(j)

	if j != "" {
		jn := strings.Split(j, ",")
		return jn[0]
	}
	return ""
}
