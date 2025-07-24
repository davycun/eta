package userkey

import (
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/eta/constants"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type UserKey struct {
	entity.BaseEntity
	UserId       string        `json:"user_id,omitempty" redis:"user_id" gorm:"column:user_id;comment:用户ID;not null"`
	AccessKey    *ctype.String `json:"access_key,omitempty" redis:"access_key" gorm:"type:varchar(255);column:access_key;comment:AccessKey"`
	AccessSecure *ctype.String `json:"access_secure,omitempty" redis:"access_secure" gorm:"type:varchar(255);column:access_secure;comment:AccessSecure"`
	FixedToken   *ctype.String `json:"fixed_token,omitempty" redis:"fixed_token" gorm:"type:varchar(255);column:fixed_token;comment:固定 token"`
}

func (u UserKey) TableName(namer schema.Namer) string {
	if namer == nil {
		return constants.TableUserKey
	}
	return namer.TableName(constants.TableUserKey)
}

func (u UserKey) AfterMigrator(db *gorm.DB, c *ctx.Context) error {
	return caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			return dorm.CreateUniqueIndex(db, constants.TableUserKey, "access_key")
		}).
		Call(func(cl *caller.Caller) error {
			return dorm.CreateUniqueIndex(db, constants.TableUserKey, "fixed_token")
		}).Err
}
