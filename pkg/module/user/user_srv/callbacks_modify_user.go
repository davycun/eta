package user_srv

import (
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/app"
	"github.com/davycun/eta/pkg/module/namer"
	"github.com/davycun/eta/pkg/module/setting"
	"github.com/davycun/eta/pkg/module/user"
	"github.com/davycun/eta/pkg/module/user/user2app"
	"github.com/duke-git/lancet/v2/slice"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"regexp"
	"time"
)

// 服务层更新用户信息
func modifyCallbackUser(cfg *hook.SrvConfig, pos hook.CallbackPosition) error {

	err := caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			return hook.BeforeCreate(cfg, pos, func(cfg *hook.SrvConfig, newValues []user.User) error {
				return beforeCreateFillUserField(cfg.Ctx, cfg.TxDB, newValues)
			})
		}).
		Call(func(cl *caller.Caller) error {
			return hook.AfterCreate(cfg, pos, func(cfg *hook.SrvConfig, newValues []user.User) error {
				namer.DelIdNameCacheByContext(cfg.Ctx)
				return afterCreateUserNewUser2App(cfg.Ctx, cfg.TxDB, newValues)
			})
		}).
		Call(func(cl *caller.Caller) error {
			return hook.BeforeUpdate(cfg, pos, filterUpdateColumn)
		}).
		Call(func(cl *caller.Caller) error {
			return hook.AfterUpdate(cfg, pos, func(cfg *hook.SrvConfig, oldValues []user.User, newValues []user.User) error {
				namer.DelIdNameCacheByContext(cfg.Ctx)
				return delUserCache(cfg.NewValues)
			})
		}).
		Call(func(cl *caller.Caller) error {
			return hook.AfterDelete(cfg, pos, func(cfg *hook.SrvConfig, oldValues []user.User) error {
				namer.DelIdNameCacheByContext(cfg.Ctx)
				return delUserCache(cfg.OldValues)
			})
		}).Err

	return err
}

// db 可以是appDB也可以是localDB，如果是nil，那么默认是localDB
func beforeCreateFillUserField(c *ctx.Context, db *gorm.DB, usList []user.User) error {
	for i, _ := range usList {
		err := entity.BeforeCreate(&usList[i].BaseEntity, c)
		if err != nil {
			return err
		}
		us := &usList[i]
		us.Valid = ctype.Boolean{Valid: true, Data: true}
		//如果在service层没有初始化完这个ID，那么就把自己的ID设置给自己，其实理论上不太合理
		if us.Category == "" {
			us.Category = constants.UserTypeSystem
		}

		if us.Password == "" {
			us.Password = ctype.ToString(us.Account) + "@Abc123"
		}
		// 校验密码
		cfg, _ := setting.GetLoginConfig(db)
		if cfg.PwdValidateReg != "" {
			if match, _ := regexp.MatchString(cfg.PwdValidateReg, us.Password); !match {
				err = errs.NewClientError("密码格式不正确")
				if err != nil {
					return err
				}
			}
		}
		us.LastUpdatePwd = ctype.NewLocalTimePrt(time.Now())

		passwd, err1 := bcrypt.GenerateFromPassword(utils.StringToBytes(us.Password), bcrypt.DefaultCost)
		if err1 != nil {
			return err1
		}
		us.Password = utils.BytesToString(passwd)
	}
	return nil
}

func filterUpdateColumn(cfg *hook.SrvConfig, oldValues []user.User, newValues []user.User) error {
	var (
		pass = []string{"password", "last_update_pwd"}
	)
	cfg.Param.Columns = slice.Filter(cfg.Param.Columns, func(index int, item string) bool {
		return !slice.Contain(pass, item)
	})

	slice.ForEach(oldValues, func(index int, item user.User) {
		oldValues[index].Password = ""
		oldValues[index].LastUpdatePwd = nil
	})
	return nil
}

func delUserCache(data any) error {
	var us []user.User
	utils.ConvertToSlice(data, &us)
	slice.ForEach(us, func(index int, item user.User) {
		user.DelUserCache(item.ID)
	})

	return nil
}

// 主要是migrate和admin管理员创建
func afterCreateUserNewUser2App(c *ctx.Context, txDb *gorm.DB, usList []user.User) error {
	//如果context存在APP，就创建于当前APP的关系，如果没有，那么就创建于默认app的关系
	if c == nil {
		c = ctx.NewContext()
	}
	if txDb == nil {
		txDb = global.GetLocalGorm()
	}
	var (
		appId   = c.GetContextAppId()
		u2aList = make([]user2app.User2App, 0, 1)
	)

	if appId == "" {
		ap, err := app.LoadDefaultApp(txDb)
		if err != nil {
			return err
		}
		if ap.ID == "" {
			return errs.NewServerError("not found the default app")
		}
		appId = ap.ID
	}

	for _, v := range usList {
		u2aList = append(u2aList, user2app.User2App{
			BaseEdgeEntity: entity.BaseEdgeEntity{FromId: v.ID, ToId: appId},
		})
	}

	cls := clause.OnConflict{
		Columns: []clause.Column{
			{Name: entity.FromIdDbName},
			{Name: entity.ToIdDbName},
		},
		DoNothing: true,
	}

	return dorm.Table(txDb, constants.TableUser2App).Clauses(cls).Create(&u2aList).Error
}
