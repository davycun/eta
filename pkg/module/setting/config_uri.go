package setting

import (
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/eta/constants"
	"gorm.io/gorm"
)

// UriConfig
// 忽略配置包括权限的忽略和日志记录的忽略，支持内存中配置（即setting中的Default）和数据库中的配置
// 数据库中的初始化是根据Default配置来的，最终校验的时候是会根据Default的配置和数据库中的配置进行Merge比较的
// 所有的uri的配置规范请参照[github.com/davycun/eta/pkg/common/utils.IsMatchedUri]函数说明针对pattern的说明
type UriConfig struct {
	IgnoreTokenUri     []string `json:"ignore_token_uris,omitempty"`  //哪些请求不需要token校验, 注意这里的格式是METHOD@URI，其中URI支持正则
	IgnoreLogUri       []string `json:"ignore_log_uris,omitempty"`    //哪些请求不需要记录日志，注意这里的格式是METHOD@URI，其中URI支持正则
	IgnoreAuthUri      []string `json:"ignore_auth_uris,omitempty"`   //哪些请求不需要权限校验，主要是菜单对应的api调用权限，注意这里的格式是METHOD@URI，其中URI支持正则
	AdminUri           []string `json:"admin_uri,omitempty"`          //哪些请求是在只有管理员才能调用，注意这里的格式是METHOD@URI，其中URI支持正则
	IgnoreLoadTableUri []string `json:"ignore_entity_uris,omitempty"` // 针对一些请求与非数据库有关的，也就是不用调用[github.com/davydcun/eta/pkg/eta/middleware.LoadTable]
	IgnoreGinLogUri    []string `json:"ignore_gin_log_uri,omitempty"` //忽略那些gin的请求日志
}

func AddDefaultIgnoreTokenUri(cf ...string) {
	cfg := GetDefault[UriConfig](ConfigUriCategory, ConfigUriName)
	cfg.IgnoreTokenUri = utils.Merge(cfg.IgnoreTokenUri, cf...)
	st := Setting{
		Namespace: constants.NamespaceEta,
		Category:  ConfigUriCategory,
		Name:      ConfigUriName,
		Content:   ctype.Json{Data: cfg, Valid: true},
	}
	Registry(st)
}

func AddDefaultIgnoreLogUri(cf ...string) {
	cfg := GetDefault[UriConfig](ConfigUriCategory, ConfigUriName)
	cfg.IgnoreLogUri = utils.Merge(cfg.IgnoreLogUri, cf...)
	st := Setting{
		Namespace: constants.NamespaceEta,
		Category:  ConfigUriCategory,
		Name:      ConfigUriName,
		Content:   ctype.Json{Data: cfg, Valid: true},
	}
	Registry(st)
}
func AddDefaultIgnoreAuthUri(cf ...string) {
	cfg := GetDefault[UriConfig](ConfigUriCategory, ConfigUriName)
	cfg.IgnoreAuthUri = utils.Merge(cfg.IgnoreAuthUri, cf...)
	st := Setting{
		Namespace: constants.NamespaceEta,
		Category:  ConfigUriCategory,
		Name:      ConfigUriName,
		Content:   ctype.Json{Data: cfg, Valid: true},
	}
	Registry(st)
}
func AddDefaultAdminUri(cf ...string) {
	cfg := GetDefault[UriConfig](ConfigUriCategory, ConfigUriName)
	cfg.AdminUri = utils.Merge(cfg.AdminUri, cf...)
	st := Setting{
		Namespace: constants.NamespaceEta,
		Category:  ConfigUriCategory,
		Name:      ConfigUriName,
		Content:   ctype.Json{Data: cfg, Valid: true},
	}
	Registry(st)
}
func AddDefaultIgnoreLoadTableUri(cf ...string) {
	cfg := GetDefault[UriConfig](ConfigUriCategory, ConfigUriName)
	cfg.IgnoreLoadTableUri = utils.Merge(cfg.IgnoreLoadTableUri, cf...)
	st := Setting{
		Namespace: constants.NamespaceEta,
		Category:  ConfigUriCategory,
		Name:      ConfigUriName,
		Content:   ctype.Json{Data: cfg, Valid: true},
	}
	Registry(st)
}
func AddDefaultIgnoreGinLogUri(cf ...string) {
	cfg := GetDefault[UriConfig](ConfigUriCategory, ConfigUriName)
	cfg.IgnoreGinLogUri = utils.Merge(cfg.IgnoreGinLogUri, cf...)
	st := Setting{
		Namespace: constants.NamespaceEta,
		Category:  ConfigUriCategory,
		Name:      ConfigUriName,
		Content:   ctype.Json{Data: cfg, Valid: true},
	}
	Registry(st)
}

func GetUriConfig(db *gorm.DB) UriConfig {
	cfg, err := GetConfig[UriConfig](db, ConfigUriCategory, ConfigUriName)
	if err != nil {
		logger.Errorf("load common config err %s", err)
	}
	dfCfg := GetDefault[UriConfig](ConfigUriCategory, ConfigUriName)

	cfg.IgnoreAuthUri = utils.Merge(cfg.IgnoreAuthUri, dfCfg.IgnoreAuthUri...)
	cfg.IgnoreLogUri = utils.Merge(cfg.IgnoreLogUri, dfCfg.IgnoreLogUri...)
	cfg.IgnoreTokenUri = utils.Merge(cfg.IgnoreTokenUri, dfCfg.IgnoreTokenUri...)
	cfg.AdminUri = utils.Merge(cfg.AdminUri, dfCfg.AdminUri...)
	cfg.IgnoreLoadTableUri = utils.Merge(cfg.IgnoreLoadTableUri, dfCfg.IgnoreLoadTableUri...)
	cfg.IgnoreGinLogUri = utils.Merge(cfg.IgnoreGinLogUri, dfCfg.IgnoreGinLogUri...)
	return cfg
}

func IsIgnoreTokenUri(db *gorm.DB, method, uri string) bool {
	cfg := GetUriConfig(db)
	return utils.IsMatchedUri(method, uri, cfg.IgnoreTokenUri...)
}

func IsIgnoreLogUri(db *gorm.DB, method, uri string) bool {
	cfg := GetUriConfig(db)
	return utils.IsMatchedUri(method, uri, cfg.IgnoreLogUri...)
}

func IsIgnoreAuthUri(db *gorm.DB, method, uri string) bool {
	cfg := GetUriConfig(db)
	return utils.IsMatchedUri(method, uri, cfg.IgnoreAuthUri...)
}

func IsAdminUri(db *gorm.DB, method, uri string) bool {
	cfg := GetUriConfig(db)
	return utils.IsMatchedUri(method, uri, cfg.AdminUri...)
}
func IsIgnoreLoadTableUri(db *gorm.DB, method, uri string) bool {
	cfg := GetUriConfig(db)
	return utils.IsMatchedUri(method, uri, cfg.IgnoreLoadTableUri...)
}
func IsIgnoreGinLogUri(db *gorm.DB, method, uri string) bool {
	cfg := GetUriConfig(db)
	return utils.IsMatchedUri(method, uri, cfg.IgnoreGinLogUri...)
}
