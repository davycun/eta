package setting

import (
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/storage"
	"github.com/davycun/eta/pkg/eta/constants"
	"gorm.io/gorm"
)

type StorageConfig struct {
	storage.Storage
	PublicFolder string `json:"public_folder,omitempty"` // [仅 default app 有效] 公共目录
}

func GetStorageConfig(db *gorm.DB) (StorageConfig, bool) {
	cfg, err := GetConfig[StorageConfig](db, ConfigStorageCategory, ConfigStorageName)
	if err != nil {
		logger.Errorf("load storage config err %s", err)
		return StorageConfig{}, false
	}
	return cfg, true
}

func AddDefaultStorageConfig(cf StorageConfig) {
	defaultSettingMap[ConfigStorageCategory+ConfigStorageName] = Setting{
		Namespace: constants.NamespaceEta,
		Category:  ConfigStorageCategory,
		Name:      ConfigStorageName,
		Content:   ctype.Json{Data: &cf, Valid: true},
	}
}
