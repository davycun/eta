package entity

import (
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/duke-git/lancet/v2/maputil"
)

const defaultCryptKey = "0123456789abcdef"

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
	DisableRetrieveEs ctype.Boolean    `json:"retrieve_with_es"`        //禁止从ES查询数据
	ParamAuth         ctype.Boolean    `json:"param_auth"`              //默认是false，也就是需要权限，如果设置为true。那么就会根据参数（DisablePermFilter）决定是否需要权限
	CryptFields       []CryptFieldInfo `json:"crypt_fields,omitempty"`
	SignFields        []SignFieldsInfo `json:"sign_fields,omitempty"`
	RaEnable          ctype.Boolean    `json:"ra_enable,omitempty"` //表示是否启用RA，是的话针对dameng会自动创建所以，后续会自动填充内容
	RaDbFields        []string         `json:"ra_db_fields,omitempty"`
}

func (f Feature) GetCryptInfoByField(field string) CryptFieldInfo {
	for _, v := range f.CryptFields {
		if v.Field == field {
			return v
		}
	}
	return CryptFieldInfo{}
}

func (f Feature) NeedCrypt() bool {
	if len(f.CryptFields) < 1 {
		return false
	}
	for _, v := range f.CryptFields {
		if v.Enable {
			return true
		}
	}
	return false
}
func (f Feature) NeedSign() bool {
	if len(f.SignFields) < 1 {
		return false
	}
	for _, v := range f.SignFields {
		if v.Enable {
			return true
		}
	}
	return false
}
func (f Feature) RaEnabled() bool {
	return ctype.Bool(f.RaEnable)
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
		f.SignFields = MergeSignFieldsInfo(f.SignFields, ft.SignFields...)
	}
	if len(ft.CryptFields) > 0 {
		f.CryptFields = MergeCryptFieldInfo(f.CryptFields, ft.CryptFields...)
	}
	if len(ft.RaDbFields) > 0 {
		f.RaDbFields = utils.Merge(f.RaDbFields, ft.RaDbFields...)
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
	Enable        bool     `json:"enable"`                                  //可以控制是否开启加密，为了能提前配置加密字段，但是不
	Algo          string   `json:"algo,omitempty" binding:"required"`       // 算法
	SecretKey     []string `json:"secret_key,omitempty" binding:"required"` // 密钥，16字节
	Field         string   `json:"field,omitempty" binding:"required"`      // 需要加密的字段，这个字段需要是 Field 里已定义的
	KeepTxtPreCnt int      `json:"keep_txt_pre_cnt"`                        // 保持文本前几位明文
	KeepTxtSufCnt int      `json:"keep_txt_suf_cnt"`                        // 保持文本后几位明文
	SliceSize     int      `json:"slice_size"`                              // 切片加密的切片大小，为0时不切片。默认0
}

func (s CryptFieldInfo) GetSecretKey() string {
	if len(s.SecretKey) < 1 {
		return defaultCryptKey
	}
	return s.SecretKey[0]
}

func MergeSignFieldsInfo(sfL []SignFieldsInfo, sf ...SignFieldsInfo) []SignFieldsInfo {
	mp := map[string]SignFieldsInfo{}
	for _, v := range sfL {
		mp[v.Field] = v
	}
	for _, v := range sf {
		mp[v.Field] = v
	}
	return maputil.Values(mp)
}
func MergeCryptFieldInfo(cfL []CryptFieldInfo, cf ...CryptFieldInfo) []CryptFieldInfo {

	mp := map[string]CryptFieldInfo{}
	for _, v := range cfL {
		mp[v.Field] = v
	}
	for _, v := range cf {
		mp[v.Field] = v
	}
	return maputil.Values(mp)
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
	return ft
}
