package user

import (
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/dept"
	"github.com/davycun/eta/pkg/module/user/user2app"
	"github.com/davycun/eta/pkg/module/user/userkey"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

const (
	SchemaPrefix = "eta_"
	TableName    = constants.TableUser
)

var (
	// NotVirtualUser 真实用户（非虚拟用户）
	NotVirtualUser     = []string{constants.UserTypeSystem, constants.UserTypeDatlas}
	DefaultUserColumns = append(entity.DefaultVertexColumns, "account", "name", "phone", "category", "email", "sex", "valid", "job_title", "is_manager")
)

// User
// TODO 支持多个第三方用户唯一ID的方式，可以进行改进，不然不断地扩充User表的字段不合适
// TODO 单独一张表记录，并且添加索引
// TODO 新增json字段统一存放，不一定合适，查询性能可能受限，管理也不方便
type User struct {
	entity.BaseEntity
	Account       *ctype.String    `json:"account,omitempty" redis:"account" gorm:"column:account;comment:登录账号;not null" binding:"required"`
	Phone         *ctype.String    `json:"phone,omitempty"  gorm:"column:phone;comment:手机号码;" redis:"phone"`
	Password      string           `json:"password,omitempty" redis:"password" gorm:"column:password;comment:密码" binding:"required"`
	LastUpdatePwd *ctype.LocalTime `json:"last_update_pwd,omitempty" gorm:"type:timestamp with time zone;comment:最近更新密码时间" redis:"last_update_pwd"`
	Category      string           `json:"category,omitempty" gorm:"column:category;comment:用户分类"` //用户分类，比如系统用户创建这类用户主要是导入数据用或者为了来源而用
	Name          string           `json:"name,omitempty"  gorm:"column:name;comment:用户姓名" binding:"required" redis:"name"`
	Email         *ctype.String    `json:"email,omitempty"  gorm:"column:email;comment:邮箱" redis:"email"`
	Avatar        string           `json:"avatar,omitempty"  gorm:"column:avatar;comment:头像" redis:"avatar"`
	Sex           string           `json:"sex,omitempty"  gorm:"column:sex;comment:性别" redis:"sex"`
	Valid         ctype.Boolean    `json:"valid,omitempty"  gorm:"column:valid;comment:启用或禁用" redis:"valid"`
	JobTitle      string           `json:"job_title,omitempty"  gorm:"column:job_title;comment:职位" redis:"job_title"`
	WechatUnionId *ctype.String    `json:"wechat_union_id,omitempty" redis:"wechat_union_id" gorm:"type:varchar(255);column:wechat_union_id;comment:微信UnionId"`
	DingUnionId   *ctype.String    `json:"ding_union_id,omitempty" redis:"ding_union_id" gorm:"type:varchar(255);column:ding_union_id;comment:钉钉UnionId"`
	ZzdAccountId  *ctype.String    `json:"zzd_account_id,omitempty" redis:"zzd_account_id" gorm:"type:varchar(255);column:zzd_account_id;comment:浙政钉AccountId"`
	WecomUserId   *ctype.String    `json:"wecom_user_id,omitempty" redis:"wecom_user_id" gorm:"type:varchar(255);column:wecom_user_id;comment:企业微信UserId"`
	DatlasUserId  *ctype.String    `json:"datlas_user_id,omitempty" redis:"datlas_user_id" gorm:"type:varchar(255);column:datlas_user_id;comment:Datlas UserId"`
	YthUserId     *ctype.String    `json:"yth_user_id,omitempty" redis:"yth_user_id" gorm:"type:varchar(255);column:yth_user_id;comment:一体化 UserId"`
	Sign          ctype.String     `json:"sign,omitempty" gorm:"sign;comment:签名值"` // 签名值

	User2App    []user2app.User2App `json:"user2app,omitempty" gorm:"-:all"`
	UserKey     []userkey.UserKey   `json:"user_key,omitempty" gorm:"-:all"`
	SignMatched ctype.Boolean       `json:"sign_matched,omitempty" gorm:"-:all"` // 签名匹配字段,数据完整性校验
	User2Dept   []dept.RelationDept `json:"user2dept,omitempty" gorm:"-:all" binding:"ignore"`
	CurrentDept dept.RelationDept   `json:"current_dept,omitempty" gorm:"-:all"  binding:"ignore"`
}

func (u User) TableName(namer schema.Namer) string {
	if namer == nil {
		return constants.TableUser
	}
	return namer.TableName(constants.TableUser)
}

func (u User) AfterMigrator(db *gorm.DB, c *ctx.Context) error {

	return caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			return dorm.CreateUniqueIndex(db, TableName, "account")
		}).
		Call(func(cl *caller.Caller) error {
			return dorm.CreateUniqueIndex(db, TableName, "phone")
		}).
		Call(func(cl *caller.Caller) error {
			return dorm.CreateUniqueIndex(db, TableName, "email")
		}).Err
}
