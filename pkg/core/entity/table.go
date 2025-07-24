package entity

import (
	"errors"
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/utils"
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
		idx.Class, dorm.Quote(dbType, idx.IndexName(tableName)), dorm.GetDbTable(db, tableName)))
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
	Name     string `json:"name,omitempty" binding:"required"`    //字段名称
	Title    string `json:"title,omitempty" binding:"required"`   //字段标题
	Type     string `json:"type,omitempty" binding:"required"`    //字段类型 参照ctype.DbTypeMap
	Comment  string `json:"comment,omitempty" binding:"required"` //注释
	Validate string `json:"validate,omitempty"`                   //验证规则(github.com/go-playground/validator)，会添加到结构体的tag中，比如 "required"
	Default  string `json:"default,omitempty"`                    //默认值
}

type Table struct {
	Feature
	Fields       []TableField           `json:"fields,omitempty" binding:"required,dive"` //当前表的字段
	Indexes      []TableIndex           `json:"index,omitempty"`                          //表的索引
	TableName    string                 `json:"table_name,omitempty"`                     //表名
	EsEntityType reflect.Type           `json:"es_entity_type,omitempty"`                 //如果有es的实体，则该字段为es的实体类型，否则为nil
	EntityType   reflect.Type           `json:"entity_type,omitempty" gorm:"-:all"`       //操作的实体结构体的类型
	Order        int                    `json:"order,omitempty"`                          //这个主要的作用是在Migrator的时候可能会有优先级问题，值越大优先级越高
	Options      map[dorm.DbType]string `json:"options,omitempty" gorm:"-:all"`           //创建表的时候的一些选项，比如表空间，表引擎等
	Settings     map[string]interface{} `json:"settings,omitempty" gorm:"-:all"`          //主要是给ES的 index的setting用
	Located      int                    `json:"located,omitempty" gorm:"-:all"`           // 表示实体表会创建的位置 ,默认是0即 创建在APP
	EnableDbType []dorm.DbType          `json:"enable_db_type,omitempty" gorm:"-:all"`    //那些实体需要再哪些类型的DB上创建对应的表
}

func (t *Table) LocatedApp() bool {
	return t.Located&LocatedApp == LocatedApp || t.Located == 0 //没有配置，默认表就放在app下
}
func (t *Table) LocatedLocal() bool {
	return t.Located&LocatedLocal == LocatedLocal
}

func (t *Table) GetTableName() string {
	if t.TableName == "" && t.EntityType != nil {
		obj := reflect.New(t.EntityType)
		t.TableName = GetTableName(obj.Interface())
	}
	return t.TableName
}

// GetEsIndexName
// 返回es的索引名称，真正的索引是需要加上schema的前缀的
func (t *Table) GetEsIndexName() string {
	return t.GetTableName()
}
func (t *Table) NewEntityPointer() any {
	if t.EntityType == nil {
		return nil
	}
	return reflect.New(t.EntityType).Interface()
}
func (t *Table) NewEntitySlicePointer() any {
	if t.EntityType == nil {
		return nil
	}
	return reflect.New(reflect.SliceOf(t.EntityType)).Interface()
}
func (t *Table) EnableRetrieveEs() bool {
	return !ctype.Bool(t.Feature.DisableRetrieveEs) && ctype.Bool(t.EsEnable)
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
	t.Feature = t.Feature.Merge(tb.Feature)
}
