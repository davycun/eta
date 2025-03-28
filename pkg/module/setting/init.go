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
func GetDefault(category, name string) Setting {
	if s, ok := defaultSettingMap[category+name]; ok {
		return s
	}
	return Setting{}
}

// AddDefaultTableConfig
// 添加默认的表配置初始化到数据库
func AddDefaultTableConfig(cf entity.Table) {
	var (
		st  = defaultSettingMap[ConfigTableCategory+ConfigTableName]
		cfg = &TableConfig{}
	)
	if ctype.IsValid(st.Content) {
		cfg = st.Content.Data.(*TableConfig)
	} else {
		st = Setting{
			Namespace: constants.NamespaceEta,
			Category:  ConfigTableCategory,
			Name:      ConfigTableName,
			Content:   ctype.Json{Data: cfg, Valid: true},
		}
	}
	if cfg.Tables == nil {
		cfg.Tables = make(map[string]entity.Table)
	}
	cfg.Tables[cf.GetTableName()] = cf
	defaultSettingMap[ConfigTableCategory+ConfigTableName] = st
}
func AddDefaultSmsConfig(cf SmsInfo, isDefault bool) {
	if cf.Vendor == "" {
		logger.Errorf("SmsInfo.Vender is empty")
	}
	var (
		st  = defaultSettingMap[ConfigSmsCategory+ConfigSmsName]
		cfg = &SmsConfig{}
	)
	if ctype.IsValid(st.Content) {
		cfg = st.Content.Data.(*SmsConfig)
	} else {
		st = Setting{
			Namespace: constants.NamespaceEta,
			Category:  ConfigSmsCategory,
			Name:      ConfigSmsName,
			Content:   ctype.Json{Data: cfg, Valid: true},
		}
	}
	if cfg.SmsInfoMap == nil {
		cfg.SmsInfoMap = make(map[string]SmsInfo)
	}
	cfg.SmsInfoMap[cf.Vendor] = cf
	if isDefault {
		cfg.Vendor = cf.Vendor
	}
	defaultSettingMap[ConfigSmsCategory+ConfigSmsName] = st
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
