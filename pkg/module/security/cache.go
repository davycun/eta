package security

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/crypt"
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetTransferCryptKey(c *gin.Context) (string, error) {
	var (
		algoASym   = c.GetHeader(constants.HeaderCryptAsymmetricAlgorithm)
		encryptKey = c.GetHeader(constants.HeaderCryptSymmetryKey)
	)

	//如果没有指定堆成加密算法就表示不需要传输加密
	if algoASym == "" {
		return "", errs.NewClientError(fmt.Sprintf("not found the algorithm in header[%s] ", constants.HeaderCryptAsymmetricAlgorithm))
	}

	if encryptKey == "" {
		return "", errs.NewClientError(fmt.Sprintf("not found the encryptKey in header[%s] ", constants.HeaderCryptSymmetryKey))
	}

	privateKey := crypt.GetPrivateKey(algoASym)
	if privateKey == "" {
		return "", errs.NewClientError(fmt.Sprintf("not support the asymmetric algorithm %s", algoASym))
	}

	//通过非对称密钥解密对称加密（传输加密）的密钥
	cryptKey, err := crypt.DecryptBase64(algoASym, privateKey, encryptKey)
	return string(cryptKey), err
}

func SaveTransferKey(c *ctx.Context, db *gorm.DB, algo string, tkStr string, key string) error {

	tk := TransferKey{}
	tk.Token = ctype.NewStringPrt(tkStr)
	tk.Key = ctype.NewStringPrt(key)
	tk.RequestId = ctype.NewStringPrt(c.GetRequestId())
	tk.Algo = ctype.NewStringPrt(algo)

	_ = entity.BeforeCreate(&tk.BaseEntity, c)
	tkList := []TransferKey{tk}

	return dorm.TableWithContext(db, c, constants.TableTransferKey).Create(&tkList).Error
}
