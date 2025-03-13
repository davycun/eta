package plugin_es

import (
	"errors"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/dorm/es"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/broker"
	"gorm.io/gorm"
	"reflect"
)

// ModifyCallbackForEs
// modify callback for ES, 同步到 ES
func ModifyCallbackForEs(cfg *hook.SrvConfig, pos hook.CallbackPosition) error {
	switch pos {
	case hook.CallbackAfter:
		switch cfg.Method {
		case iface.MethodCreate, iface.MethodUpdate, iface.MethodUpdateByFilters:
			return Sync2Es(cfg.TxDB, cfg.GetTable(), cfg.NewValues, false, false)
		case iface.MethodDelete, iface.MethodDeleteByFilters:
			return Sync2Es(cfg.TxDB, cfg.GetTable(), cfg.OldValues, true, false)
		}
	}
	return nil
}

func Sync2Es(txDb *gorm.DB, tb *entity.Table, entityList any, delete bool, autoCommit bool) error {
	if entityList == nil {
		return nil
	}
	el := reflect.ValueOf(entityList)
	if el.Kind() == reflect.Ptr {
		el = el.Elem()
	}
	if el.Kind() != reflect.Slice {
		return errors.New("entityList must be a slice")
	}
	if el.Len() <= 0 || global.GetES() == nil || !ctype.Bool(tb.EsEnable) {
		return nil
	}

	// 避免临时调试，把es启用关闭，而恰好有数据更新，导致ES不同步，所以这里只要是提供了ES服务就进行同步
	rsIdxName := entity.GetEsIndexNameByDb(txDb, el.Index(0).Interface())
	if !autoCommit {
		var (
			c      = dorm.GetDbContext(txDb)
			userId = ""
		)
		if c != nil {
			userId = c.GetContextUserId()
		}
		event := broker.NewEvent(userId, entityList, func(event *broker.Event) {
			if delete {
				event.OptType = broker.EventOptTypeDelete
			} else {
				event.OptType = broker.EventOptTypeUpdate
			}
			event.TableName = rsIdxName
		})
		//eventIds := []string{event.Id}
		//dorm.Store(txDb, constants.BrokerSync2Es, eventIds)
		//return broker.Publish(context.Background(), constants.BrokerSync2Es, event, false)
		dorm.Store(txDb, constants.BrokerSync2Es, event) //只是为了做事务回滚的时候删除
	}
	return sync2Es(rsIdxName, entityList, delete)
}

func sync2Es(idxName string, entityList any, delete bool) error {
	if delete {
		return es.NewApi(global.GetES(), idxName).Delete(entityList)
	} else {
		return es.NewApi(global.GetES(), idxName).Upsert(entityList)
	}
}
