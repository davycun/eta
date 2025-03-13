package userkey

import (
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/history"
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
	//不用传入的namer是为了保障user表都是平台schema下
	return global.GetLocalGorm().NamingStrategy.TableName(constants.TableUserKey)
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

type History struct {
	entity.History
	Entity UserKey `json:"entity,omitempty" gorm:"embedded;embeddedPrefix:h_"`
}

func (h History) TableName(namer schema.Namer) string {
	//不用传入的namer是为了保障user表都是平台schema下
	return global.GetLocalGorm().NamingStrategy.TableName(constants.TableUserHistory)
}

func (h History) AfterMigrator(db *gorm.DB, c *ctx.Context) error {
	scm := dorm.GetDbSchema(db)
	return history.CreateTrigger(db, scm, constants.TableUser)
}
