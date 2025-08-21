package entity

import (
	"github.com/davycun/eta/pkg/common/caller"
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

// GetValue
// 获取struct或者map中指定key的Value
func GetValue(val reflect.Value, key string) reflect.Value {

	if !val.IsValid() {
		return reflect.Value{}
	}
	val = utils.GetRealValue(val)

	var (
		fieldName = utils.UnderlineToHump(key) //TODO 理论上要根据jsonTag或者gormTag来确定名字
		jsonName  = utils.HumpToUnderline(key)
		valType   = val.Type()
	)

	switch valType.Kind() {
	case reflect.Pointer:
		return GetValue(val.Elem(), key)
	case reflect.Struct:
		//包装的结构体，无需地推查找，直接返回
		fieldValue := val.FieldByName(fieldName)
		if fieldValue.IsValid() {
			return fieldValue
		}
		for i := 0; i < val.NumField(); i++ {
			var (
				fieldVal  = val.Field(i)
				field     = valType.Field(i)
				fieldType = utils.GetRealType(field.Type)
				gormTag   = tag.ParseGormTag(field.Tag.Get(tag.GormTagName))
			)
			if jsonName == tag.ParseJsonTag(field.Tag.Get(tag.JsonTagName)).GetFirstKey() || gormTag.Get("column") == jsonName {
				return fieldVal
			}
			if fieldType.Kind() == reflect.Struct && (field.Anonymous || gormTag.Exists("embedded")) {
				//只有组合字段才继续查找
				fv := GetValue(fieldVal, jsonName)
				if fv.IsValid() {
					return fv
				}
			}
		}
	case reflect.Map:
		if !val.CanInterface() {
			return reflect.Value{}
		}
		var (
			valInter = val.Interface()
		)

		if x, ok := valInter.(Getter); ok {
			if v, ok1 := x.Get(jsonName); ok1 {
				return reflect.ValueOf(v)
			}
		}
		switch x := valInter.(type) {
		case ctype.Map:
			if v, ok := x[jsonName]; ok {
				return reflect.ValueOf(v)
			}
		case *ctype.Map:
			if v, ok := (*x)[jsonName]; ok {
				return reflect.ValueOf(v)
			}
		case gin.H:
			if v, ok := x[jsonName]; ok {
				return reflect.ValueOf(v)
			}
		case *gin.H:
			if v, ok := (*x)[jsonName]; ok {
				return reflect.ValueOf(v)
			}
		case map[string]any:
			if v, ok := x[jsonName]; ok {
				return reflect.ValueOf(v)
			}
		case *map[string]any:
			if v, ok := (*x)[jsonName]; ok {
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
			gormTg    = tag.ParseGormTag(structField.Tag.Get(tag.GormTagName))
			jsonTg    = tag.ParseJsonTag(structField.Tag.Get(tag.JsonTagName))
			fieldName = structField.Name
			colName   = gormTg.Get("column")
		)

		if gormTg.Exists("-") {
			continue
		}
		if colName == "" {
			colName = jsonTg.GetFirstKey()
		}
		if colName == "" {
			colName = utils.HumpToUnderline(structField.Name)
		}

		//自定义的包装类型，基本都实现了这个接口
		if _, ok := fieldVal.Interface().(schema.GormDataTypeInterface); ok {
			gormField[fieldName] = colName
			continue
		}
		if gormTg.Exists("serializer") {
			gormField[fieldName] = colName
			continue
		}

		switch fieldTyp.Kind() {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int,
			reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint, reflect.Float32, reflect.Float64, reflect.String, reflect.Bool:
			gormField[fieldName] = colName
		case reflect.Map, reflect.Slice, reflect.Array:
		case reflect.Struct:
			if gormTg.Exists("column") {
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
	if obj == nil {
		return []TableField{}
	}

	switch x := obj.(type) {
	case reflect.Type:
		tp = utils.GetRealType(x)
	case reflect.Value:
		tp = utils.GetRealType(x.Type())
	case *reflect.Value:
		tp = utils.GetRealType(x.Type())
	default:
		tp = utils.GetRealType(reflect.TypeOf(obj))
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
			fieldTyp    = utils.GetRealType(structField.Type)
		)
		if !structField.IsExported() {
			continue
		}

		var (
			fieldVal = reflect.New(fieldTyp)
			gormTg   = tag.ParseGormTag(structField.Tag.Get(tag.GormTagName))
			jsonTg   = tag.ParseJsonTag(structField.Tag.Get(tag.JsonTagName))
		)

		tbField.Name = gormTg.Get("column")
		tbField.Title = gormTg.Get("comment")
		tbField.BindingTag = structField.Tag.Get(tag.BindingTagName)

		if gormTg.Exists("-") {
			continue
		}
		if tbField.Name == "" {
			tbField.Name = jsonTg.GetFirstKey()
		}
		if tbField.Name == "" {
			tbField.Name = utils.HumpToUnderline(structField.Name)
		}

		caller.NewCaller().
			Call(func(cl *caller.Caller) error {
				if x := gormTg.Get("type"); x != "" {
					tbField.Type = x
					cl.Stop()
				}
				return nil
			}).
			Call(func(cl *caller.Caller) error {
				//自定义的包装类型，基本都实现了这个接口
				if x, ok := fieldVal.Interface().(schema.GormDataTypeInterface); ok {
					tbField.Type = x.GormDataType()
					cl.Stop()
				}
				return nil
			}).
			Call(func(cl *caller.Caller) error {
				if gormTg.Exists("serializer") {
					tbField.Type = ctype.TypeJsonName
					cl.Stop()
				}
				return nil
			}).
			Call(func(cl *caller.Caller) error {
				switch fieldTyp.Kind() {
				case reflect.Int8, reflect.Int16, reflect.Uint8, reflect.Uint16:
					tbField.Type = ctype.TypeIntegerName
				case reflect.Int32, reflect.Int64, reflect.Int, reflect.Uint32, reflect.Uint64, reflect.Uint:
					tbField.Type = ctype.TypeBigIntegerName
				case reflect.Float32, reflect.Float64:
					tbField.Type = ctype.TypeNumericName
				case reflect.String:
					tbField.Type = ctype.TypeStringName
				case reflect.Bool:
					tbField.Type = ctype.TypeBoolName
				case reflect.Map, reflect.Slice, reflect.Array:
				case reflect.Struct:
					if gormTg.Exists("column") { //理论上不会，只能是理解为默认json了
						tbField.Type = ctype.TypeJsonName
						cl.Stop()
					} else {
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
					}
				default:
				}
				return nil
			})
		//有可能类型是embeddedPrefix，所以当前tbField.Type为空，则跳过
		if tbField.Type != "" {
			tableFieldList = append(tableFieldList, tbField)
		}
	}
	tableFieldCache[tp] = tableFieldList
	return tableFieldList
}
