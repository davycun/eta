package oauth2

import "C"
import (
	"fmt"
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/core/migrate"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/app"
	"github.com/davycun/eta/pkg/module/dept"
	"github.com/davycun/eta/pkg/module/security"
	"github.com/davycun/eta/pkg/module/setting"
	"github.com/davycun/eta/pkg/module/user"
	"github.com/davycun/eta/pkg/module/user/user2app"
	"github.com/davycun/eta/pkg/module/user/user2dept"
	"github.com/davycun/eta/pkg/module/user/user_srv"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
	"time"
)

// ProcessResult
// 主要针对登录过程都通过之后的结果处理
// 1.如果是AccessKey方式，需要考虑生效期内不重新生成token。其余方式每次登录都是生成一个新的token，不影响老的token使用
// 2.
func ProcessResult(c *ctx.Context, us *user.User, loginType string, result *LoginResult) error {

	err := loadOldToken(c, us, loginType, result)
	if err != nil {
		return err
	}
	//如果已经登录过，就直接返回结果，不需要重新生成token，目前是针对AccessToken，其余登录方式都会重新生成新token
	if result.Authorization != "" {
		return nil
	}

	var (
		appId    = c.GetContextAppId() //可能在可能header中传入（此部分只是针对ignoreUri有效）
		loginApp = app.App{}
	)

	err = caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			if us.ID == "" {
				return errs.NewClientError("用户不存在")
			}
			return nil
		}).
		Call(func(cl *caller.Caller) error {
			if appId == "" {
				loginApp, err = user2app.LoadDefaultAppByUserId(global.GetLocalGorm(), us.ID)
				return err
			}
			u2a, err1 := user2app.LoadUser2App(global.GetLocalGorm(), us.ID, appId)
			if err1 != nil {
				return err1
			}
			if u2a.ToId == "" {
				return errs.NewClientError(fmt.Sprintf("用户[%s]没有对应的APP[%s]的权限", ctype.ToString(us.Account), appId))
			}
			loginApp, err = app.LoadAppById(global.GetLocalGorm(), appId)
			if err != nil {
				return err
			}
			if loginApp.ID == "" {
				return errs.NewClientError(fmt.Sprintf("not found the app[%s]", appId))
			}
			return err
		}).
		Call(func(cl *caller.Caller) error {
			if loginApp.ID == "" {
				return errs.NewClientError("用户没有分配任何应用")
			}
			appDb, err1 := global.LoadGorm(loginApp.Database)
			c.SetAppGorm(appDb)
			app.SetContextApp(c, &loginApp)
			return err1
		}).
		Call(func(cl *caller.Caller) error {
			return reGenerateToken(c, us, loginApp, loginType, result)
		}).Err
	return err
}

// 只有access_token 方式不重新生成token
// 其余方式都是只要发生登录行为就会生成一个新的token，表示允许多地登录
func loadOldToken(c *ctx.Context, us *user.User, loginType string, result *LoginResult) error {
	if loginType != constants.LoginTypeAccessToken {
		return nil
	}

	tkList, err := user.LoadTokenByUserId(us.ID)
	if err != nil || len(tkList) < 1 {
		return err
	}

	tk := tkList[0]
	result.Authorization = tk.Token
	result.ExpiresIn = tk.ExpiredAt.Unix() - time.Now().Unix()
	return nil
}
func reGenerateToken(c *ctx.Context, us *user.User, ap app.App, loginType string, result *LoginResult) error {
	var (
		cfg, _     = setting.GetLoginConfig(c.GetAppGorm())
		appDb, err = global.LoadGorm(ap.Database)
		token      = fmt.Sprintf("%s_%s", loginType, ulid.Make().String())
	)

	err = caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			c.SetAppGorm(appDb)
			c.SetContextAppId(ap.ID)
			c.SetContextUserId(ap.ID)
			c.SetContextUserName(us.Name)
			c.SetContextIsManager(user2app.UserIsManagerForApp(us.ID, ap.ID))
			return err
		}).
		Call(func(cl *caller.Caller) error {
			return LoadUserDept(appDb, us, c)
		}).
		Call(func(cl *caller.Caller) error {
			//存储token -> user
			dur := time.Duration(cfg.TokenExpireIn) * time.Second
			tk := user.TokenInfo{
				UserId:    us.ID,
				Token:     token,
				DeptId:    us.CurrentDept.Dept.ID,
				AppId:     ap.ID,
				ExpiredAt: time.Now().Add(dur),
			}
			if loginType == constants.LoginTypeAccessToken {
				tk.OnlyOne = true
			}
			//如果是加密传输，那就同时保存下密钥
			ct := c.GetGinContext()
			if ct != nil {
				tk.Key, _ = security.GetTransferCryptKey(ct)
			}
			if tk.Key != "" {
				algo := ct.GetHeader(constants.HeaderCryptSymmetryAlgorithm)
				_ = security.SaveTransferKey(c, global.GetLocalGorm(), algo, tk.Token, tk.Key)
			}
			return user.StoreToken(tk)
		}).Err

	result.Authorization = token
	result.ExpiresIn = cfg.TokenExpireIn
	result.Data = ctype.Map{
		"user": user_srv.GetUserProp(*us),
		"app":  user_srv.GetAppProp(ap),
	}
	return err
}

func LoadUserDept(appDb *gorm.DB, u *user.User, c *ctx.Context) (err error) {
	u.User2Dept, err = dept.LoadUser2DeptByUserId(c, u.ID)
	if err != nil {
		mig := migrate.NewMigrator(appDb, c)
		err = mig.Migrate(user2dept.User2Dept{}, dept.Department{})
		if err != nil {
			return err
		}
		u.User2Dept, err = dept.LoadUser2DeptByUserId(c, u.ID)
	}

	for i, v := range u.User2Dept {
		if v.IsMain {
			u.CurrentDept = u.User2Dept[i]
		}
	}

	if len(u.User2Dept) == 1 || (u.CurrentDept.Dept.ID == "" && len(u.User2Dept) > 0) {
		u.CurrentDept = u.User2Dept[0]
	}
	c.SetContextCurrentDeptId(u.CurrentDept.Dept.ID)
	return err
}
