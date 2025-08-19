package reload

import (
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dsync"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/core/service"
)

func (s *Service) ReloadDb2Es(args *dto.Param, rs *dto.Result) error {
	dto.InitPage(args)

	entityCodeList := args.Columns
	args.Columns = []string{}
	if len(entityCodeList) <= 0 {
		return nil
	}

	//srvList, err := service_loader.LoadSrvList(s.GetContext(), true, entityCodeList...)

	srvList := make([]iface.Service, 0, len(entityCodeList))
	for _, v := range entityCodeList {
		sr, err1 := service.NewService(v, s.GetContext().Clone(), s.GetDB())
		if err1 != nil {
			return err1
		}
		srvList = append(srvList, sr)
	}

	args.OrderBy = []dorm.OrderBy{{Column: entity.IdDbName, Asc: true}}
	so := *(args.Extra.(*dsync.SyncOption))

	for _, srv := range srvList {
		logger.Infof("%s 的 db2es reload 操作开始...", srv.GetTableName())
		// 每个表都需要重新构建一下参数。这里的参数，
		args.Extra = &so
		buildNewSyncOption(srv, args)

		sa := &dsync.SyncArgs{Args: args, Srv: srv}
		err := caller.NewCaller().
			Call(func(cl *caller.Caller) error {
				return checkEsIndex(srv)
			}).
			Call(func(cl *caller.Caller) error {
				return BeforeReload(sa)
			}).
			Call(func(cl *caller.Caller) error {
				option := GetSyncOption(srv, args)
				syncSrv := dsync.NewAntsSyncService(DbLoader, EsSaver, *option)
				return syncSrv.Sync(sa, sa)
			}).Err

		if err != nil {
			logger.Errorf("reload db2es Sync 失败. %s", err.Error())
			return err
		}

		err = AfterReload(sa)
		if err != nil {
			logger.Errorf("reload db2es AfterReload 失败. %s", err.Error())
			return err
		}

		logger.Infof("%s 的 db2es reload 操作完成!!!", srv.GetTableName())
	}
	return nil
}
