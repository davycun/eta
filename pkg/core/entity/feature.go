package entity

import (
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"reflect"
)

// HistoryInterface
// Entity通过实现这个接口可自动进行数据历史操作记录的记录
type HistoryInterface interface {
	History() bool
}

// FieldUpdater
// Entity通过实现这个接口来实现字段更新人的自动记录（trigger实现）
type FieldUpdater interface {
	FieldUpdater() bool
}

// RaInterface
// 实体实现，构建“content”和“detail_content”字段的值
// 暂时只是支持本实体的字段的全文检索，如果需要自定义内容，可以借助ES新增宽表实体实现
type RaInterface interface {
	// RaDbFields 构建“ra_content”的数据库字段
	RaDbFields() []string
}

// SignInterface
// Entity通过实现这个接口来实现自动签名
type SignInterface interface {
	SignFields() []SignFieldsInfo
}

// CryptInterface
// Entity通过实现这个接口来实现自动加解密
type CryptInterface interface {
	CryptFields() []CryptFieldInfo
}

// EsInterface
// Entity实现这个接口，表示会自动同步数据到ES
type EsInterface interface {
	EsEnable() bool
}

// WideInterface
// Entity实现这个接口，表示会自动同步数据到ES
type WideInterface interface {
	WideEsIndexName() string
	WideTableName() string
}

// Embedded
// 在通过From表加载关系和实体的时候，同在定义了一个关系实体，然后通过join查询，通过embedded给列名前端才能组装成实体对象
// 实现这个接口就可以再Loader中自动获取，比如RelationAddr，RelationPeople等
type Embedded interface {
	EmbeddedPrefix() string
}

// Feature
// 如果是动态表，可以通过一下信息来记录动态表的一些特性
type Feature struct {
	dorm.JsonType
	History           ctype.Boolean    `json:"history,omitempty"`       //表示是否启用History
	FieldUpdater      ctype.Boolean    `json:"field_updater,omitempty"` //表示是否启用FieldUpdater（字段更新人）
	DisableRetrieveEs ctype.Boolean    `json:"retrieve_with_es"`
	CryptFields       []CryptFieldInfo `json:"crypt_fields,omitempty"`
	SignFields        []SignFieldsInfo `json:"sign_fields,omitempty"`
	RaDbFields        []string         `json:"ra_db_fields,omitempty"`
	EsEnable          ctype.Boolean    `json:"es_enable,omitempty"` //是否有对应的ES index
	EsExtraFields     []TableField     `json:"es_extra_fields"`     //除了常规的table的字段外，ES需要更多的字段（适用于宽表)
}

func (f Feature) Merge(ft Feature) Feature {
	if ctype.IsValid(ft.History) {
		f.History = ft.History
	}
	if ctype.IsValid(ft.FieldUpdater) {
		f.FieldUpdater = ft.FieldUpdater
	}
	if ctype.IsValid(ft.DisableRetrieveEs) {
		f.DisableRetrieveEs = ft.DisableRetrieveEs
	}
	if len(ft.SignFields) > 0 {
		f.SignFields = ft.SignFields
	}
	if len(ft.CryptFields) > 0 {
		f.CryptFields = ft.CryptFields
	}
	if len(ft.RaDbFields) > 0 {
		f.RaDbFields = ft.RaDbFields
	}
	if ctype.Bool(ft.EsEnable) {
		f.EsEnable = ft.EsEnable
	}
	if len(ft.EsExtraFields) > 0 {
		f.EsExtraFields = ft.EsExtraFields
	}
	return f
}

type SignFieldsInfo struct {
	Enable      bool     `json:"enable,omitempty"`                    //可以控制是否开启，可以提前做一些配置加签，但是不开启
	Algo        string   `json:"algo,omitempty" binding:"required"`   // 算法，hmac_sm3 hmac_sha256 md5
	Key         string   `json:"key,omitempty" binding:"required"`    // 加签密钥
	Field       string   `json:"field,omitempty" binding:"required"`  // 存储签名的字段，这个字段需要是表里已经存在的
	Fields      []string `json:"fields,omitempty" binding:"required"` // 需要加签的字段
	VerifyField string   `json:"verify_field,omitempty"`              // 校验结果字段。这个字段不为空，则在查询这条数据的时候，会额外增加一个字段表示验签结果，在返回结果中体现
}

type CryptFieldInfo struct {
	Enable        bool     `json:"enable,omitempty"`                        //可以控制是否开启加密，为了能提前配置加密字段，但是不
	Algo          string   `json:"algo,omitempty" binding:"required"`       // 算法
	SecretKey     []string `json:"secret_key,omitempty" binding:"required"` // 密钥，16字节
	Field         string   `json:"field,omitempty" binding:"required"`      // 需要加密的字段，这个字段需要是 Field 里已定义的
	KeepTxtPreCnt int      `json:"keep_txt_pre_cnt,omitempty"`              // 保持文本前几位明文
	KeepTxtSufCnt int      `json:"keep_txt_suf_cnt,omitempty"`              // 保持文本后几位明文
	SliceSize     int      `json:"slice_size,omitempty"`                    // 切片加密的切片大小，为0时不切片。默认0
}

type SignInfoList []SignFieldsInfo
type CryptInfoList []CryptFieldInfo

func (d SignInfoList) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return dorm.JsonGormDBDataType(db, field)
}
func (d SignInfoList) GormDataType() string {
	return dorm.JsonGormDataType()
}
func (d CryptInfoList) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return dorm.JsonGormDBDataType(db, field)
}
func (d CryptInfoList) GormDataType() string {
	return dorm.JsonGormDataType()
}

func LoadFeature(obj any) Feature {
	var (
		ft = Feature{}
	)
	if obj == nil {
		return ft
	}
	if _, ok := obj.(HistoryInterface); ok {
		ft.History = ctype.NewBoolean(true, true)
	}
	if _, ok := obj.(FieldUpdater); ok {
		ft.FieldUpdater = ctype.NewBoolean(true, true)
	}
	if x, ok := obj.(CryptInterface); ok {
		ft.CryptFields = x.CryptFields()
	}
	if x, ok := obj.(SignInterface); ok {
		ft.SignFields = x.SignFields()
	}
	if x, ok := obj.(RaInterface); ok {
		ft.RaDbFields = x.RaDbFields()
	}
	if x, ok := obj.(EsInterface); ok {
		ft.EsEnable = ctype.NewBoolean(x.EsEnable(), true)
	}
	return ft
}

func LoadTable(obj any) Table {
	var (
		tb = Table{}
	)
	if obj == nil {
		return tb
	}
	tp := reflect.TypeOf(obj)
	if tp.Kind() == reflect.Pointer {
		tp = tp.Elem()
	}
	tb.EntityType = tp
	tb.TableName = GetTableName(obj)
	tb.Feature = LoadFeature(obj)
	tb.Fields = GetTableFields(obj)

	return tb
}
