package integration

import (
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/dorm/xa"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/run"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/eta/ecf"
	"gorm.io/gorm"
	"sync"
)

type txService struct {
	NewSrv     iface.NewService
	EntityCode string
	C          *ctx.Context
	EC         *iface.EntityConfig
	Param      *dto.Param
	Result     *dto.Result
	Command    string
}

func transactionCall(c *ctx.Context, param *CommandParam, srvList []txService, result *CommandResult) error {
	var (
		err       error
		appTxDb   = xa.Transaction(c.GetAppGorm())
		localTxDb = xa.Transaction(global.GetLocalGorm())
		wg        = sync.WaitGroup{}
	)

	defer func() {
		//TODO 非XA事务，有问题
		xa.CommitOrRollback(appTxDb, err)
		xa.CommitOrRollback(localTxDb, err)
	}()

	for i := range srvList {
		if param.Order {
			err = errs.Cover(err, callTxServices(&srvList[i], localTxDb, appTxDb, result))
			if err != nil {
				return err
			}
		} else {
			wg.Add(1)
			run.Go(func() {
				defer wg.Done()
				if err != nil {
					return
				}
				err = errs.Cover(err, callTxServices(&srvList[i], localTxDb, appTxDb, result))
			})
		}
	}
	wg.Wait()
	return err
}

func callTxServices(txSrv *txService, localTxDb, appTxDb *gorm.DB, result *CommandResult) error {
	var (
		//txSrv       = &srvList[i]
		err         error
		ec          = ecf.GetContextEntityConfig(txSrv.C)
		srvListTemp = make([]iface.Service, 0, 1)
	)
	//如果某个表同时在localDB和appDB
	if ec.LocatedLocal() {
		cx := txSrv.C.Clone()
		cx.SetContextGorm(localTxDb)
		srvListTemp = append(srvListTemp, txSrv.NewSrv(cx, localTxDb, ec))
	}
	if ec.LocatedApp() {
		cx := txSrv.C.Clone()
		cx.SetContextGorm(appTxDb)
		srvListTemp = append(srvListTemp, txSrv.NewSrv(cx, appTxDb, ec))
	}

	for _, srv := range srvListTemp {
		err = callService(srv, txSrv, result)
		if err != nil {
			return err
		}
	}
	return nil
}

func callService(srv iface.Service, txSrv *txService, result *CommandResult) error {

	var (
		err        error
		args       = txSrv.Param
		command    = txSrv.Command
		rs         = txSrv.Result
		entityCode = txSrv.EntityCode
	)

	switch iface.Method(command) {
	case iface.MethodCreate:
		err = srv.Create(args, rs)
	case iface.MethodUpdate:
		err = srv.Update(args, rs)
	case iface.MethodUpdateByFilters:
		err = srv.UpdateByFilters(args, rs)
	case iface.MethodDelete:
		err = srv.Delete(args, rs)
	case iface.MethodDeleteByFilters:
		err = srv.DeleteByFilters(args, rs)
	case iface.MethodQuery:
		err = srv.Query(args, rs)
	case iface.MethodCount:
		err = srv.Count(args, rs)
	case iface.MethodDetail:
		err = srv.Detail(args, rs)
	case iface.MethodPartition:
		err = srv.Partition(args, rs)
	default:
		logger.Errorf("the method[%s] is not support", command)
		return err

	}
	if err == nil {
		result.Items = append(result.Items, CommandResultItem{EntityCode: entityCode, Command: command, Result: rs})
	}

	return err
}
