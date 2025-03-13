package authorize

import (
	"errors"
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/controller"
	"github.com/davycun/eta/pkg/module/app"
	"github.com/davycun/eta/pkg/module/dept"
	"github.com/davycun/eta/pkg/module/third"
	"github.com/davycun/eta/pkg/module/user"
	"github.com/davycun/eta/pkg/module/user/user2app"
	"strings"
)

var (
	appGormUrl = []string{"/user/*", "/app/*", "/security/*"}
)

func Authorize(c *ctx.Context) {
	l := NewAuthorizationService(c).AuthToken().LoadUser().LoadApp().LoadCurrentDept().Store()
	if l.Err != nil {
		controller.Fail(c.GetGinContext(), l.Status, l.Err.Error(), nil)
	}
}

var (
	unLogin = errors.New("用户未登录")
)

type AuthorizeService struct {
	Err     error
	Status  int
	c       *ctx.Context
	Token   user.TokenInfo
	U       user.User
	ap      app.App
	curDept dept.RelationDept
}

func NewAuthorizationService(c *ctx.Context) *AuthorizeService {
	return &AuthorizeService{
		c:      c,
		U:      user.User{},
		ap:     app.App{},
		Status: 500,
	}
}

func (a *AuthorizeService) AuthToken() *AuthorizeService {
	a.Token, a.Err = user.LoadTokenByToken(ctx.GetToken(a.c))
	if a.Token.Token == "" || a.Token.UserId == "" {
		a.Status = 401
		a.Err = unLogin
	} else {
		a.c.SetContextToken(a.Token.Token)
		user.SetContextToken(a.c, a.Token)
	}
	return a
}

func (a *AuthorizeService) LoadUser() *AuthorizeService {
	if a.Err != nil {
		return a
	}
	a.U, a.Err = user.LoadUserById(global.GetLocalGorm(), a.Token.UserId)
	if a.Err != nil {
		a.Status = 401
		a.Err = unLogin
	}
	user.SetContextUser(a.c, &a.U)
	a.c.SetContextUserId(a.U.ID)
	a.c.SetContextUserName(a.U.Name)
	a.c.SetContextIsManager(user2app.UserIsManagerForApp(a.Token.UserId, a.Token.AppId))
	return a
}
func (a *AuthorizeService) LoadApp() *AuthorizeService {
	if a.Err != nil {
		return a
	}
	a.ap, a.Err = app.LoadAppById(global.GetLocalGorm(), a.Token.AppId)
	third.CheckAppVersion(a.ap)

	db, err1 := global.LoadGorm(a.ap.GetDatabase())
	if err1 != nil {
		a.Err = err1
	} else {
		appDoris := global.LoadAppDoris(a.ap.GetDatabase().Schema, a.ap.GetDatabase().LogLevel, a.ap.GetDatabase().SlowThreshold)
		a.c.SetAppDoris(appDoris)
		a.c.SetAppGorm(db)
		a.c.SetContextAppId(a.ap.ID)
		app.SetContextApp(a.c, &a.ap)
	}
	return a
}
func (a *AuthorizeService) LoadCurrentDept() *AuthorizeService {
	if a.Err != nil {
		return a
	}
	var (
		err error
		uri = a.c.GetGinContext().Request.RequestURI
	)

	a.curDept, err = dept.LoadUser2DeptByUserIdDeptId(a.c, a.Token.UserId, a.Token.DeptId)

	//如果不是migrate才处理错误，否则可能会导致部门相关表变动无法migrate的问题
	if !strings.Contains(uri, "app/migrate") {
		if err != nil {
			a.Err = err
			return a
		}
	}
	if a.curDept.Dept.ID != "" {
		a.c.SetContextCurrentDeptId(a.curDept.Dept.ID)
		user.SetContextDept(a.c, &a.curDept)
	}
	return a
}
func (a *AuthorizeService) Store() *AuthorizeService {

	if a.Err != nil {
		return a
	}
	var (
		uri     = a.c.GetGinContext().Request.RequestURI
		db, err = global.LoadGorm(a.ap.GetDatabase())
	)
	if err != nil {
		a.Err = err
		return a
	}

	if (utils.IsMatchedUri(uri, appGormUrl...)) && uri != "/app/migrate" {
		a.c.SetContextGorm(global.GetLocalGorm())
	} else {
		a.c.SetContextGorm(db)
	}
	return a
}
