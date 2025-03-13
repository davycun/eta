package validator

import (
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/validity"
)

func AddValidate() {

	validate := global.GetValidator()
	validate.RegisterCustomTypeFunc(validity.ValidateValuer, ctype.LocalTime{}, ctype.Integer{}, ctype.Boolean{})
	err := validate.RegisterValidation(validity.MobileTagName, validity.MobileV, true)
	err = validate.RegisterValidation(validity.IdCardTagName, validity.IdCardV, true)
	err = validate.RegisterValidation(validity.IgnoreTagName, validity.Ignore, true)
	if err != nil {
		logger.Errorf("register validation for %s error", validity.MobileTagName)
	}
}
