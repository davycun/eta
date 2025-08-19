package security_srv

import (
	"github.com/davycun/eta/pkg/common/global"
)

func InitModule() {
	global.GetGin().GET("/security/public_key", publicKey)
	global.GetGin().PUT("/security/update_transfer_key", updateTransferKey)
}
