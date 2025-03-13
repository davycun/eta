package user_srv

import (
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/dorm/filter"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/service"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/app"
	"github.com/davycun/eta/pkg/module/user"
	"github.com/davycun/eta/pkg/module/user/user2app"
	"github.com/duke-git/lancet/v2/slice"
)

// 创建root用户和app的关系，相当于默认把root用户加入到所有app中
func modifyCallbackApp(cfg *hook.SrvConfig, pos hook.CallbackPosition) error {
	err := caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			return hook.AfterCreate(cfg, pos, func(cfg *hook.SrvConfig, newValues []app.App) error {
				return afterCreateApp(cfg, newValues)
			})
		}).
		Call(func(cl *caller.Caller) error {
			return hook.AfterDelete(cfg, pos, func(cfg *hook.SrvConfig, oldValues []app.App) error {
				return afterDeleteApp(cfg, oldValues)
			})
		}).Err

	return err
}

// 1.创建root用户和app的关系，相当于默认把root用户加入到所有app中
// 2.前提条件式默认用户已经提前创建，默认用户和默认app都是在migrate的时候创建
func afterCreateApp(cfg *hook.SrvConfig, appList []app.App) error {
	//把admin用户挂到新增的app上，并且设置为管理员，
	var (
		u2aList = make([]user2app.User2App, 0)
	)
	us, err := user.LoadDefaultUser(cfg.TxDB)
	if err != nil || us.ID == "" {
		return err
	}
	//如果是默认app，那么就设置默认用户为管理员
	for _, v := range appList {
		u2a := user2app.User2App{
			BaseEdgeEntity: entity.BaseEdgeEntity{
				FromId: us.ID,
				ToId:   v.ID,
			},
			IsDefault: v.IsDefault,
			IsManager: ctype.NewBooleanPrt(true), //admin用户是所有APP的管理员
		}
		u2aList = append(u2aList, u2a)
	}
	srv := service.NewService(constants.TableUser2App, cfg.Ctx, cfg.TxDB)
	return hook.Create(srv, &dto.Result{}, u2aList)
}

// 1.删除相关联的用户
// 2.登出所有用当前app登录的用户
func afterDeleteApp(cfg *hook.SrvConfig, appList []app.App) error {
	//删除相关联的用户，同时登出所有关联的用户
	var (
		userSvr   = service.NewService(constants.TableUser2App, cfg.Ctx, cfg.TxDB)
		param     = &dto.Param{}
		result    = &dto.Result{}
		appIdList []string
		err       error
		u2aList   []user2app.User2App
	)

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
