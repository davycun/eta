package setting

import (
	"github.com/davycun/eta/pkg/common/logger"
	"gorm.io/gorm"
)

type OcrConfig struct {
	BaseCredentials
	Vendor string `json:"vendor,omitempty" binding:"oneof=aliyun mas ''"`
}

func GetOcrConfig(db *gorm.DB) (OcrConfig, bool) {
	cfg, err := GetConfig[OcrConfig](db, ConfigOcrCategory, ConfigOcrName)
	if err != nil {
		logger.Errorf("load ocr config err %s", err)
		return OcrConfig{}, false
	}
	return cfg, true
}
