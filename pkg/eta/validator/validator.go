package validator

import (
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/validity"
)

func InitValidate() error {

	validate := global.GetValidator()
	validate.RegisterCustomTypeFunc(validity.ValidateValuer, ctype.LocalTime{}, ctype.Integer{}, ctype.Boolean{})
	err := caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			return validate.RegisterValidation(validity.MobileTagName, validity.MobileV, true)
		}).
		Call(func(cl *caller.Caller) error {
			return validate.RegisterValidation(validity.IdCardTagName, validity.IdCardV, true)
		}).
		Call(func(cl *caller.Caller) error {
			return validate.RegisterValidation(validity.IgnoreTagName, validity.Ignore, true)
		}).Err
	return err
}
