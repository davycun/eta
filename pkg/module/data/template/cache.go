package template

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/errs"
	"gorm.io/gorm"
)

func LoadByCode(db *gorm.DB, code string) (temp Template, err error) {
	var (
		tempList []Template
	)
	err = db.Model(&tempList).Where(map[string]any{"code": code}).Find(&tempList).Error
	if err != nil {
		return
	}
	if len(tempList) < 1 {
		return temp, errs.NewClientError(fmt.Sprintf("template[%s] not found", code))
	}
	temp = tempList[0]

	if temp.Status != Ready {
		return Template{}, errs.NewClientError(fmt.Sprintf("the template[%s] is not ready", code))
	}
	return
}
