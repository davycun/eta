package ctype

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/modern-go/reflect2"
	"gorm.io/gorm"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	NullValue = []byte("null")
	// 系统支持的数据数据类型，也是前端传入的数据类型
	TpId          = "id"
	TpArrayInt    = "array_int"
	TpArrayString = "array_string"
	TpBool        = "bool"
	TpNumeric     = "numeric"
	TpGeometry    = "geometry"
	TpInteger     = "integer"
	TpBigInteger  = "bigint"
	TpJson        = "json"
	TpString      = "string"
	TpText        = "text"
	TpTime        = "time"
	TpFile        = "file"
	DbTypeMap     = map[string]dbTypeItem{
		TpId:          {sys: TpId, pg: "varchar", dm: "varchar", mysql: "VARCHAR(255)"},
		TpArrayInt:    {sys: TpArrayInt, pg: "integer[]", dm: "ARR_INT_CLS", mysql: "JSON"},
		TpArrayString: {sys: TpArrayString, pg: "text[]", dm: "ARR_STR_CLS", mysql: "JSON"},
		TpBool:        {sys: TpBool, pg: "boolean", dm: "BIT", mysql: "BOOL"},
		TpNumeric:     {sys: TpNumeric, pg: "numeric", dm: "NUMERIC", mysql: "NUMERIC"},
		TpGeometry:    {sys: TpGeometry, pg: "geometry", dm: "SYSGEO.ST_Geometry", mysql: "GEOMETRY"},
		TpInteger:     {sys: TpInteger, pg: "integer", dm: "INTEGER", mysql: "INT"},
		TpBigInteger:  {sys: TpBigInteger, pg: "bigint", dm: "BIGINT", mysql: "BIGINT"},
		TpJson:        {sys: TpJson, pg: "jsonb", dm: "CLOB", mysql: "JSON"},
		TpString:      {sys: TpString, pg: "varchar", dm: "VARCHAR", mysql: "TEXT"},
		TpText:        {sys: TpText, pg: "text", dm: "TEXT", mysql: "TEXT"},
		TpTime:        {sys: TpTime, pg: "timestamp with time zone", dm: "TIMESTAMP WITH TIME ZONE", mysql: "TIMESTAMP"},
		TpFile:        {sys: TpFile, pg: "text[]", dm: "ARR_STR_CLS", mysql: "JSON"},
	}
	FieldType = map[string]reflect.Type{
		TpId:          reflect2.TypeOf((*int64)(nil)).Type1().Elem(),
		TpArrayInt:    reflect2.TypeOf((*Int64Array)(nil)).Type1().Elem(),
		TpArrayString: reflect2.TypeOf((*StringArray)(nil)).Type1().Elem(),
		TpBool:        reflect2.TypeOf((*Boolean)(nil)).Type1().Elem(),
		TpNumeric:     reflect2.TypeOf((*Float)(nil)).Type1().Elem(),
		TpGeometry:    reflect2.TypeOf((*Geometry)(nil)).Type1().Elem(),
		TpInteger:     reflect2.TypeOf((*Integer)(nil)).Type1().Elem(),
		TpBigInteger:  reflect2.TypeOf((*Integer)(nil)).Type1().Elem(),
		TpJson:        reflect2.TypeOf((*Json)(nil)).Type1().Elem(),
		TpString:      reflect2.TypeOf((*String)(nil)).Type1().Elem(),
		TpText:        reflect2.TypeOf((*Text)(nil)).Type1().Elem(),
		TpTime:        reflect2.TypeOf((*LocalTime)(nil)).Type1().Elem(),
		TpFile:        reflect2.TypeOf((*StringArray)(nil)).Type1().Elem(),
	}
)

type dbTypeItem struct {
	sys   string
	pg    string
	dm    string
	mysql string
}

// NewType isPrt 是否需要返回指针。如果是否，那么返回对应类型的零值
func NewType[T sql.ColumnType | *sql.ColumnType | *string | string](name T, isPrt bool) interface{} {

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

	rsTp := getType(dbType, precision, scale, isPrt)
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

func GetDbTypeName(db *gorm.DB, tp string) (string, error) {
	var (
		dtItem, ok = DbTypeMap[tp]
		dbTp       = dorm.GetDbType(db)
		dbUser     = dorm.GetDbUser(db)
	)
	if !ok {
		if strings.HasPrefix(strings.ToLower(tp), "numeric(") {
			matched, _ := regexp.MatchString(`numeric\([1-9]\d*,[1-9]\d*\)`, tp)
			if !matched {
				return "", notSupportType(tp)
			}
			return tp, nil
		}
		return "", notSupportType(tp)
	}

	switch dbTp {
	case dorm.PostgreSQL:
		return dtItem.pg, nil
	case dorm.DaMeng:
		targetTp := dtItem.dm
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

func getType(columnType string, precision int, scale int, isPrt bool) interface{} {

	switch columnType {
	case "character varying", "character", "char", "varchar", TpString:
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
	case "boolean", "bit", TpBool:
		if isPrt {
			return new(Boolean)
		} else {
			return Boolean{}
		}
	case TpInteger, "int4", "int8", "int16", "int32", "int64", "smallint", "smallserial", "serial", "int", "bigint", "bigserial":
		if isPrt {
			return new(Integer)
		} else {
			return Integer{}
		}
	case TpNumeric, "decimal", "double":
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
	case TpTime, "timestamp", "datetime", "date", "timestamp with time zone", "timestamp without time zone", "datetime with time zone", "datetime without time zone":
		if isPrt {
			return new(LocalTime)
		} else {
			return LocalTime{}
		}
	case "integer[]", "_int4", "_int", TpArrayInt:
		if isPrt {
			return new(Int64Array)
		} else {
			return Int64Array{}
		}
	case "text[]", "_text", "varchar[]", "character varying[]", "arr_str_cls", "_varchar", TpArrayString, TpFile:
		if isPrt {
			return new(StringArray)
		} else {
			return StringArray{}
		}
	case "jsonb", TpJson:
		if isPrt {
			return new(Json)
		} else {
			return Json{}
		}
	case "st_geometry", TpGeometry:
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
