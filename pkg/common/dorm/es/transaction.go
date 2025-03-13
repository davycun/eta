package es

import (
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/broker"
	"gorm.io/gorm"
)

func CommitOrRollback(txDb *gorm.DB, err error) {
	var (
		err1  error
		event *broker.Event
	)
	if !dorm.InTransaction(txDb) {
		return
	}
	ev, b := dorm.LoadAndDelete(txDb, constants.BrokerSync2Es)
	if b {
		event = ev.(*broker.Event)
	}
	if err != nil {
		err1 = txDb.Rollback().Error
		//注意下面的删除或者新增更新要反着操作，因为这里是rollback
		if event != nil {
			switch event.OptType {
			case broker.EventOptTypeDelete:
				err1 = NewApi(global.GetES(), event.TableName).Upsert(event.Data)
			case broker.EventOptTypeInsert, broker.EventOptTypeUpdate:
				err1 = NewApi(global.GetES(), event.TableName).Delete(event.Data)
			}
		}
	} else {
		err1 = txDb.Commit().Error
	}
	if err1 != nil {
		logger.Errorf("DB Rollback Or Commit failed. %v", err)
	}
}
