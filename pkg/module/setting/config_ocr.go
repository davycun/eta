package setting

import (
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/eta/constants"
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

func AddDefaultOcrConfig(cf OcrConfig) {
	defaultSettingMap[ConfigOcrCategory+ConfigOcrName] = Setting{
		Namespace: constants.NamespaceEta,
		Category:  ConfigOcrCategory,
		Name:      ConfigOcrName,
		Content:   ctype.Json{Data: &cf, Valid: true},
	}
}
