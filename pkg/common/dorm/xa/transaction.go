package xa

import (
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/es"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/logger"
	"gorm.io/gorm"
)

const (
	SyncEsData = "xa_data_sync2es"
)

type TxData struct {
	TargetData   any
	RollbackData any
	Delete       bool
	EsIndexName  string //es索引名字
}

func CommitOrRollback(txDb *gorm.DB, err error) {
	var (
		err1   error
		txData *TxData
	)
	if !dorm.InTransaction(txDb) {
		return
	}
	dt, b := dorm.LoadAndDelete(txDb, SyncEsData)
	if b {
		txData = dt.(*TxData)
	}
	if err != nil {
		err1 = txDb.Rollback().Error
		//注意下面的删除或者新增更新要反着操作，因为这里是rollback
		if txData != nil {
			if txData.Delete {
				err1 = es.NewApi(global.GetES(), txData.EsIndexName).Upsert(txData.RollbackData)
			} else {
				err1 = es.NewApi(global.GetES(), txData.EsIndexName).Delete(txData.RollbackData)
			}
		}
	} else {
		err1 = txDb.Commit().Error
	}
	if err1 != nil {
		logger.Errorf("DB Rollback Or Commit failed. %v", err)
	}
}
