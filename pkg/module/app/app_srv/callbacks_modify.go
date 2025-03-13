package app_srv

import (
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/eta/migrator"
	"github.com/davycun/eta/pkg/module/app"
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
				return afterCreateApp(cfg.Ctx, cfg.TxDB, newValues)
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
		if newV.Database.Host == "" || newV.Database.Port == 0 || newV.Database.User == "" {
			newV.SetDatabase(global.GetLocalDatabase())
			newV.Database.Schema = SchemaPrefix + newV.ID
		}
		if newV.Database.Schema == "" {
			newV.Database.Schema = SchemaPrefix + newV.ID
		}
		if !newV.Valid.Valid { // 没有传这个参数，默认为 true
			newV.Valid = ctype.NewBoolean(true, true)
		}

		//mysql的dbName和Schema保持一致
		if dorm.GetDbType(txDb) == dorm.Mysql {
			newV.Database.DBName = newV.Database.Schema
		}
	}
	//if len(newValues) == 1 && c.GetAppGorm() == nil {
	//	appDb, err := global.LoadGorm(newValues[0].GetDatabase())
	//	if err != nil {
	//		logger.Errorf("beforeCreateApp create db err %s", err)
	//	}
	//	c.SetAppGorm(appDb)
	//}
	return nil
}

// 创建APP之后进行migrate
func afterCreateApp(c *ctx.Context, txDb *gorm.DB, newValues []app.App) error {
	for _, v := range newValues {
		//这里要注意，创建用户的db应该用当前创建app的db，初始化app的db应该用app信息中的database
		appDb, err := global.LoadGorm(v.GetDatabase())
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
