package template

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/eta/constants"
	dynamicstruct "github.com/ompluscator/dynamic-struct"
	"gorm.io/gorm"
	"reflect"
	"strings"
)

// GetTable
// 获取实际模版数据的表名
func (p *Template) GetTable() *entity.Table {
	if p.Table.TableName == "" && p.Code != "" {
		p.Table.TableName = constants.TableTemplatePrefix + p.Code
	}
	if p.Table.EntityType == nil {
		p.Table.EntityType = p.GetEntityType()
	}
	return &p.Table
}

func (p Template) GetEntityType() reflect.Type {
	if len(p.Table.Fields) < 1 {
		return nil
	}
	return reflect.TypeOf(p.GetEntity(false)).Elem()
}
func (p Template) GetRsDataType() reflect.Type {
	return p.GetEntityType()
}

func (p Template) GetEntity(batch bool) any {
	var (
		//TODO 这里应该考虑下template 的EntityType 来确定用BaseEntity 还是BaseEdgeEntity
		builder = dynamicstruct.ExtendStruct(entity.BaseEntity{})
	)
	for _, v := range p.Table.Fields {
		var (
			tag       = fmt.Sprintf(`json:"%s,omitempty" gorm:"column:%s;" %s`, v.Name, v.Name, getBinding(v.Validate))
			fieldName = constants.Column2StructFieldName(v.Name)
		)
		dbType := strings.ToLower(v.Type)
		builder.AddField(fieldName, ctype.NewTypeValue(dbType, true), tag)
	}
	for _, v := range p.Table.SignFields {
		vType := "bool"
		var (
			vName     = v.VerifyField
			tag       = fmt.Sprintf(`json:"%s,omitempty" gorm:"-:all"`, vName)
			fieldName = constants.Column2StructFieldName(vName)
		)
		if vName == "" {
			continue
		}
		dbType := strings.ToLower(vType)
		builder.AddField(fieldName, ctype.NewTypeValue(dbType, false), tag)
	}

	// 为了能从实体获取表名
	//builder.AddField(constants.TemplateTableNameField, "", `json:"table_name,omitempty" gorm:"-:all"`)
	// 为了能从实体获获取到 RaDbFields
	//builder.AddField(constants.TemplateRaDbFields, []string{}, `json:"ra_db_fields,omitempty" gorm:"-:all"`)

	if batch {
		e := builder.Build().NewSliceOfStructs()
		return e
	} else {
		obj := builder.Build().New()
		return obj
	}
}

func getGormType(db *gorm.DB, tp string) (string, error) {
	if tp == ctype.TypeJsonName {
		return "serializer:json", nil
	}
	rs, err := ctype.GetDbTypeName(db, tp)
	return "type:" + rs, err
}

func getBinding(bind string) string {
	if bind != "" {
		return fmt.Sprintf(`binding:"%s"`, bind)
	}
	return ""
}
