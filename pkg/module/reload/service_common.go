package reload

import (
	"context"
	"errors"
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/filter"
	"github.com/davycun/eta/pkg/common/dorm/xa"
	"github.com/davycun/eta/pkg/common/dsync"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/core/service"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/eta/plugin/plugin_es"
	"reflect"
)

func DbLoader(args any) (data any, over bool, err error) {
	var (
		srvArgs = args.(*dsync.SyncArgs)
		srv     = srvArgs.Srv
		param   = srvArgs.Args
		extra   = dto.GetExtra[dsync.SyncOption](param)
	)
	if !replaceEidFilter(param.Filters, extra.StartEid) && extra.StartEid != 0 {
		param.Filters = append(param.Filters, filter.Filter{
			LogicalOperator: filter.And,
			Column:          entity.EIdDbName,
			Operator:        filter.GT,
			Value:           extra.StartEid,
		})
	}

	var (
		entPtr      = srv.NewEntityPointer()
		entSlicePtr = srv.NewEntitySlicePointer()
		tableName   = entity.GetTableName(entPtr)
	)

	logger.Debugf("db loader %s startId: %d", tableName, extra.StartEid)
	if len(param.Columns) > 0 {
		param.Columns = utils.Merge(param.Columns, entity.EIdDbName)
	}

	listSql, _, err := service.BuildParamSql(srv.GetDB(), param, srv.GetEntityConfig())
	if err != nil {
		return nil, true, err
	}

	err = dorm.RawFetch(listSql, srv.GetDB(), entSlicePtr)
	if err != nil {
		return nil, true, err
	}

	elem := reflect.ValueOf(entSlicePtr).Elem()
	if elem.Kind() == reflect.Slice {
		if elem.Len() <= 0 || elem.Len() < param.GetLimit() {
			over = true
		}
		if elem.Len() > 0 {
			extra.StartEid = elem.Index(elem.Len() - 1).FieldByName("EID").Int()
		} else {
			entSlicePtr = nil
		}
	} else {
		logger.Warnf("data is not slice")
		over = true
		entSlicePtr = nil
	}
	return entSlicePtr, over, err
}

func replaceEidFilter(flt []filter.Filter, startEid int64) bool {
	if startEid == 0 {
		return false
	}
	flag := false
	for i, _ := range flt {
		f := &flt[i]
		if f.Column == "eid" && f.Operator == filter.GT {
			f.Value = startEid
			flag = true
		}
		if len(f.Filters) > 0 {
			replaceEidFilter(f.Filters, startEid)
		}
	}
	return flag
}

func EsSaver(args any, data any) error {
	var (
		err       error
		sa        = args.(*dsync.SyncArgs)
		srv       = sa.Srv
		db        = srv.GetDB()
		tb        = srv.GetTable()
		tbName    = srv.GetTableName()
		indexName = entity.GetEsIndexNameByDb(db, tbName)
		convert   = plugin_es.GetConvert(tbName, tbName)
	)
	txData := &xa.TxData{
		Delete:       false,
		EsIndexName:  indexName,
		RollbackData: data,
		TargetData:   data,
	}

	if convert != nil {
		cfg := hook.NewSrvConfig(iface.CurdRetrieve, "db2es", iface.NewSrvOptionsFromService(srv), sa.Args, nil)
		txData, err = convert(cfg, txData)
		if err != nil {
			return err
		}
	}

	return plugin_es.Sync2Es(db, tb, txData, true)
}

func operateTrigger(srv iface.Service, enableTrigger bool) error {
	var (
		err            error
		tableName      = entity.GetTableName(srv.NewEntityPointer())
		triggerHistory = fmt.Sprintf(`trigger_%s%s`, tableName, constants.TableHistorySubFix)
		triggerUpdater = fmt.Sprintf(`trigger_%s_updater`, tableName)
		triggerRa      = fmt.Sprintf(`trigger_%s_ra`, tableName)
		db             = srv.GetDB()
	)
	for _, tg := range []string{triggerHistory, triggerUpdater, triggerRa} {
		if exists, _ := dorm.TriggerExists(db, tg, tableName); !exists {
			logger.Debugf("操作触发器[%s]不存在, 不执行任何操作", tg)
			continue
		}
		schTg := fmt.Sprintf("%s.%s", dorm.GetDbSchema(db), tg)
		if enableTrigger {
			err = dorm.TriggerEnable(srv.GetDB(), schTg, dorm.GetScmTableName(db, tableName))
		} else {
			err = dorm.TriggerDisable(srv.GetDB(), schTg, dorm.GetScmTableName(db, tableName))
		}
		if err != nil {
			logger.Debugf("操作触发器[%s]出错, enableTrigger=%v", schTg, enableTrigger)
		}
	}
	return nil
}

func checkEsIndex(srv iface.Service) error {
	var (
		db        = srv.GetDB()
		entPtr    = srv.NewEntityPointer()
		indexName = entity.GetEsIndexNameByDb(db, entPtr)
	)

	if global.GetES() == nil || !srv.GetTable().EsRetrieveEnabled() {
		return errors.New("不支持ES索引")
	}

	exists, err := global.GetES().EsTypedApi.Indices.Exists(indexName).Do(context.Background())
	if err != nil {
		logger.Warnf("indics exists err %s", err)
	}
	if exists {
		return nil
	}
	exists, err = global.GetES().EsTypedApi.Indices.ExistsAlias(indexName).Do(context.Background())
	if err != nil {
		logger.Warnf("indics alias exists err %s", err)
	}
	if exists {
		return nil
	}
	return errors.New(fmt.Sprintf("indics [%s] does not exist,", indexName))
}
