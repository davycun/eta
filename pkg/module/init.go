package module

import (
	"github.com/davycun/eta/pkg/module/app/app_srv"
	"github.com/davycun/eta/pkg/module/authorize"
	"github.com/davycun/eta/pkg/module/authorize/role"
	"github.com/davycun/eta/pkg/module/cache"
	"github.com/davycun/eta/pkg/module/dept/dept_srv"
	"github.com/davycun/eta/pkg/module/dict/dict_srv"
	"github.com/davycun/eta/pkg/module/forward/forward_srv"
	"github.com/davycun/eta/pkg/module/integration"
	"github.com/davycun/eta/pkg/module/menu/menu_srv"
	"github.com/davycun/eta/pkg/module/namer/namer_srv"
	"github.com/davycun/eta/pkg/module/optlog"
	"github.com/davycun/eta/pkg/module/reload"
	"github.com/davycun/eta/pkg/module/security/security_srv"
	"github.com/davycun/eta/pkg/module/setting/setting_srv"
	"github.com/davycun/eta/pkg/module/sms/sms_srv"
	"github.com/davycun/eta/pkg/module/storage"
	"github.com/davycun/eta/pkg/module/subscribe"
	"github.com/davycun/eta/pkg/module/template/template_srv"
	"github.com/davycun/eta/pkg/module/user/login/captcha"
	"github.com/davycun/eta/pkg/module/user/login/oauth2"
	"github.com/davycun/eta/pkg/module/user/user2app"
	"github.com/davycun/eta/pkg/module/user/user2role"
	"github.com/davycun/eta/pkg/module/user/user_srv"
)

func InitModules() error {
	app_srv.InitModule()
	authorize.InitModule()
	role.InitModule()
	cache.InitModule()
	template_srv.InitModule()
	dept_srv.InitModule()
	dict_srv.InitModule()
	forward_srv.InitModule()
	integration.InitModule()
	menu_srv.InitModule()
	namer_srv.InitModule()
	optlog.InitModule()
	reload.InitModule()
	security_srv.InitModule()
	setting_srv.InitModule()
	sms_srv.InitModule()
	storage.InitModule()
	subscribe.InitModule()

	captcha.InitModule()
	oauth2.InitModule()
	user2app.InitModule()
	user2role.InitModule()
	user_srv.InitModule()
	return nil
}
