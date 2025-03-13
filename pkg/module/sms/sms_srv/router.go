package sms_srv

import "github.com/davycun/eta/pkg/common/global"

func Router() {
	global.GetGin().POST("/sms/send_verify_code", SendVerifyCode)
}
