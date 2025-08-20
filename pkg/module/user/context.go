package user

import (
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/app"
	"github.com/davycun/eta/pkg/module/dept"
	"github.com/davycun/eta/pkg/module/user/user2app"
	"gorm.io/gorm"
)

const (
	userContextKey     = "userContextKey"
	curDeptContextKey  = "curDeptContextKey"
	curTokenContextKey = "curTokenContextKey"
)

func init() {
	ctx.AddNewContextOptions(initContext)
}

func initContext(c *ctx.Context) {
	var (
		userId = c.GetContextUserId()
		err    error
		u      User
		ap     app.App
		appDb  *gorm.DB
		dpt    dept.RelationDept
	)
	if userId == "" {
		return
	}

	err = caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			u, err = LoadUserById(global.GetLocalGorm(), userId)
			c.SetContextUserId(userId)
			c.SetContextUserName(u.Name)
			SetContextUser(c, &u)
			if c.GetContextCurrentDeptId() == "" {
				c.SetContextCurrentDeptId(dpt.ID)
			}
			return err
		}).
		Call(func(cl *caller.Caller) error {
			ap, err = user2app.LoadDefaultAppByUserId(global.GetLocalGorm(), u.ID)
			c.SetContextAppId(ap.ID)
			app.SetContextApp(c, &ap)
			c.SetContextIsManager(user2app.UserIsManagerForApp(u.ID, ap.ID))
			return err
		}).
		Call(func(cl *caller.Caller) error {
			appDb, err = global.LoadGormSetAppId(ap.ID, ap.Database)
			if c.GetAppGorm() == nil {
				c.SetAppGorm(appDb)
			}
			if c.GetContextGorm() == nil {
				c.SetContextGorm(appDb)
			}
			return err
		}).
		Call(func(cl *caller.Caller) error {
			c.SetAppDoris(global.LoadAppDoris(ap.Database.Schema, ap.Database.LogLevel, ap.Database.SlowThreshold))
			return nil
		}).
		Call(func(cl *caller.Caller) error {

			u2d, err2 := dept.LoadUser2DeptByUserId(c, userId)
			if err2 != nil {
				return err2
			}

			if len(u2d) < 1 {
				dpt = dept.GetDefaultUser2Dept(userId, u.Name)
			} else {
				for i, v := range u2d {
					if v.IsMain || v.IsManager {
						dpt = u2d[i]
					}
				}
				if dpt.ToId == "" {
					dpt = u2d[0]
				}
			}
			_, b := GetContextDept(c)
			if !b {
				SetContextDept(c, &dpt)
			}
			return nil
		}).Err

	if err != nil {
		logger.Errorf("New Context err %s", err)
	}
}

func GetContextUser(c *ctx.Context) (*User, bool) {
	var (
		u = &User{}
	)
	value, exists := c.Get(userContextKey)
	if exists {
		u = value.(*User)
	}
	return u, exists
}
func SetContextUser(c *ctx.Context, u *User) {
	c.Set(userContextKey, u)
}
func GetContextDept(c *ctx.Context) (*dept.RelationDept, bool) {
	var (
		u = &dept.RelationDept{}
	)
	value, exists := c.Get(curDeptContextKey)

	if exists {
		u = value.(*dept.RelationDept)
	}
	return u, exists
}
func SetContextDept(c *ctx.Context, u *dept.RelationDept) {
	c.Set(curDeptContextKey, u)
}

func SetContextToken(c *ctx.Context, tk TokenInfo) {
	c.SetContextToken(tk.Token)
	c.Set(curTokenContextKey, tk)
}
func GetContextToken(c *ctx.Context) (TokenInfo, bool) {
	value, exists := c.Get(curTokenContextKey)

	if exists {
		return value.(TokenInfo), true
	}
	return TokenInfo{}, false
}

func NewContext(u User, ap app.App, currentDeptId string) (*ctx.Context, error) {
	c := &ctx.Context{}
	c.SetContextAppId(ap.ID)
	c.SetContextUserId(u.ID)
	c.SetContextUserName(u.Name)
	c.SetContextIsManager(user2app.UserIsManagerForApp(u.ID, ap.ID))
	c.SetContextCurrentDeptId(currentDeptId)
	db, err := global.LoadGormSetAppId(ap.ID, ap.GetDatabase())
	if err != nil {
		return c, err
	}

	//DB相关
	c.SetContextGorm(db)
	c.SetAppGorm(db)
	c.SetAppDoris(global.LoadAppDoris(ap.Database.Schema, ap.Database.LogLevel, ap.Database.SlowThreshold))

	//存储用户
	SetContextUser(c, &u)
	app.SetContextApp(c, &ap)
	//协程变量
	ctx.SetCurrentContext(c)
	return c, nil
}

func GetToken(c *ctx.Context) string {
	token := c.GetContextToken()
	if token != "" {
		return token
	}

	token = c.GetGinContext().GetHeader(constants.HeaderAuthorization)
	if token == "" {
		token = c.GetGinContext().Query(constants.HeaderAuthorization)
		if token == "" {
			token = c.GetGinContext().Query("authorization")
		}
	}
	if token == "" {
		token, _ = c.GetGinContext().Cookie(constants.HeaderAuthorization)
		if token == "" {
			token, _ = c.GetGinContext().Cookie("authorization")
		}
	}
	if token != "" {
		c.SetContextToken(token)
	}

	return token
}
