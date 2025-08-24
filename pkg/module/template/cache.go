package template

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/core/loader"
	"github.com/davycun/eta/pkg/eta/constants"
	"gorm.io/gorm"
	"strings"
)

var (
	allData = loader.NewCacheLoader[Template, Template](constants.TableTemplate, constants.CacheAllDataTemplate).AddExtraKeyName("code")
)

func SetCacheHasAll(db *gorm.DB, b bool) {
	allData.SetHasAll(db, b)
}

func DelCache(db *gorm.DB, dataList ...Template) {
	for _, v := range dataList {
		if v.ID == "" {
			continue
		}
		allData.Delete(db, v.ID)
	}
}

func LoadByCode(db *gorm.DB, code string) (Template, error) {

	if db == nil {
		return Template{}, nil
	}

	all, err := allData.LoadAll(db)
	if err != nil {
		return Template{}, err
	}
	code = strings.TrimLeft(code, constants.TableTemplatePrefix)
	for _, v := range all {
		if v.Code == code {
			if v.Status != Ready {
				return Template{}, errs.NewClientError(fmt.Sprintf("the template[%s] is not ready", code))
			}
			return v, nil
		}
	}
	return Template{}, errs.NewClientError(fmt.Sprintf("template[%s] not found", code))
}
