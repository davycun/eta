package entity

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/dynamicstruct"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/tag"
	"github.com/davycun/eta/pkg/common/utils"
)

func NewStructOrSlicePointer(batch bool, gormEnable bool, tbFields ...TableField) any {
	builder := NewStructBuilder(gormEnable, tbFields...)
	if builder == nil {
		return nil
	}
	if batch {
		e := builder.Build().NewSliceOfStructs()
		return e
	} else {
		obj := builder.Build().New()
		return obj
	}
}

func NewStructBuilder(gormEnable bool, tbFields ...TableField) dynamicstruct.Builder {
	var (
		builder = dynamicstruct.NewStruct()
	)
	for _, v := range tbFields {
		var (
			jsonTag   = tag.NewJsonTag().Add(v.Name, "").Add("omitempty", "")
			gormTag   = tag.ParseGormTag(v.GormTag)
			bindTag   = tag.ParseBindingTag(v.BindingTag)
			esTag     = tag.ParseEsTag(v.EsTag)
			fieldName = utils.Column2StructFieldName(v.Name)
		)

		if gormEnable {
			if !gormTag.Exists("column") {
				gormTag.Add("column", v.Name)
			}
			if !gormTag.Exists("comment") {
				gormTag.Add("comment", v.Comment)
			}
			if v.Type == ctype.TypeJsonName {
				gormTag.Add("serializer", "json")
			}
		} else {
			gormTag.Add("-", "all")
		}

		if !esTag.Exists("type") {
			esType, _ := ctype.GetEsType(v.Type)
			if esType == "" {
				logger.Errorf("not support es type: %s", v.Type)
			} else {
				esTag.Add("type", esType)
			}
		}
		tg := fmt.Sprintf("%s %s %s %s", jsonTag.String(), gormTag.String(), esTag.String(), bindTag.String())
		builder.AddField(fieldName, ctype.NewTypeValue(v.Type, true), tg)
	}
	return builder
}
