package security_srv

import (
	"github.com/davycun/eta/pkg/common/global"
)

func InitModule() {
	global.GetGin().GET("/security/public_key", publicKey)
	global.GetGin().POST("/security/public_key", publicKey) //为了兼容以前Delta前端
	global.GetGin().GET("/crypto/public_key", publicKey)
	global.GetGin().POST("/crypto/public_key", publicKey)                    //为了兼容以前Delta前端
	global.GetGin().POST("/security/update_transfer_key", updateTransferKey) //更新传输加解密的key
}
