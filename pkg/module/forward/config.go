package forward

import (
	"errors"
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/setting"
	"gorm.io/gorm"
)

const (
	settingCategoryForward = "setting_category_forward"
	settingNameForward     = "setting_name_forward"
)

type Config struct {
	Vendors map[string]setting.BaseCredentials
}

func (c Config) GetVendor(vendor string) (setting.BaseCredentials, error) {
	if c.Vendors == nil {
		return setting.BaseCredentials{}, errors.New(fmt.Sprintf("can not found the vendor %s setting", vendor))
	}
	vd, ok := c.Vendors[vendor]
	if !ok {
		return setting.BaseCredentials{}, errors.New(fmt.Sprintf("can not found the vendor %s setting", vendor))
	}
	return vd, nil
}

func GetVendor(db *gorm.DB, vendor string) (setting.BaseCredentials, error) {

	cfg, err := setting.GetConfig[Config](db, settingCategoryForward, settingNameForward)
	if err != nil {
		return setting.BaseCredentials{}, err
	}
	return cfg.GetVendor(vendor)
}

// AddDefaultVendor
// 提供外部初始化扩展，主要是在程序初始化时调用，把一些默认的配置写入到数据库
func AddDefaultVendor(vendor string, cred setting.BaseCredentials) {
	var (
		set = setting.GetDefault(settingCategoryForward, settingNameForward)
		cfg = Config{}
	)
	if ctype.IsValid(set.Content) {
		switch c := set.Content.Data.(type) {
		case Config:
			cfg = c
		case *Config:
			cfg = *c
		}
	}
	if cfg.Vendors == nil {
		cfg.Vendors = make(map[string]setting.BaseCredentials)
	}

	if _, ok := cfg.Vendors[vendor]; ok {
		logger.Warnf("the target %s for the forwarding has already been set will be overwritten", vendor)
	}
	cfg.Vendors[vendor] = cred

	set.Namespace = constants.NamespaceEta
	set.Category = settingCategoryForward
	set.Name = settingNameForward
	set.Content = ctype.NewJson(&cfg)
	setting.Registry(set)
}
