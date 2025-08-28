package reload

import (
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/dsync"
	"github.com/davycun/eta/pkg/common/logger"
)

func db2es(rds *RdService) error {

	logger.Infof("%s 的 db2es reload 操作开始...", rds.GetService().GetTableName())
	// 每个表都需要重新构建一下参数。这里的参数，
	//args.Extra = &so
	buildNewSyncOption(rds.GetService(), rds.GetParam())

	srvArgs := &dsync.SyncArgs{Args: rds.GetParam(), Srv: rds.GetService()}
	err := caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			return checkEsIndex(rds.GetService())
		}).
		Call(func(cl *caller.Caller) error {
			return callbackBefore(RdTypeDb2Es, rds)
		}).
		Call(func(cl *caller.Caller) error {
			option := GetSyncOption(rds.GetService(), rds.GetParam())
			syncSrv := dsync.NewAntsSyncService(DbLoader, EsSaver, *option)
			return syncSrv.Sync(srvArgs, srvArgs)
		}).
		Call(func(cl *caller.Caller) error {
			return callbackAfter(RdTypeDb2Es, rds)
		}).Err

	if err != nil {
		logger.Errorf("reload db2es Sync 失败. %s", err.Error())
	} else {
		logger.Infof("%s 的 db2es reload 操作完成!!!", rds.GetService().GetTableName())
	}
	return err
}
