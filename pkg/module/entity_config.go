package module

import (
	"github.com/davycun/eta/pkg/core/controller"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/core/service"
	"github.com/davycun/eta/pkg/core/service/ecf"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/app"
	"github.com/davycun/eta/pkg/module/app/app_srv"
	"github.com/davycun/eta/pkg/module/authorize/auth"
	"github.com/davycun/eta/pkg/module/authorize/role"
	"github.com/davycun/eta/pkg/module/broker/publish"
	"github.com/davycun/eta/pkg/module/broker/subscribe"
	"github.com/davycun/eta/pkg/module/data/template"
	"github.com/davycun/eta/pkg/module/dept"
	"github.com/davycun/eta/pkg/module/dept/dept_srv"
	"github.com/davycun/eta/pkg/module/dict"
	"github.com/davycun/eta/pkg/module/dict/dict_srv"
	"github.com/davycun/eta/pkg/module/menu"
	"github.com/davycun/eta/pkg/module/menu/menu_srv"
	"github.com/davycun/eta/pkg/module/optlog"
	"github.com/davycun/eta/pkg/module/security"
	"github.com/davycun/eta/pkg/module/setting"
	"github.com/davycun/eta/pkg/module/sms"
	"github.com/davycun/eta/pkg/module/task"
	"github.com/davycun/eta/pkg/module/user"
	"github.com/davycun/eta/pkg/module/user/user2app"
	"github.com/davycun/eta/pkg/module/user/user2dept"
	"github.com/davycun/eta/pkg/module/user/user2role"
	"github.com/davycun/eta/pkg/module/user/user_srv"
	"github.com/davycun/eta/pkg/module/user/userkey"
	"reflect"
)

func RegistryEntityConfig() {
	ecf.Registry(entityConfig()...)
}

func entityConfig() []ecf.EntityConfig {
	NS := constants.NamespaceEta
	return []ecf.EntityConfig{
		//APP模块
		{Namespace: NS, Name: "eta_app", Migrate: true, BaseUrl: "/app", ServiceType: reflect.TypeOf(app_srv.Service{}), Table: entity.Table{EntityType: reflect.TypeOf(app.App{}), Located: entity.LocatedLocal, Order: 1000}},
		{Namespace: NS, Name: "eta_app_history", Migrate: true, DisableApi: true, Table: entity.Table{EntityType: reflect.TypeOf(app.History{}), Located: entity.LocatedLocal}},
		//用户模块
		{Namespace: NS, Name: "eta_user", Migrate: true, BaseUrl: "/user", ServiceType: reflect.TypeOf(user_srv.Service{}), Table: entity.Table{EntityType: reflect.TypeOf(user.User{}), Located: entity.LocatedLocal, Order: 994}},
		{Namespace: NS, Name: "eta_user_history", Migrate: true, DisableApi: true, Table: entity.Table{EntityType: reflect.TypeOf(user.History{}), Located: entity.LocatedLocal, Order: 995}},
		{Namespace: NS, Name: "eta_user_key", Migrate: true, BaseUrl: "/userkey", NewService: service.NewDefaultService, Table: entity.Table{EntityType: reflect.TypeOf(userkey.UserKey{}), Located: entity.LocatedLocal, Order: 996}},
		{Namespace: NS, Name: "eta_user_key_history", Migrate: true, DisableApi: true, Table: entity.Table{EntityType: reflect.TypeOf(userkey.History{}), Located: entity.LocatedLocal}, Order: 997},
		{Namespace: NS, Name: "eta_user2app", Migrate: true, BaseUrl: "/user2app", NewService: service.NewDefaultService, Table: entity.Table{EntityType: reflect.TypeOf(user2app.User2App{}), Located: entity.LocatedLocal, Order: 999}},
		{Namespace: NS, Name: "eta_user2app_history", Migrate: true, DisableApi: true, Table: entity.Table{EntityType: reflect.TypeOf(user2app.History{}), Located: entity.LocatedLocal}, Order: 998},
		{Namespace: NS, Name: "eta_user2dept", Migrate: true, BaseUrl: "/user2dept", NewService: service.NewDefaultService, Table: entity.Table{EntityType: reflect.TypeOf(user2dept.User2Dept{})}},
		{Namespace: NS, Name: "eta_user2role", Migrate: true, BaseUrl: "/user2role", NewService: service.NewDefaultService, Table: entity.Table{EntityType: reflect.TypeOf(user2role.User2Role{})}},
		//权限模块
		{Namespace: NS, Name: "eta_auth2role", Migrate: true, BaseUrl: "/auth2role", NewService: service.NewDefaultService, Table: entity.Table{EntityType: reflect.TypeOf(auth.Auth2Role{})}},
		{Namespace: NS, Name: "eta_permission", Migrate: true, BaseUrl: "/permission", NewService: service.NewDefaultService, Table: entity.Table{EntityType: reflect.TypeOf(auth.Permission{})}},
		{Namespace: NS, Name: "eta_role", Migrate: true, BaseUrl: "/role", NewService: service.NewDefaultService, Table: entity.Table{EntityType: reflect.TypeOf(role.Role{})}},
		//订阅发布
		{Namespace: NS, Name: "eta_publish", Migrate: true, BaseUrl: "/publish", NewService: service.NewDefaultService, Table: entity.Table{EntityType: reflect.TypeOf(publish.Record{})}},
		{Namespace: NS, Name: "eta_subscriber", Migrate: true, BaseUrl: "/subscriber", NewService: service.NewDefaultService, Table: entity.Table{EntityType: reflect.TypeOf(subscribe.Subscriber{})}},
		//数据中心模块
		{Namespace: NS, Name: "eta_template", Migrate: true, BaseUrl: "/template", NewService: service.NewDefaultService, Table: entity.Table{EntityType: reflect.TypeOf(template.Template{})}},
		{Namespace: NS, Name: "eta_template_history", Migrate: true, BaseUrl: "/template_history", NewService: service.NewDefaultService, DisableMethod: iface.GetAllModifyMethod(), Table: entity.Table{EntityType: reflect.TypeOf(template.History{})}},
		{Namespace: NS, Name: "eta_data", Migrate: false, BaseUrl: "/data/:code", NewService: service.NewDefaultService, NewController: controller.NewDefaultController},
		//组织结构模块
		{Namespace: NS, Name: "eta_department", Migrate: true, BaseUrl: "/department", NewService: dept_srv.NewService, Table: entity.Table{EntityType: reflect.TypeOf(dept.Department{})}},
		{Namespace: NS, Name: "eta_department_history", Migrate: true, BaseUrl: "/department_history", NewService: service.NewDefaultService, Table: entity.Table{EntityType: reflect.TypeOf(dept.History{})}},
		//字典模块
		{Namespace: NS, Name: "eta_dictionary", Migrate: true, BaseUrl: "/dictionary", ServiceType: reflect.TypeOf(dict_srv.Service{}), Table: entity.Table{EntityType: reflect.TypeOf(dict.Dictionary{})}},
		//菜单模块
		{Namespace: NS, Name: "eta_menu", Migrate: true, BaseUrl: "/menu", ServiceType: reflect.TypeOf(menu_srv.Service{}), Table: entity.Table{EntityType: reflect.TypeOf(menu.Menu{})}},
		//日志模块
		{Namespace: NS, Name: "eta_optlog", Migrate: true, BaseUrl: "/optlog", NewService: service.NewDefaultService, Table: entity.Table{EntityType: reflect.TypeOf(optlog.OptLog{})}},
		//安全模块之密钥记录
		{Namespace: NS, Name: "eta_token_key", Migrate: true, BaseUrl: "/token_key", DisableApi: true, Table: entity.Table{EntityType: reflect.TypeOf(security.TransferKey{}), Located: entity.LocatedLocal}},
		{Namespace: NS, Name: "eta_token_key_history", Migrate: true, DisableApi: true, Table: entity.Table{EntityType: reflect.TypeOf(security.History{}), Located: entity.LocatedLocal}},
		//配置管理模块
		{Namespace: NS, Name: "eta_setting", Migrate: true, BaseUrl: "/setting", NewService: service.NewDefaultService, Table: entity.Table{EntityType: reflect.TypeOf(setting.Setting{}), Located: entity.LocatedAll}},
		//短信模块
		{Namespace: NS, Name: "eta_sms_task", Migrate: true, BaseUrl: "/sms_task", NewService: service.NewDefaultService, Table: entity.Table{EntityType: reflect.TypeOf(sms.Task{})}},
		{Namespace: NS, Name: "eta_sms_target", Migrate: true, BaseUrl: "/sms_target", NewService: service.NewDefaultService, Table: entity.Table{EntityType: reflect.TypeOf(sms.Target{})}},
		//任务管理模块
		{Namespace: NS, Name: "eta_task", Migrate: true, BaseUrl: "/task", NewService: service.NewDefaultService, DisableMethod: []iface.Method{iface.MethodDeleteByFilters, iface.MethodDelete}, Table: entity.Table{EntityType: reflect.TypeOf(task.DataTask{})}},
	}
}
