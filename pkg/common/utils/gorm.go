package utils

import (
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/duke-git/lancet/v2/slice"
	"gorm.io/gorm/schema"
	"reflect"
	"sync"
)

// DbFieldNameMap 实体类名称Map。dbName:FieldName
func DbFieldNameMap(eType reflect.Type, nameStrategy schema.Namer) map[string]string {
	var (
		val  = reflect.New(eType)
		vs   = val.Interface()
		sMap sync.Map
		mp   = make(map[string]string)
	)
	parse, err1 := schema.Parse(vs, &sMap, nameStrategy)
	if err1 != nil {
		logger.Errorf("parse error: %v", err1)
		return mp
	}
	slice.ForEach(parse.Fields, func(index int, item *schema.Field) {
		mp[item.DBName] = item.Name
	})
	return mp
}
