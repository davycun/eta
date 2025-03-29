package ctype

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm"
	"gorm.io/gorm"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	TypeIdName          = "id"
	TypeArrayIntName    = "array_int"
	TypeArrayStringName = "array_string"
	TypeBoolName        = "bool"
	TypeNumericName     = "numeric"
	TypeGeometryName    = "geometry"
	TypeIntegerName     = "integer"
	TypeBigIntegerName  = "bigint"
	TypeJsonName        = "json"
	TypeStringName      = "string"
	TypeTextName        = "text"
	TypeTimeName        = "time"
	TypeTimestampName   = "timestamp"
	TypeTimestampTzName = "timestamp_tz"
	TypeFileName        = "file"
)

// 系统支持的数据数据类型，也是前端传入的数据类型
var (
	nullValue = []byte("null")
	dbTypeMap = map[string]dbTypeItem{
		TypeIdName:          {sys: TypeIdName, pg: "varchar", dm: "varchar", mysql: "VARCHAR(255)"},
		TypeArrayIntName:    {sys: TypeArrayIntName, pg: "integer[]", dm: "ARR_INT_CLS", mysql: "JSON"},
		TypeArrayStringName: {sys: TypeArrayStringName, pg: "varchar[]", dm: "ARR_STR_CLS", mysql: "JSON"},
		TypeBoolName:        {sys: TypeBoolName, pg: "boolean", dm: "BIT", mysql: "BOOL"},
		TypeNumericName:     {sys: TypeNumericName, pg: "numeric", dm: "NUMERIC", mysql: "NUMERIC"},
		TypeGeometryName:    {sys: TypeGeometryName, pg: "geometry", dm: "SYSGEO.ST_Geometry", mysql: "GEOMETRY"},
		TypeIntegerName:     {sys: TypeIntegerName, pg: "integer", dm: "INTEGER", mysql: "INT"},
		TypeBigIntegerName:  {sys: TypeBigIntegerName, pg: "bigint", dm: "BIGINT", mysql: "BIGINT"},
		TypeJsonName:        {sys: TypeJsonName, pg: "jsonb", dm: "CLOB", mysql: "JSON"},
		TypeStringName:      {sys: TypeStringName, pg: "varchar", dm: "VARCHAR", mysql: "TEXT"},
		TypeTextName:        {sys: TypeTextName, pg: "text", dm: "TEXT", mysql: "TEXT"},
		TypeTimeName:        {sys: TypeTimeName, pg: "timestamp with time zone", dm: "TIMESTAMP WITH TIME ZONE", mysql: "TIMESTAMP"},
		TypeTimestampName:   {sys: TypeTimestampName, pg: "timestamp with time zone", dm: "TIMESTAMP WITH TIME ZONE", mysql: "TIMESTAMP"},
		TypeTimestampTzName: {sys: TypeTimestampTzName, pg: "timestamp with time zone", dm: "TIMESTAMP WITH TIME ZONE", mysql: "TIMESTAMP"},
		TypeFileName:        {sys: TypeFileName, pg: "varchar[]", dm: "ARR_STR_CLS", mysql: "JSON"},
	}
	fieldType = map[string]reflect.Type{
		TypeIdName:          reflect.TypeFor[int64](),
		TypeArrayIntName:    reflect.TypeFor[Int64Array](),
		TypeArrayStringName: reflect.TypeFor[StringArray](),
		TypeBoolName:        reflect.TypeFor[Boolean](),
		TypeNumericName:     reflect.TypeFor[Float](),
		TypeGeometryName:    reflect.TypeFor[Geometry](),
		TypeIntegerName:     reflect.TypeFor[Integer](),
		TypeBigIntegerName:  reflect.TypeFor[Integer](),
		TypeJsonName:        reflect.TypeFor[Json](),
		TypeStringName:      reflect.TypeFor[String](),
		TypeTextName:        reflect.TypeFor[Text](),
		TypeTimeName:        reflect.TypeFor[LocalTime](),
		TypeTimestampName:   reflect.TypeFor[LocalTime](),
		TypeTimestampTzName: reflect.TypeFor[LocalTime](),
		TypeFileName:        reflect.TypeFor[StringArray](),
	}
)

type dbTypeItem struct {
	sys   string
	pg    string
	dm    string
	mysql string
}

// NewTypeValue
// isPrt 是否需要返回指针。如果是否，那么返回对应类型的零值
func NewTypeValue[T sql.ColumnType | *sql.ColumnType | *string | string](name T, isPrt bool) interface{} {

	var (
		dbType       = ""
		columnType   = sql.ColumnType{}
		isColumnType = false
		tp           interface{}
		precision    = 10
		scale        = 2
	)
	tp = name
	switch tp.(type) {
	case *sql.ColumnType:
		columnType = *tp.(*sql.ColumnType)
		isColumnType = true
	case sql.ColumnType:
		columnType = tp.(sql.ColumnType)
		isColumnType = true
	case *string:
		dbType = strings.ToLower(*tp.(*string))
		isColumnType = false
	case string:
		dbType = strings.ToLower(tp.(string))
		isColumnType = false
	}

	if isColumnType {
		dbType = strings.ToLower(columnType.DatabaseTypeName())
		size, scl, ok := columnType.DecimalSize()
		if ok {
			precision = int(size)
			scale = int(scl)
		}
	}
	// character varying(n) 、character(n)、char(n)、numeric(5,2)...
	if i := strings.Index(dbType, "("); i >= 0 {
		s := strings.Index(dbType, ",")
		if s > 0 {
			e := strings.Index(dbType, ")")
			p := dbType[s+1 : e]
			precision, _ = strconv.Atoi(p)
		}
		dbType = dbType[:i]
	}

	rsTp := newTypeValue(dbType, precision, scale, isPrt)
	if rsTp != nil {
		return rsTp
	}
	if isColumnType {
		if isPrt {
			return reflect.New(columnType.ScanType()).Interface()
		} else {
			return reflect.New(columnType.ScanType()).Elem().Interface()
		}
	}

	//最后什么都不行就用string接收吧？？
	if isPrt {
		return new(String)
	} else {
		return String{}
	}
}

func GetFieldType(tpStr string) (reflect.Type, bool) {
	tp, ok := fieldType[tpStr]
	return tp, ok
}

func GetDbTypeName(db *gorm.DB, tp string) (string, error) {
	var (
		dtItem, ok = dbTypeMap[tp]
		dbTp       = dorm.GetDbType(db)
		dbUser     = dorm.GetDbUser(db)
	)
	if !ok {
		//numeric(5,3)、varchar(256)
		pattern := `(.+)(\(([0-9]*)(,?)([0-9]*)\))`
		if matched, _ := regexp.MatchString(pattern, tp); matched {
			return tp, nil
		}
		return "", notSupportType(tp)
	}

	switch dbTp {
	case dorm.PostgreSQL:
		return dtItem.pg, nil
	case dorm.DaMeng:
		targetTp := strings.ToUpper(dtItem.dm)
		if targetTp == "ARR_INT_CLS" || targetTp == "ARR_STR_CLS" {
			targetTp = dbUser + "." + targetTp
		}
		return targetTp, nil
	case dorm.Mysql:
		return dtItem.mysql, nil
	default:
		return "", errors.New(fmt.Sprintf("not support type %s for database %s", tp, dbTp))
	}
}

func GetSupportType() []string {
	return []string{TypeIdName,
		TypeArrayIntName,
		TypeArrayStringName,
		TypeBoolName,
		TypeNumericName,
		TypeGeometryName,
		TypeIntegerName,
		TypeBigIntegerName,
		TypeJsonName,
		TypeStringName,
		TypeTextName,
		TypeTimeName,
		TypeTimestampName,
		TypeTimestampTzName,
		TypeFileName,
	}
}

// 根据字符串名字创建对应类型的实例
func newTypeValue(columnType string, precision int, scale int, isPrt bool) interface{} {

	switch columnType {
	case "character varying", "character", "char", "varchar", TypeStringName:
		if isPrt {
			return new(String)
		} else {
			return String{}
		}
	case "text":
		if isPrt {
			return new(Text)
		} else {
			return Text{}
		}
	case "boolean", "bit", TypeBoolName:
		if isPrt {
			return new(Boolean)
		} else {
			return Boolean{}
		}
	case TypeIntegerName, "int4", "int8", "int16", "int32", "int64", "smallint", "smallserial", "serial", "int", "bigint", "bigserial":
		if isPrt {
			return new(Integer)
		} else {
			return Integer{}
		}
	case TypeNumericName, "decimal", "double":
		if isPrt {
			f := new(Float)
			f.Precision = precision
			f.Scale = scale
			return f
		} else {
			f := Float{}
			f.Precision = precision
			f.Scale = scale
			return f
		}
	case TypeTimeName, TypeTimestampName, TypeTimestampTzName, "datetime", "date", "timestamp with time zone", "timestamp without time zone", "datetime with time zone", "datetime without time zone":
		if isPrt {
			return new(LocalTime)
		} else {
			return LocalTime{}
		}
	case "integer[]", "_int4", "_int", TypeArrayIntName:
		if isPrt {
			return new(Int64Array)
		} else {
			return Int64Array{}
		}
	case "text[]", "_text", "varchar[]", "character varying[]", "arr_str_cls", "_varchar", TypeArrayStringName, TypeFileName:
		if isPrt {
			return new(StringArray)
		} else {
			return StringArray{}
		}
	case "jsonb", TypeJsonName:
		if isPrt {
			return new(Json)
		} else {
			return Json{}
		}
	case "st_geometry", TypeGeometryName:
		if isPrt {
			return new(Geometry)
		} else {
			return Geometry{}
		}
	case "blob":
		///TODO error
		return make([]byte, 1024)
	default:
		return nil
	}
}

func notSupportType(tp string) error {
	return errors.New(fmt.Sprintf(`"not support type for %s"`, tp))
}
