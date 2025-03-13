package setting

import (
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/storage"
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
