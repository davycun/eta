package security_srv

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/crypt"
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/core/controller"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/security"
	"github.com/davycun/eta/pkg/module/user"
	"github.com/gin-gonic/gin"
)

// 获取指定非对称算法的公钥
func publicKey(c *gin.Context) {
	algo := c.GetHeader(constants.HeaderCryptAsymmetricAlgorithm)
	if algo == "" {
		algo = c.Query("algo")
	}
	if algo == "" {
		//没有传参就用默认的sm2的公钥
		algo = crypt.AlgoASymSm2Pkcs8C132
	}

	if algo != crypt.AlgoASymSm2Pkcs8C132 && algo != crypt.AlgoASymRsaPKCS1v15 {
		controller.ProcessResult(c, nil, errs.NewClientError(fmt.Sprintf("algo only support algorithm %s,%s", crypt.AlgoASymRsaPKCS1v15, crypt.AlgoASymSm2Pkcs8C132)))
		return
	}
	rs := &dto.Result{}
	c.Header(constants.HeaderCryptAsymmetricAlgorithm, algo)
	rs.Data = map[string]string{
		"algo":       algo,
		"public_key": security.GetPublicKey(algo),
	}
	controller.ProcessResult(c, rs, nil)
}

// 更细传输加密的密钥
func updateTransferKey(c *gin.Context) {
	key, err := security.GetTransferCryptKey(c)
	if err != nil {
		controller.ProcessResult(c, nil, err)
		return
	}
	if key == "" {
		controller.ProcessResult(c, nil, errs.NewClientError(fmt.Sprintf("not found the transfer key in header[%s] ", constants.HeaderCryptSymmetryKey)))
		return
	}
	ct := ctx.GetContext(c)
	tkStr := ct.GetContextToken()
	tk, err := user.LoadTokenByToken(tkStr)
	if err != nil {
		controller.ProcessResult(c, nil, err)
	}
	tk.Key = key
	err = user.StoreToken(tk)
	controller.ProcessResult(c, nil, err)
}
