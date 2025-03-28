package module

import (
	"github.com/davycun/eta/pkg/module/app/app_srv"
	"github.com/davycun/eta/pkg/module/authorize"
	"github.com/davycun/eta/pkg/module/cache"
	"github.com/davycun/eta/pkg/module/data/template/template_srv"
	"github.com/davycun/eta/pkg/module/dept/dept_srv"
	"github.com/davycun/eta/pkg/module/dict/dict_srv"
	"github.com/davycun/eta/pkg/module/forward/forward_srv"
	"github.com/davycun/eta/pkg/module/integration"
	"github.com/davycun/eta/pkg/module/menu/menu_srv"
	"github.com/davycun/eta/pkg/module/namer/namer_srv"
	"github.com/davycun/eta/pkg/module/plugin"
	"github.com/davycun/eta/pkg/module/security/security_srv"
	"github.com/davycun/eta/pkg/module/storage"
	"github.com/davycun/eta/pkg/module/user/login/captcha"
	"github.com/davycun/eta/pkg/module/user/login/oauth2"
	"github.com/davycun/eta/pkg/module/user/user_srv"
	"github.com/davycun/eta/pkg/module/ws_api"
)

func Router() {
	app_srv.Router()
	authorize.Router()
	cache.Router()
	template_srv.Router()
	dept_srv.Router()
	dict_srv.Router()
	integration.Router()
	menu_srv.Router()
	namer_srv.Router()
	security_srv.Router()
	storage.Router()
	user_srv.Router()
	captcha.Router()
	oauth2.Router()
	///单独初始化一下
	plugin.InitPlugin()
	ws_api.Router()
	//转发模块
	forward_srv.Router()
}
