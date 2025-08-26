package entity

import (
	"errors"
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/dynamicstruct"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/duke-git/lancet/v2/maputil"
	"gorm.io/gorm"
	"reflect"
	"strings"
)

var (
	LocatedApp   int = 1                         //针对的是Table.Located 的常量，表示对应实体表将创建在APP库（默认的情况）
	LocatedLocal int = 2                         //针对的是Table.Located 的常量，表示对应实体表将创建在Local（平台/本地）库
	LocatedAll   int = LocatedApp | LocatedLocal //针对的是Table.Located 的常量，表在平台库和APP库都创建
)

type TableIndex struct {
	Fields []string `json:"fields,omitempty"` //索引涉及到的字段
	Type   string   `json:"type,omitempty"`   //B-List、Hash、GiST、SP-GiST、Gin、BRIN，这个是针对postgresql的，这个根据不同的数据要注意分别有哪些类型
	Class  string   `json:"class,omitempty"`  //索引类别，比如UNIQUE，或者达梦的[CONTEXT、BITMAP、ARRAY、SPATIAL]
	Option string   `json:"option,omitempty"`
}

// ValidType
// 校验索引类型是否合法
func (idx TableIndex) ValidType() bool {
	//达梦支持
	return true
}

// Build
// 构建索引创建语句
// PG支持的Type：BTREE、HASH、GIST、GIN
// 达梦：暂不支持
// MYSQL：BTREE、HASH
func (idx TableIndex) Build(db *gorm.DB, tableName string) (string, error) {
	var (
		dbType = dorm.GetDbType(db)
		bd     = strings.Builder{}
	)
	//TODO Class 要区别对待不同的数据库
	bd.WriteString(fmt.Sprintf(`CREATE %s INDEX  %s ON %s`,
		idx.Class, dorm.Quote(dbType, idx.IndexName(tableName)), dorm.GetScmTableName(db, tableName)))
	if idx.Type != "" {
		switch dbType {
		case dorm.PostgreSQL:
			bd.WriteString(" USING " + idx.Type + " ")
			bd.WriteString(fmt.Sprintf(`(%s)`, dorm.JoinColumns(dbType, "", idx.Fields)))
		case dorm.DaMeng:
			//not support
			bd.WriteString(fmt.Sprintf(`(%s)`, dorm.JoinColumns(dbType, "", idx.Fields)))
		case dorm.Mysql:
			if !strings.EqualFold(idx.Type, "BTREE") && !strings.EqualFold(idx.Type, "HASH") {
				return "", errors.New("mysql only support BTREE and HASH index type")
			}
			bd.WriteString(fmt.Sprintf(`(%s)`, dorm.JoinColumns(dbType, "", idx.Fields)))
			bd.WriteString(" USING " + idx.Type + " ")
		default:
			//not support
		}
	} else {
		bd.WriteString(fmt.Sprintf(`(%s)`, dorm.JoinColumns(dbType, "", idx.Fields)))
	}
	bd.WriteString(idx.Option)
	return bd.String(), nil
}
func (idx TableIndex) IndexName(tbName string) string {
	bd := strings.Builder{}
	bd.WriteString("idx_" + tbName)
	for _, v := range idx.Fields {
		bd.WriteString("_" + v)
	}
	return bd.String()
}

type TableField struct {
	Name       string `json:"name,omitempty" binding:"required"`    //字段名称
	Title      string `json:"title,omitempty" binding:"required"`   //字段标题
	Type       string `json:"type,omitempty" binding:"required"`    //字段类型 参照ctype.DbTypeMap
	Comment    string `json:"comment,omitempty" binding:"required"` //注释
	Default    string `json:"default,omitempty"`                    //默认值
	GormTag    string `json:"gorm_tag,omitempty"`
	BindingTag string `json:"binding_tag,omitempty"` //验证规则(github.com/go-playground/validator)，会添加到结构体的tag中，比如 binding:"required"
	EsTag      string `json:"es_tag,omitempty"`      //es的tag内容，比如type:text;analyzer:digit_analyzer 应该需要结构化，暂不适用此字段
}

type Table struct {
	Feature
	Fields     []TableField           `json:"fields,omitempty" binding:"required,dive"` //当前表的字段
	Indexes    []TableIndex           `json:"indexes,omitempty"`                        //表的索引
	TableName  string                 `json:"table_name,omitempty"`                     //表名
	EntityType reflect.Type           `json:"entity_type,omitempty" gorm:"-:all"`       //操作的实体结构体的类型
	Order      int                    `json:"order,omitempty"`                          //这个主要的作用是在Migrator的时候可能会有优先级问题，值越大优先级越高
	Options    map[dorm.DbType]string `json:"options,omitempty" gorm:"-:all"`           //创建表的时候的一些选项，比如表空间，表引擎等
	Located    int                    `json:"located,omitempty" gorm:"-:all"`           // 表示实体表会创建的位置 ,默认是0即 创建在APP
	//EnableDbType []dorm.DbType          `json:"enable_db_type,omitempty" gorm:"-:all"`    //那些实体需要再哪些类型的DB上创建对应的表
	EsEnable     ctype.Boolean          `json:"es_enable,omitempty"`                //是否有对应的ES index
	EsSettings   map[string]interface{} `json:"es_settings,omitempty" gorm:"-:all"` //主要是给ES的 index的setting用
	EsFields     []TableField           `json:"es_fields,omitempty"`
	EsEntityType reflect.Type           `json:"es_entity_type,omitempty"` //如果有es的实体，则该字段为es的实体类型，否则为nil
}

func (t *Table) LocatedApp() bool {
	return t.Located&LocatedApp == LocatedApp || t.Located == 0 //没有配置，默认表就放在app下
}
func (t *Table) LocatedLocal() bool {
	return t.Located&LocatedLocal == LocatedLocal
}
func (t *Table) LocatedAll() bool {
	return t.Located == LocatedAll
}

func (t *Table) GetTableName() string {
	if t.TableName == "" {
		t.TableName = GetTableName(t.NewEntityPointer())
	}
	return t.TableName
}
func (t *Table) GetFields() []TableField {
	if len(t.Fields) > 0 {
		return t.Fields
	}
	t.Fields = GetTableFields(t.NewEntityPointer())
	return t.Fields
}
func (t *Table) EsEnabled() bool {
	if ctype.IsValid(t.EsEnable) {
		return ctype.Bool(t.EsEnable)
	}
	obj := t.NewEsEntityPointer()
	if x, ok := obj.(EsInterface); ok {
		t.EsEnable = ctype.NewBoolean(x.EsEnable(), true)
	}
	return ctype.Bool(t.EsEnable)
}
func (t *Table) EsRetrieveEnabled() bool {
	return t.EsEnabled() && !ctype.Bool(t.DisableRetrieveEs)
}
func (t *Table) UseParamAuth() bool {
	return ctype.Bool(t.ParamAuth)
}

func (t *Table) NewEntityPointer() any {
	if t.EntityType == nil {
		return t.newEntityOrSlicePointerFromFields(false, false)
	}
	return reflect.New(t.EntityType).Interface()
}
func (t *Table) NewEsEntityPointer() any {
	if t.EsEntityType == nil {
		s := t.newEntityOrSlicePointerFromFields(false, true)
		if s == nil {
			return t.NewEntityPointer()
		}
	}
	return reflect.New(t.EsEntityType).Interface()
}
func (t *Table) NewEntitySlicePointer() any {
	if t.EntityType == nil {
		return t.newEntityOrSlicePointerFromFields(true, false)
	}
	return reflect.New(reflect.SliceOf(t.EntityType)).Interface()
}
func (t *Table) NewEsEntitySlicePointer() any {
	if t.EsEntityType == nil {
		obj := t.newEntityOrSlicePointerFromFields(true, true)
		if obj == nil {
			return t.NewEntitySlicePointer()
		}
	}
	return reflect.New(reflect.SliceOf(t.EsEntityType)).Interface()
}

func (t *Table) newEntityOrSlicePointerFromFields(batch bool, isEs bool) any {

	if len(t.Fields) < 1 && len(t.EsFields) < 1 {
		return nil
	}

	var (
		bd1       = dynamicstruct.ExtendStruct(BaseEntity{})
		fieldList = t.Fields
	)

	//在构建ES结构体的时候，如果不存在ES字段，则使用Entity的字段
	if isEs && len(t.EsFields) > 0 {
		fieldList = t.EsFields
	}

	bd2 := NewStructBuilder(true, fieldList...)
	bd1.Merge(bd2)
	for _, v := range t.SignFields {
		if v.VerifyField == "" {
			continue
		}
		bd1.AddField(utils.Column2StructFieldName(v.VerifyField), new(bool), fmt.Sprintf(`json:"%s,omitempty" gorm:"-:all"`, v.VerifyField))
	}

	obj := bd1.Build().New()
	if t.EntityType == nil && !isEs {
		t.EntityType = reflect.TypeOf(obj).Elem()
	}

	if t.EsEntityType == nil && len(t.EsFields) < 1 {
		t.EsEntityType = reflect.TypeOf(obj).Elem()
	}

	if batch {
		return bd1.Build().NewSliceOfStructs()
	}
	return obj
}

func (t *Table) Merge(tb *Table) {
	if tb == nil {
		return
	}
	if tb.TableName != "" {
		t.TableName = tb.TableName
	}
	if tb.TableName != "" {
		t.TableName = tb.TableName
	}
	if tb.EntityType != nil {
		t.EntityType = tb.EntityType
	}
	if len(tb.Fields) > 0 {
		t.Fields = utils.Merge(t.Fields, tb.Fields...)
	}
	//TODO 待深度合并
	if len(tb.Indexes) > 0 {
		t.Indexes = tb.Indexes
	}
	if tb.EsEntityType != nil {
		t.EsEntityType = tb.EsEntityType
	}
	if ctype.IsValid(tb.EsEnable) {
		t.EsEnable = tb.EsEnable
	}
	if len(tb.EsFields) > 0 {
		t.EsFields = utils.Merge(t.EsFields, tb.EsFields...)
	}
	if len(tb.EsSettings) > 0 {
		t.EsSettings = maputil.Merge(t.EsSettings, tb.EsSettings)
	}

	t.Feature = t.Feature.Merge(tb.Feature)
}

func LoadTable(obj any) Table {
	var (
		tb = Table{}
	)
	if obj == nil {
		return tb
	}
	tb.EntityType = utils.GetRealType(reflect.TypeOf(obj))
	tb.TableName = GetTableName(obj)
	tb.Feature = LoadFeature(obj)
	tb.Fields = GetTableFields(obj)

	return tb
}
