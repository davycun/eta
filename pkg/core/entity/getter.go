package entity

import (
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/tag"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/schema"
	"reflect"
)

var (
	gormFieldCache  = map[reflect.Type]map[string]string{} //reflect.Type -> []string
	tableFieldCache = map[reflect.Type][]TableField{}      //reflect.Type -> []string
)

type Getter interface {
	Get(field string) (any, bool)
}

func GetString(src any, key string) string {

	dt, b := Get(src, key)
	if b {
		return ctype.ToString(dt)
	}
	return ""
}

func Get(src interface{}, key string) (interface{}, bool) {

	if src == nil {
		return nil, false
	}

	jsonKey := utils.HumpToUnderline(key)
	if x, ok := src.(Getter); ok {
		return x.Get(jsonKey)
	}

	if x, ok := src.(reflect.Value); ok {
		val := GetValue(x, key)
		if val.IsValid() && val.CanInterface() {
			return val.Interface(), true
		}
		return nil, false
	}

	switch x := src.(type) {
	case *map[string]any:
		d, ok := (*x)[jsonKey]
		return d, ok
	case map[string]any:
		d, ok := x[jsonKey]
		return d, ok
	case gin.H:
		d, ok := x[jsonKey]
		return d, ok
	case *gin.H:
		d, ok := (*x)[jsonKey]
		return d, ok
	case ctype.Map:
		d, ok := x[jsonKey]
		return d, ok
	case *ctype.Map:
		d, ok := (*x)[jsonKey]
		return d, ok
	}

	val := GetValue(reflect.ValueOf(src), key)
	if val.IsValid() && val.CanInterface() {
		return val.Interface(), true
	}
	return nil, false
}

func GetValue(val reflect.Value, key string) reflect.Value {

	if !val.IsValid() {
		return reflect.Value{}
	}
	val = utils.GetRealValue(val)

	var (
		fieldKey   = utils.UnderlineToHump(key) //TODO 理论上要根据jsonTag或者gormTag来确定名字
		jsonKey    = utils.HumpToUnderline(key)
		fieldValue = val.FieldByName(fieldKey)
		valType    = val.Type()
	)
	if fieldValue.IsValid() {
		return fieldValue
	}

	switch valType.Kind() {
	case reflect.Pointer:
		return GetValue(val.Elem(), key)
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			var (
				fieldVal  = val.Field(i)
				field     = valType.Field(i)
				fieldType = field.Type
				jsonTag   = tag.NewJsonTag(field.Tag.Get("json")).GetName()
				gormTag   = tag.New(field.Tag.Get("gorm"))
			)

			if jsonKey == jsonTag || gormTag.Get("column") == jsonKey {
				return fieldVal
			}

			//包装的结构体，无需地推查找，直接返回
			if val.CanInterface() {
				if _, ok := val.Interface().(schema.GormDataTypeInterface); ok {
					continue
				}
			}
			//直接对应数据库的一个字段，无需当做结构体再次查找
			if gormTag.Get("serializer") != "" {
				continue
			}
			if fieldType.Kind() == reflect.Struct && (field.Anonymous || gormTag.Get("embedded") != "") {
				//只有组合字段才继续查找
				fv := GetValue(fieldVal, jsonKey)
				if fv.IsValid() {
					return fv
				}
			}
		}
	case reflect.Map:
		var (
			valInter = val.Interface()
		)
		if x, ok := valInter.(Getter); ok {
			v, ok1 := x.Get(jsonKey)
			if ok1 {
				return reflect.ValueOf(v)
			}
		}
		switch x := valInter.(type) {
		case ctype.Map:
			v, ok := x[jsonKey]
			if ok {
				return reflect.ValueOf(v)
			}

		case *ctype.Map:
			v, ok := (*x)[jsonKey]
			if ok {
				return reflect.ValueOf(v)
			}
		case gin.H:
			v, ok := x[jsonKey]
			if ok {
				return reflect.ValueOf(v)
			}
		case *gin.H:
			v, ok := (*x)[jsonKey]
			if ok {
				return reflect.ValueOf(v)
			}
		case map[string]any:
			v, ok := x[jsonKey]
			if ok {
				return reflect.ValueOf(v)
			}
		case *map[string]any:
			v, ok := (*x)[jsonKey]
			if ok {
				return reflect.ValueOf(v)
			}
		}
	default:

	}

	return reflect.Value{}
}

func GetGormFieldName(obj any) map[string]string {
	tp := reflect.TypeOf(obj)
	if tp.Kind() == reflect.Pointer {
		tp = tp.Elem()
	}
	if x, ok := gormFieldCache[tp]; ok {
		return x
	}

	var (
		gormField = make(map[string]string) //FieldName -> DbColumnName
	)

	fieldCount := tp.NumField()
	for i := 0; i < fieldCount; i++ {
		structField := tp.Field(i)
		if !structField.IsExported() {
			continue
		}

		fieldTyp := structField.Type
		if fieldTyp.Kind() == reflect.Pointer {
			fieldTyp = fieldTyp.Elem()
		}

		var (
			fieldVal  = reflect.New(fieldTyp)
			gormTg    = tag.New(structField.Tag.Get("gorm"))
			jsonTg    = tag.NewJsonTag(structField.Tag.Get("json"))
			fieldName = structField.Name
			colName   = gormTg.Get("column")
		)

		if gormTg.Get("-") != "" {
			continue
		}
		if colName == "" {
			colName = jsonTg.GetName()
		}
		if colName == "" {
			colName = utils.HumpToUnderline(structField.Name)
		}

		//自定义的包装类型，基本都实现了这个接口
		if _, ok := fieldVal.Interface().(schema.GormDataTypeInterface); ok {
			gormField[fieldName] = colName
			continue
		}
		if gormTg.Get("serializer") != "" {
			gormField[fieldName] = colName
			continue
		}

		switch fieldTyp.Kind() {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int,
			reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint, reflect.Float32, reflect.Float64, reflect.String, reflect.Bool:
			gormField[fieldName] = colName
		case reflect.Map, reflect.Slice, reflect.Array:
		case reflect.Struct:
			if gormTg.Get("column") != "" {
				gormField[fieldName] = colName
				continue
			}
			//TODO 这里返回使用可能会有问题，可能需要递归调用返回嵌入字段，同时修改返回值指定要更新的嵌入字段。有待测试
			pre := gormTg.Get("embeddedPrefix")
			if pre != "" {
				gormField[fieldName] = colName
				continue
			}
			ff := GetGormFieldName(fieldVal.Interface())
			for k, v := range ff {
				gormField[k] = v
			}
			continue
		default:

		}
	}
	gormFieldCache[tp] = gormField
	return gormField
}

func GetTableFields(obj any) []TableField {

	var (
		tp reflect.Type
	)
	if x, ok := tp.(reflect.Type); ok {
		tp = x
	} else if obj != nil {
		tp = reflect.TypeOf(obj)
	} else {
		return []TableField{}
	}

	if tp.Kind() == reflect.Pointer {
		tp = tp.Elem()
	}
	if x, ok := tableFieldCache[tp]; ok {
		return x
	}

	var (
		fieldCount     = tp.NumField()
		tableFieldList = make([]TableField, 0, fieldCount)
	)

	for i := 0; i < fieldCount; i++ {
		var (
			tbField     = TableField{}
			structField = tp.Field(i)
			fieldTyp    = structField.Type
		)
		if !structField.IsExported() {
			continue
		}

		if fieldTyp.Kind() == reflect.Pointer {
			fieldTyp = fieldTyp.Elem()
		}

		var (
			fieldVal = reflect.New(fieldTyp)
			gormTg   = tag.New(structField.Tag.Get("gorm"))
			jsonTg   = tag.NewJsonTag(structField.Tag.Get("json"))
		)

		tbField.Name = gormTg.Get("column")
		tbField.Title = gormTg.Get("comment")
		tbField.Validate = structField.Tag.Get("binding")

		if gormTg.Get("-") != "" {
			continue
		}
		if tbField.Name == "" {
			tbField.Name = jsonTg.GetName()
		}
		if tbField.Name == "" {
			tbField.Name = utils.HumpToUnderline(structField.Name)
		}

		//自定义的包装类型，基本都实现了这个接口
		if x, ok := fieldVal.Interface().(schema.GormDataTypeInterface); ok {
			tbField.Type = x.GormDataType()
			tableFieldList = append(tableFieldList, tbField)
			continue
		}
		if gormTg.Get("serializer") != "" {
			tbField.Type = ctype.TpJson
			tableFieldList = append(tableFieldList, tbField)
			continue
		}

		switch fieldTyp.Kind() {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
			tbField.Type = ctype.TpInteger
		case reflect.Float32, reflect.Float64:
			tbField.Type = ctype.TpNumeric
		case reflect.String:
			tbField.Type = ctype.TpString
		case reflect.Bool:
			tbField.Type = ctype.TpBool
		case reflect.Map, reflect.Slice, reflect.Array:
		case reflect.Struct:
			if gormTg.Get("column") != "" { //理论上不会，只能是理解为默认json了
				tbField.Type = ctype.TpJson
				tableFieldList = append(tableFieldList, tbField)
				continue
			}

			pre := gormTg.Get("embeddedPrefix")
			if pre != "" {
				tf := GetTableFields(fieldVal.Interface())
				for _, v := range tf {
					v.Name = pre + v.Name
					tableFieldList = append(tableFieldList, v)
				}
			} else if structField.Anonymous {
				tf := GetTableFields(fieldVal.Interface())
				tableFieldList = append(tableFieldList, tf...)
			}
		default:
		}
	}
	tableFieldCache[tp] = tableFieldList
	return tableFieldList
}
