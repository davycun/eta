package user_srv

import (
	"github.com/davycun/eta/pkg/core/controller"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/core/migrate"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/core/service/sqlbd"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/user"
)

func InitModule() {
	//注册回调
	hook.AddModifyCallback(constants.TableUser, modifyCallbackUser)
	hook.AddModifyCallback(constants.TableUser2Dept, modifyCallbackUser2Dept)
	hook.AddModifyCallback(constants.TableDept, modifyCallbackDept)
	hook.AddRetrieveCallback(constants.TableUser, retrieveCallbackUser)

	//注册sql
	sqlbd.AddSqlBuilder(constants.TableUser, buildListSql, iface.MethodList)

	//注册Migrate
	migrate.AddCallback(constants.TableUser, afterMigratorUser)

	//添加路由
	controller.Publish(constants.TableUser, "/list", controller.ApiConfig{
		Handler: func(srv iface.Service, args any, rs any) error {
			s := srv.(*Service)
			return s.Retrieve(args.(*dto.Param), rs.(*dto.Result), iface.MethodList)
		},
		GetParam: func() any {
			return &dto.Param{RetrieveParam: dto.RetrieveParam{Extra: &user.ListParam{}}}
		},
	})
	controller.Publish(constants.TableUser, "/change_passwd", controller.ApiConfig{
		Handler: func(srv iface.Service, args, rs any) error {
			s := srv.(*Service)
			return s.ChangePassword(args.(*user.ModifyPasswordParam), rs.(*dto.Result))
		},
		GetParam: func() any {
			return &user.ModifyPasswordParam{}
		},
	})
	controller.Publish(constants.TableUser, "/set_current_dept", controller.ApiConfig{
		Handler: func(srv iface.Service, args, rs any) error {
			s := srv.(*Service)
			return s.SetCurrentDept(args.(*dto.Param), rs.(*dto.Result))
		},
	})
	controller.Publish(constants.TableUser, "/current", controller.ApiConfig{
		Handler: func(srv iface.Service, args, rs any) error {
			s := srv.(*Service)
			return s.Current(args.(*dto.Param), rs.(*dto.Result))
		},
	})

	controller.Publish(constants.TableUser, "/change_phone", controller.ApiConfig{
		Handler: func(srv iface.Service, args, rs any) error {
			s := srv.(*Service)
			return s.ChangePhone(args.(*user.ModifyPhoneParam), rs.(*dto.Result))
		},
		GetParam: func() any {
			return &user.ModifyPhoneParam{}
		},
	})
	controller.Publish(constants.TableUser, "/reset_password", controller.ApiConfig{
		Handler: func(srv iface.Service, args, rs any) error {
			s := srv.(*Service)
			return s.ResetPassword(args.(*user.ResetPasswordParam), rs.(*dto.Result))
		},
		GetParam: func() any {
			return &user.ResetPasswordParam{}
		},
	})

}
