package sms_srv

import (
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/eta/constants"
)

func InitModule() {

	hook.AddModifyCallback(constants.TableSmsTask, modifyCallbacks)
	global.GetGin().POST("/sms/send_verify_code", SendVerifyCode)
}
