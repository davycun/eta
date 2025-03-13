package security

import (
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/eta/constants"
	"gorm.io/gorm/schema"
)

// TransferKey
// 传输加密的密钥
type TransferKey struct {
	entity.BaseEntity
	Token     *ctype.String `json:"token,omitempty" redis:"token" gorm:"column:token;comment:token"`
	Algo      *ctype.String `json:"algo,omitempty" redis:"algo" gorm:"column:algo;comment:算法"  binding:"require"`
	Key       *ctype.String `json:"key,omitempty" redis:"key" gorm:"column:key;comment:密钥" binding:"require"`
	RequestId *ctype.String `json:"request_id,omitempty" redis:"request_id" gorm:"column:request_id;comment:requestId"`
}

func (t TransferKey) TableName(namer schema.Namer) string {
	if namer == nil {
		return constants.TableTransferKey
	}
	return namer.TableName(constants.TableTransferKey)
}
