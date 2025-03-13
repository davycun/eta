package oauth2

import (
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/module/user"
	"github.com/duke-git/lancet/v2/slice"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func LoginByAccount(c *ctx.Context, args any) (user.User, error) {

	var (
		err      error
		userList = make([]user.User, 0, 1)
		param    = args.(*LoginByUsernameParam)
	)

	err = caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			//是否锁定校验
			return LoginFailLockCheck(c.GetAppGorm(), param.Username)
		}).
		Call(func(cl *caller.Caller) error {
			return global.GetLocalGorm().Model(&userList).Where(map[string]any{"account": param.Username, "valid": true}).Find(&userList).Error
		}).
		//密码校验
		Call(func(cl *caller.Caller) error {
			if len(userList) < 1 {
				return errs.NewClientError("用户不存在或被禁用")
			}
			userList = slice.Filter(userList, func(index int, u user.User) bool {
				start := time.Now().UnixMilli()
				err = bcrypt.CompareHashAndPassword(utils.StringToBytes(u.Password), utils.StringToBytes(param.Password))
				logger.Infof("bcrypt.CompareHashAndPassword()耗时: %dms", time.Now().UnixMilli()-start)
				return err == nil
			})
			if len(userList) < 1 {
				err1 := LoginFailLockCounterIncr(c.GetAppGorm(), param.Username)
				if err1 != nil {
					logger.Errorf("Login fail lock err %s", err1)
				}
				return errs.NewClientError("用户密码错误")
			} else if len(userList) > 1 {
				return errs.NewServerError("账号重复，请联系管理员")
			}
			return nil
		}).Err

	if err != nil {
		return user.User{}, err
	}

	// 登录成功，清除失败计数
	err = LoginFailLockCounterClear(c.GetAppGorm(), param.Username)
	if err != nil {
		logger.Errorf("登录失败计数器异常, %v", err)
	}
	return userList[0], nil
}
