package app_srv

import (
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/dorm/filter"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/service"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/eta/migrator"
	"github.com/davycun/eta/pkg/module/app"
	"github.com/davycun/eta/pkg/module/user"
	"github.com/davycun/eta/pkg/module/user/user2app"
	"github.com/duke-git/lancet/v2/slice"
	"gorm.io/gorm"
)

const (
	SchemaPrefix = "eta_"
)

func modifyCallback(cfg *hook.SrvConfig, pos hook.CallbackPosition) error {
	var (
		err error
	)

	caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			return hook.BeforeCreate(cfg, pos, func(cfg *hook.SrvConfig, newValues []app.App) error {
				return beforeCreateApp(cfg.Ctx, cfg.TxDB, newValues)
			})
		}).
		Call(func(cl *caller.Caller) error {
			return hook.AfterCreate(cfg, pos, func(cfg *hook.SrvConfig, newValues []app.App) error {
				return afterCreateAppMigrateApp(cfg.Ctx, newValues)
			})
		}).
		Call(func(cl *caller.Caller) error {
			return hook.AfterCreate(cfg, pos, func(cfg *hook.SrvConfig, newValues []app.App) error {
				return afterCreateAppCreateUser2App(cfg, newValues)
			})
		}).
		Call(func(cl *caller.Caller) error {
			return hook.AfterDelete(cfg, pos, func(cfg *hook.SrvConfig, oldValues []app.App) error {
				return afterDeleteAppDeleteUser2App(cfg, oldValues)
			})
		}).
		Call(func(cl *caller.Caller) error {
			return hook.AfterDelete(cfg, pos, func(cfg *hook.SrvConfig, oldValues []app.App) error {
				return delAppCache(oldValues)
			})
		}).
		Call(func(cl *caller.Caller) error {
			return hook.AfterUpdate(cfg, pos, func(cfg *hook.SrvConfig, oldValues []app.App, newValues []app.App) error {
				return delAppCache(oldValues)
			})
		})
	return err
}

func beforeCreateApp(c *ctx.Context, txDb *gorm.DB, newValues []app.App) error {
	if c == nil {
		c = ctx.NewContext()
	}
	for i, _ := range newValues {
		newV := &newValues[i]
		err := entity.BeforeCreate(&(newV.BaseEntity), c)
		if err != nil {
			return err
		}
		scmPrefix := global.GetConfig().Database.SchemaPrefix
		if scmPrefix == "" {
			scmPrefix = global.GetConfig().Database.Schema + "_"
		}
		if scmPrefix == "" {
			scmPrefix = SchemaPrefix
		}
		if newV.Database.Host == "" || newV.Database.Port == 0 || newV.Database.User == "" {
			newV.SetDatabase(global.GetLocalDatabase())
			newV.Database.Schema = scmPrefix + newV.ID
		}
		if newV.Database.Schema == "" {
			newV.Database.Schema = scmPrefix + newV.ID
		}
		if !newV.Valid.Valid { // 没有传这个参数，默认为 true
			newV.Valid = ctype.NewBoolean(true, true)
		}

		//mysql的dbName和Schema保持一致
		if dorm.GetDbType(txDb) == dorm.Mysql {
			newV.Database.DBName = newV.Database.Schema
		}
	}
	return nil
}

// 创建APP之后进行migrate
func afterCreateAppMigrateApp(c *ctx.Context, newValues []app.App) error {
	for _, v := range newValues {
		//这里要注意，创建用户的db应该用当前创建app的db，初始化app的db应该用app信息中的database
		appDb, err := global.LoadGormSetAppId(v.ID, v.GetDatabase())
		if err != nil {
			return err
		}
		c1 := c.Clone()
		c1.SetAppGorm(appDb)
		err = migrator.MigrateApp(appDb, c1, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

// 1.创建root用户和app的关系，相当于默认把root用户加入到所有app中
// 2.前提条件式默认用户已经提前创建，默认用户和默认app都是在migrate的时候创建
func afterCreateAppCreateUser2App(cfg *hook.SrvConfig, appList []app.App) error {
	us, err := user.LoadDefaultUser(cfg.TxDB)
	//如果是初次初始化，那么优先Migrate APP的时候，这里可能为空，那么默认app和默认用户的关系就由Migrate User来处理
	//如果先MigrateUser，那么这里会执行，MigrateUserAfter中的CreateUser2App针对root用户不会执行
	if err != nil || us.ID == "" {
		return err
	}
	u2aList := make([]user2app.User2App, 0, 1)
	//如果是默认app，那么就设置默认用户为管理员
	for _, v := range appList {
		u2a := user2app.User2App{
			BaseEdgeEntity: entity.BaseEdgeEntity{
				FromId: us.ID,
				ToId:   v.ID,
			},
			IsManager: ctype.NewBooleanPrt(true), //admin用户是所有APP的管理员
		}
		u2aList = append(u2aList, u2a)
	}
	return service.NewSrvWrapper(constants.TableUser2App, cfg.Ctx, cfg.TxDB).SetData(u2aList).Create()
}

// 1.删除相关联的用户
// 2.登出所有用当前app登录的用户
func afterDeleteAppDeleteUser2App(cfg *hook.SrvConfig, appList []app.App) error {
	//删除相关联的用户，同时登出所有关联的用户
	var (
		userSvr, err = service.NewService(constants.TableUser2App, cfg.Ctx, cfg.TxDB)
		param        = &dto.Param{}
		result       = &dto.Result{}
		appIdList    []string
		u2aList      []user2app.User2App
	)
	if err != nil {
		return err
	}

	err = caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			appIdList = slice.Map(appList, func(index int, item app.App) string {
				return item.ID
			})
			return nil
		}).
		Call(func(cl *caller.Caller) error {
			u2aList, err = user2app.LoadUser2AppByAppId(cfg.TxDB, appIdList...)
			return err
		}).
		Call(func(cl *caller.Caller) error {
			//删除用户与对应app的所有关系
			param.Filters = []filter.Filter{
				{
					LogicalOperator: filter.And,
					Column:          entity.ToIdDbName,
					Operator:        filter.IN,
					Value:           appIdList,
				},
			}
			return userSvr.DeleteByFilters(param, result)
		}).
		Call(func(cl *caller.Caller) error {
			//登出所有相关的用户
			for _, v := range u2aList {
				err = user.LogOutUser(v.FromId, v.ToId)
				if err != nil {
					return err
				}
			}
			return err
		}).Err
	return err
}

func delAppCache(apps []app.App) error {
	//TODO 这里需要思考，如果把App的Database信息更新为一个新的DB（相当于抛弃了之前的DB，类似创建了一个新DB）
	//TODO 首先是否允许这样的事情发生，如果允许，那么需要进行migrateApp
	slice.ForEach(apps, func(index int, item app.App) {
		app.DelAppCache(item.ID)
		global.DeleteGorm(item.Database)
		dorisCfg := global.GetConfig().Doris
		dorisCfg.Schema = item.Database.Schema
		global.DeleteGorm(dorisCfg)
	})
	return nil
}
