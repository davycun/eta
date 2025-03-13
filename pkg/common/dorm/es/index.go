package es

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/utils"
	"gorm.io/gorm/schema"
	"reflect"
)

func GetIndexName(scm string, obj any) string {

	var (
		tbName = ""
	)

	switch x := obj.(type) {
	case string:
		tbName = x
	case *string:
		tbName = *x
	case schema.TablerWithNamer:
		tbName = x.TableName(nil)
	case schema.Tabler:
		tbName = x.TableName()
	default:
		tbName = utils.HumpToUnderline(reflect.TypeOf(obj).Name())

	}

	if scm != "" && tbName != "" {
		return fmt.Sprintf("%s_%s", scm, tbName)
	}
	return tbName
}
