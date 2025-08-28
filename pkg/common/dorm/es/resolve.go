package es

import (
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/dorm/es/es_api"
	"github.com/davycun/eta/pkg/common/tag"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/duke-git/lancet/v2/maputil"
	"reflect"
	"strings"
)

func GetEsMapping(ent interface{}) map[string]interface{} {
	return field2EsProps(getStructEsFields(reflect.TypeOf(ent))...)
}

func getStructEsFields(tp reflect.Type) []reflect.StructField {
	var (
		rs = make([]reflect.StructField, 0, 10)
	)

	switch tp.Kind() {
	case reflect.Pointer:
		return getStructEsFields(tp.Elem())
	case reflect.Slice:
		return getStructEsFields(tp.Elem())
	case reflect.Struct:
		for i := 0; i < tp.NumField(); i++ {
			fd := tp.Field(i)
			get := fd.Tag.Get(tag.EsTagName)
			if get == "" {
				//没有编写es tag的字段表示丢弃或者是组合了其他结构体
				if fd.Type.Kind() == reflect.Struct {
					rs = append(rs, getStructEsFields(fd.Type)...)
				}
				continue
			}
			rs = append(rs, tp.Field(i))
		}
	default:

	}
	return rs
}

func field2EsProps(sfs ...reflect.StructField) map[string]interface{} {
	var (
		rs = make(map[string]interface{})
	)
	for _, fd := range sfs {
		var (
			name = getFieldEsName(fd)
		)

		if name == "" {
			continue
		}

		tgTxt := fd.Tag.Get(tag.EsTagName)
		if tgTxt == "" {
			continue
		}
		//tg := tag.ParseTagDefault(EsTagName, tgTxt)
		tg := tag.ParseEsTag(tgTxt)
		if tg.Exists("ignore") {
			continue
		}

		mp := tg.GetAll()
		kmp := make(map[string]interface{})

		//kmp["type"] = tg.Get("type")
		for x, y := range mp {
			switch strings.TrimSpace(strings.ToLower(y)) {
			case "true":
				kmp[x] = true
				continue
			case "false":
				kmp[x] = false
				continue
			}
			kmp[x] = y
		}

		if tg.Get("type") == "nested" || tg.Get("type") == "object" {
			//如果是json字段，那么就存储raw json数据，而不进行动态扩展字段，比如BaseEntity中的Extra或者EtlExtra等
			//主要是避免一些自动生成的类型与实际数据不匹配，比如日期类型，这样会导致插入数据错误
			//不允许动态扩展字段之后，带来的不便就是这个字段内的字段不允许再被filter了
			if fd.Type == reflect.TypeFor[ctype.Json]() || fd.Type == reflect.TypeFor[*ctype.Json]() {
				kmp["dynamic"] = false
			} else {
				kmp["properties"] = field2EsProps(getStructEsFields(fd.Type)...)
			}
		} else {
			// 所有字段都加一个 keyword 属性
			if !maputil.HasKey(kmp, "fields") {
				kmp["fields"] = make(map[string]interface{})
			}
			if f, ok := kmp["fields"].(map[string]interface{}); ok {
				if !maputil.HasKey(f, "keyword") {
					kmp["fields"].(map[string]interface{})["keyword"] = es_api.DefaultKeyword()
				}
			}
		}
		rs[name] = kmp
	}

	return rs
}

func getFieldEsName(sf reflect.StructField) string {
	var (
		name string
	)
	tg := sf.Tag.Get("json")
	if tg == "" {
		return utils.HumpToUnderline(sf.Name)
	}
	ss := strings.Split(tg, ",")
	for _, v := range ss {
		switch v {
		case "omitempty":
			continue
		default:
			name = strings.TrimSpace(v)
			break
		}
	}
	if name == "" {
		return utils.HumpToUnderline(sf.Name)
	}

	return name
}
