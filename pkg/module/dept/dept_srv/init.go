package dept_srv

import (
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/core/service/sqlbd"
	"github.com/davycun/eta/pkg/eta/constants"
)

func init() {
	hook.AddAuthCallback(constants.TableDept, AuthRetrieve, func(option *hook.CallbackOption) {
		option.Methods = []iface.Method{iface.MethodList}
	})
	sqlbd.AddSqlBuilder(constants.TableDept, buildListSql, iface.MethodList)
}
