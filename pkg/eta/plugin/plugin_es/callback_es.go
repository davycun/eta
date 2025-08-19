package plugin_es

import (
	"errors"
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/dorm/es"
	"github.com/davycun/eta/pkg/common/dorm/xa"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/service/hook"
	"gorm.io/gorm"
	"reflect"
)

// ModifyCallbackForEs
// modify callback for ES, 同步到 ES
func ModifyCallbackForEs(cfg *hook.SrvConfig, pos hook.CallbackPosition) error {

	err := caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			return hook.AfterCreateAny(cfg, pos, func(cfg *hook.SrvConfig) error {
				txData := &xa.TxData{
					Delete:       false,
					TargetData:   cfg.NewValues,
					RollbackData: cfg.NewValues,
					EsIndexName:  cfg.GetEsIndexName(),
				}
				return convertAndSync2Es(cfg, txData)
			})
		}).
		Call(func(cl *caller.Caller) error {
			return hook.AfterUpdateAny(cfg, pos, func(cfg *hook.SrvConfig) error {
				txData := &xa.TxData{
					Delete:       false,
					TargetData:   cfg.NewValues,
					RollbackData: cfg.OldValues,
					EsIndexName:  cfg.GetEsIndexName(),
				}
				return convertAndSync2Es(cfg, txData)
			})
		}).
		Call(func(cl *caller.Caller) error {
			return hook.AfterDeleteAny(cfg, pos, func(cfg *hook.SrvConfig) error {
				txData := &xa.TxData{
					Delete:       false,
					TargetData:   cfg.OldValues,
					RollbackData: cfg.OldValues,
					EsIndexName:  cfg.GetEsIndexName(),
				}
				return convertAndSync2Es(cfg, txData)
			})
		}).Err

	return err
}

// Sync2Es
// 要确保tb 和 entityList 对应的是同一类实体
func Sync2Es(txDb *gorm.DB, tb *entity.Table, txData *xa.TxData, autoCommit bool) error {
	if txData == nil || txData.TargetData == nil {
		return nil
	}
	el := reflect.ValueOf(txData.TargetData)
	if el.Kind() == reflect.Ptr {
		el = el.Elem()
	}
	if el.Kind() != reflect.Slice {
		return errors.New("entityList must be a slice")
	}
	if el.Len() <= 0 || global.GetES() == nil || !ctype.Bool(tb.EsEnable) {
		return nil
	}

	if !autoCommit {
		dorm.Store(txDb, xa.SyncEsData, txData) //只是为了做事务回滚的时候删除
	}
	return sync2Es(txData)
}

func sync2Es(txData *xa.TxData) error {
	if txData.Delete {
		return es.NewApi(global.GetES(), txData.EsIndexName).Delete(txData.TargetData)
	} else {
		return es.NewApi(global.GetES(), txData.EsIndexName).Upsert(txData.TargetData)
	}
}
