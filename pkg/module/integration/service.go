package integration

import (
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/es"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/iface"
)

type txService struct {
	NewSrv     iface.NewService
	EntityCode string
	C          *ctx.Context
	TB         *entity.Table
	Param      *dto.Param
	Result     *dto.Result
	Command    string
}

func TransactionCall(c *ctx.Context, srvList []txService, result *CommandResult) error {
	var (
		err       error
		appTxDb   = dorm.Transaction(c.GetAppGorm())
		localTxDb = dorm.Transaction(global.GetLocalGorm())
	)

	defer func() {
		//TODO 非XA事务，有问题
		es.CommitOrRollback(appTxDb, err)
		es.CommitOrRollback(localTxDb, err)
	}()

	for _, txSrv := range srvList {
		var (
			srv iface.Service
			tb  = entity.GetContextTable(txSrv.C)
			rs  = CommandResultItem{EntityCode: txSrv.EntityCode, Command: txSrv.Command, Result: txSrv.Result}
		)
		//如果某个表同时在localDB和appDB，那么暴露的服务只操作appDB的表
		if tb.LocatedApp() {
			srv = txSrv.NewSrv(txSrv.C, appTxDb, tb)
		} else {
			srv = txSrv.NewSrv(txSrv.C, localTxDb, tb)
		}
		//if tb.LocalDB {
		//	srv = txSrv.NewSrv(txSrv.C, localTxDb, tb)
		//} else {
		//	srv = txSrv.NewSrv(txSrv.C, appTxDb, tb)
		//}

		switch iface.Method(txSrv.Command) {
		case iface.MethodCreate:
			err = srv.Create(txSrv.Param, txSrv.Result)
		case iface.MethodUpdate:
			err = srv.Update(txSrv.Param, txSrv.Result)
		case iface.MethodUpdateByFilters:
			err = srv.UpdateByFilters(txSrv.Param, txSrv.Result)
		case iface.MethodDelete:
			err = srv.Delete(txSrv.Param, txSrv.Result)
		case iface.MethodDeleteByFilters:
			err = srv.DeleteByFilters(txSrv.Param, txSrv.Result)
		}
		if err != nil {
			return err
		}
		result.Items = append(result.Items, rs)
	}
	return nil
}
