package setting

import (
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/eta/constants"
)

var (
	defaultSettingMap = make(map[string]Setting) //category+name -> Setting
)

// Registry 添加默认的配置，在setting表初次创建的时候会把内容写入数据库
func Registry(settingList ...Setting) {
	for _, v := range settingList {
		defaultSettingMap[v.Category+v.Name] = v
	}
}

// GetDefault
// 可以提供默认值的获取后修改在Registry
func GetDefault[T any](category, name string) T {
	var (
		t     T
		s, ok = defaultSettingMap[category+name]
	)
	if !ok {
		return t
	}
	if ctype.IsValid(s.Content) {
		switch v := s.Content.Data.(type) {
		case T:
			t = v
		case *T:
			t = *v
		}
	}
	return t
}

// AddDefaultTableConfig
// 添加默认的表配置初始化到数据库
func AddDefaultTableConfig(cf entity.Table) {
	var (
		cfg = GetDefault[TableConfig](ConfigTableCategory, ConfigTableName)
	)
	if cfg.Tables == nil {
		cfg.Tables = make(map[string]entity.Table)
	}
	cfg.Tables[cf.GetTableName()] = cf
	st := Setting{
		Namespace: constants.NamespaceEta,
		Category:  ConfigTableCategory,
		Name:      ConfigTableName,
		Content:   ctype.Json{Data: cfg, Valid: true},
	}
	Registry(st)
}
func AddDefaultSmsConfig(cf SmsInfo, isDefault bool) {
	if cf.Vendor == "" {
		logger.Errorf("SmsInfo.Vender is empty")
	}
	var (
		cfg = GetDefault[SmsConfig](ConfigSmsCategory, ConfigSmsName)
	)
	if cfg.SmsInfoMap == nil {
		cfg.SmsInfoMap = make(map[string]SmsInfo)
	}
	cfg.SmsInfoMap[cf.Vendor] = cf
	if isDefault {
		cfg.Vendor = cf.Vendor
	}
	st := Setting{
		Namespace: constants.NamespaceEta,
		Category:  ConfigSmsCategory,
		Name:      ConfigSmsName,
		Content:   ctype.Json{Data: cfg, Valid: true},
	}
	Registry(st)
}
func AddDefaultLoginConfig(cf LoginConfig) {
	defaultSettingMap[ConfigLoginCategory+ConfigLoginName] = Setting{
		Namespace: constants.NamespaceEta,
		Category:  ConfigLoginCategory,
		Name:      ConfigLoginName,
		Content:   ctype.Json{Data: &cf, Valid: true},
	}
}
func AddDefaultCommonConfig(cf LoginConfig) {
	defaultSettingMap[ConfigCommonCategory+ConfigCommonName] = Setting{
		Namespace: constants.NamespaceEta,
		Category:  ConfigCommonCategory,
		Name:      ConfigCommonName,
		Content:   ctype.Json{Data: &cf, Valid: true},
	}
}
func AddDefaultStorageConfig(cf StorageConfig) {
	defaultSettingMap[ConfigStorageCategory+ConfigStorageName] = Setting{
		Namespace: constants.NamespaceEta,
		Category:  ConfigStorageCategory,
		Name:      ConfigStorageName,
		Content:   ctype.Json{Data: &cf, Valid: true},
	}
}
func AddDefaultOcrConfig(cf OcrConfig) {
	defaultSettingMap[ConfigOcrCategory+ConfigOcrName] = Setting{
		Namespace: constants.NamespaceEta,
		Category:  ConfigOcrCategory,
		Name:      ConfigOcrName,
		Content:   ctype.Json{Data: &cf, Valid: true},
	}
}
