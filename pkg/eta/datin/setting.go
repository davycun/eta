package datin

import (
	"github.com/davycun/eta/pkg/module/setting"
)

func initSetting() {
	initIgnoreTokenUri()
	initIgnoreLogUri()
	initIgnoreMenuAuthUri()
	initAdminUri()
	initIgnoreLoadTableUri()
	initIgnoreGinLogUri()
}

// uri的配置请参照[github.com/davycun/eta/pkg/utils.IsMatchedUri]的说明
func initIgnoreTokenUri() {
	ignoreList := []string{
		"*@/oauth2/.*",
		"*@/storage/download/.*",
		"*@/storage/upload/.*",
	}
	setting.AddDefaultIgnoreTokenUri(ignoreList...)
}

// uri的配置请参照[github.com/davycun/eta/pkg/utils.IsMatchedUri]的说明
func initIgnoreLogUri() {
	ignoreUri := []string{"*@/authorize/.*", "*@/optlog/.*", "*@/api/.*", "*@/ws/.*", "*@/cache/.*", "*@/forward/.*"}
	setting.AddDefaultIgnoreLogUri(ignoreUri...)
}

// uri的配置请参照[github.com/davycun/eta/pkg/utils.IsMatchedUri]的说明
// 这里的配置主要是在[github.com/davycun/eta/pkg/module/menu/menu_srv.ApiCallAuth]中用到
func initIgnoreMenuAuthUri() {
	baseUri := []string{
		"*@/role/.*", "*@/permission/.*", "*@/auth2role/.*", "*@/cache/.*", "*@/oauth2/.*", "*@/authorize/.*",
		"*@/user/set_current_dept", "*@/user/current", "*@/user/update", "*@/user/id_name", "*@/user/modify_password", "*@/user/reset_password",
		"*@/app/migrate", "*@/menu/.*", "*@/optlog/.*", "*@/setting/.*", "*@/storage/.*", "*@/template/.*", "*@/tasks/.*", "*@/api/.*", "*@/crypto/.*", "*@/ws/.*",
		"*@/citizen/address/.*", "*@/citizen/addr2label/.*", "*@/citizen/address_history/.*", "*@/citizen/address_history/.*", "*@/citizen/bd2label/.*", "*@/citizen/bd2addr/.*",
		"*@/neurond/.*", "*@/tourist_forecast/.*",
	}
	setting.AddDefaultIgnoreAuthUri(baseUri...)
}
func initAdminUri() {
	//
}
func initIgnoreLoadTableUri() {
	uris := []string{
		"*@/forward/.*",
	}
	setting.AddDefaultIgnoreLoadTableUri(uris...)
}
func initIgnoreGinLogUri() {
	uris := []string{
		"*@/forward/.*",
	}
	setting.AddDefaultIgnoreGinLogUri(uris...)
}
