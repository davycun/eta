package reload

import (
	"context"
	"errors"
	"fmt"
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/es"
	"github.com/davycun/eta/pkg/common/dorm/filter"
	"github.com/davycun/eta/pkg/common/dsync"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/core/loader"
	"github.com/davycun/eta/pkg/core/ra"
	"github.com/davycun/eta/pkg/eta/constants"
	"reflect"
	"strings"
	"time"

	"github.com/duke-git/lancet/v2/maputil"
	"github.com/duke-git/lancet/v2/slice"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var (
	silentLogger = gormLogger.New(logger.Logger, gormLogger.Config{LogLevel: gormLogger.Silent})
	slowLogger   = gormLogger.New(logger.Logger, gormLogger.Config{LogLevel: gormLogger.Warn, SlowThreshold: 30 * time.Second})
)

func BeforeReload(args *dsync.SyncArgs) error {
	return caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			return operateTrigger(args.Srv, false)
		}).
		Call(func(cl *caller.Caller) error {
			if reLoadSrv, ok := args.Srv.(dsync.ReloadInjector); ok {
				return reLoadSrv.ReloadBefore(args)
			}
			return nil
		}).Err
}
func AfterReload(args *dsync.SyncArgs) error {
	return caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			if reLoadSrv, ok := args.Srv.(dsync.ReloadInjector); ok {
				return reLoadSrv.ReloadAfter(args)
			}
			return nil
		}).
		Call(func(cl *caller.Caller) error {
			return operateTrigger(args.Srv, true)
		}).Err

}

func QueryLoader(args any) (data any, over bool, err error) {
	var (
		sa = args.(*dsync.SyncArgs)
		//so        = sa.Args.Extra.(*dsync.SyncOption)
		so        = dto.GetExtra[dsync.SyncOption](sa.Args)
		srv       = sa.Srv
		tableName = entity.GetTableName(srv.NewEntityPointer())
	)
	buildParam := func(p dto.Param, newId string) dto.Param {
		fs := []filter.Filter{{
			LogicalOperator: filter.And,
			Column:          entity.IdDbName,
			Operator:        filter.GT,
			Value:           newId,
			Filters:         p.Filters,
		}}
		p.Filters = fs
		p.AutoCount = false
		p.OnlyCount = false
		p.Columns = []string{}
		return p
	}
	logger.Debugf("query loader %s startId: %s", tableName, so.StartId)

	// 如果有存储加密字段，查询出来的数据需要解密成明文
	param := buildParam(*sa.Args, so.StartId)
	result := dto.Result{}
	err = srv.Query(&param, &result)
	elem := reflect.ValueOf(result.Data).Elem()
	if elem.Kind() == reflect.Slice {
		if elem.Len() <= 0 || elem.Len() < sa.Args.GetLimit() {
			over = true
		}
		if elem.Len() > 0 {
			so.StartId = elem.Index(elem.Len() - 1).FieldByName(entity.IdFieldName).String()
		} else {
			result.Data = nil
		}
	}
	if err != nil {
		return nil, false, err
	}
	return result.Data, over, err
}
func EsLoader(args any) (data any, over bool, err error) {
	var (
		sa  = args.(*dsync.SyncArgs)
		so  = dto.GetExtra[dsync.SyncOption](sa.Args)
		srv = sa.Srv
		db  = srv.GetDB()
		ct  = srv.GetContext()
	)

	buildParam := func(p dto.Param, newId string) dto.Param {
		fs := []filter.Filter{{
			LogicalOperator: filter.And,
			Column:          entity.IdDbName,
			Operator:        filter.GT,
			Value:           newId,
			Filters:         p.Filters,
		}}
		p.Filters = fs
		p.AutoCount = false
		p.OnlyCount = false
		p.Columns = []string{}
		return p
	}
	// 如果有存储加密字段，查询出来的数据需要解密成明文
	param := buildParam(*sa.Args, so.StartId)
	var (
		entPtr      = srv.NewEntityPointer()
		entSlicePtr = srv.NewEntitySlicePointer()
		tableName   = entity.GetTableName(entPtr)
	)
	logger.Debugf("es2db loader %s startId: %s", tableName, so.StartId)

	esApi := es.NewApi(global.GetES(), entity.GetEsIndexNameByDb(db, entPtr),
		es.CodecOpt(dorm.GetDbSchema(ct.GetAppGorm()), tableName))
	esApi.AddFilters(ra.KeywordToFilters(ct.GetAppGorm(), tableName, param.SearchContent, dorm.ES)...)
	esApi.AddFilters(param.Filters...)
	esApi.OrderBy(param.OrderBy...).Limit(param.GetLimit()).Find(entSlicePtr)
	err = esApi.Err

	elem := reflect.ValueOf(entSlicePtr).Elem()
	if elem.Kind() == reflect.Slice {
		if elem.Len() <= 0 || elem.Len() < sa.Args.GetLimit() {
			over = true
		}
		if elem.Len() > 0 {
			so.StartId = elem.Index(elem.Len() - 1).FieldByName(entity.IdFieldName).String()
		} else {
			entSlicePtr = nil
		}
	}
	if err != nil {
		return nil, false, err
	}
	return entSlicePtr, over, err
}
func DbLoader(args any) (data any, over bool, err error) {
	var (
		sa     = args.(*dsync.SyncArgs)
		so     = dto.GetExtra[dsync.SyncOption](sa.Args)
		srv    = sa.Srv
		ct     = srv.GetContext()
		db     = srv.GetDB()
		dbType = dorm.GetDbType(db)
	)

	buildParam := func(p dto.Param, newId string) dto.Param {
		fs := []filter.Filter{{
			LogicalOperator: filter.And,
			Column:          entity.IdDbName,
			Operator:        filter.GT,
			Value:           newId,
			Filters:         p.Filters,
		}}
		p.Filters = fs
		p.AutoCount = false
		p.OnlyCount = false
		p.Columns = []string{}
		return p
	}
	// 如果有存储加密字段，查询出来的数据需要解密成明文
	param := buildParam(*sa.Args, so.StartId)
	var (
		entPtr      = srv.NewEntityPointer()
		entSlicePtr = srv.NewEntitySlicePointer()
		tableName   = entity.GetTableName(entPtr)
	)
	logger.Debugf("db loader %s startId: %s", tableName, so.StartId)

	tx := dorm.WithContext(db, ct).Model(entSlicePtr)
	if so.StartId != "" {
		tx = tx.Where(fmt.Sprintf(`%s > ?`, dorm.Quote(dbType, entity.IdDbName)), so.StartId)
	}
	if len(param.Filters) > 0 {
		tx = tx.Where(filter.ResolveWhere(sa.Args.Filters, dbType))
	}
	tx = tx.Order(param.ResolveOrderByString(tableName, "", false))
	err = tx.Limit(param.GetLimit()).Find(entSlicePtr).Error

	elem := reflect.ValueOf(entSlicePtr).Elem()
	if elem.Kind() == reflect.Slice {
		if elem.Len() <= 0 || elem.Len() < sa.Args.GetLimit() {
			over = true
		}
		if elem.Len() > 0 {
			so.StartId = elem.Index(elem.Len() - 1).FieldByName(entity.IdFieldName).String()
		} else {
			entSlicePtr = nil
		}
	} else {
		logger.Warnf("data is not slice")
		over = true
		entSlicePtr = nil
	}
	if err != nil {
		return nil, false, err
	}
	return entSlicePtr, over, err
}

func EsSaver(args any, data any) error {
	var (
		sa            = args.(*dsync.SyncArgs)
		srv           = sa.Srv
		so            = dto.GetExtra[dsync.SyncOption](sa.Args)
		db            = srv.GetDB()
		pa            []string
		existsDataIdx = make([]int, 0)
		indexName     = entity.GetEsIndexNameByDb(db, srv.NewEntityPointer())
	)
	if utils.IsEmptySlice(data) {
		return nil
	}

	//vals := service.ResolveValue(reflect.ValueOf(data))
	vals := utils.ConvertToValueArray(data)
	ids := make(map[string]int)
	for i, val := range vals {
		ids[entity.GetString(val.Interface(), entity.IdDbName)] = i
	}
	if !so.Upsert {
		ld := loader.NewKeyLoader(db, loader.KeyLoaderConfig{
			TableName: entity.GetTableName(srv.NewEntityPointer()),
			IndexName: indexName,
			KeyColumn: entity.IdDbName,
			Keys:      maputil.Keys(ids),
		})
		err := ld.LoadFromEs(&pa)
		if err != nil {
			return err
		}
		existsDataIdx = slice.Map(pa, func(index int, item string) int {
			if v, ok := ids[item]; ok {
				return v
			}
			return -1
		})
		existsDataIdx = slice.Filter(existsDataIdx, func(index int, item int) bool { return item >= 0 })
	}

	dataList, err := utils.SliceRemoveElemByIndexes(data, existsDataIdx)
	if err != nil {
		return err
	}
	if utils.IsEmptySlice(dataList) {
		return nil
	}

	err = es.NewApi(global.GetES(), indexName).Upsert(dataList)
	if err != nil {
		return err
	}

	return err
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
			err = dorm.TriggerEnable(srv.GetDB(), schTg, dorm.GetDbTable(db, tableName))
		} else {
			err = dorm.TriggerDisable(srv.GetDB(), schTg, dorm.GetDbTable(db, tableName))
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

	if global.GetES() == nil || !srv.GetTable().EnableRetrieveEs() {
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

// 更新数据库
func dbUpdateFeature(args *dsync.SyncArgs, data any) error {
	if utils.IsEmptySlice(data) {
		return nil
	}
	var (
		so   = args.Args.Extra.(*dsync.SyncOption)
		srv  = args.Srv
		cols = make([]string, 0)
		//tb           = data.GetContextTable(srv.GetContext())
		ec           = iface.GetContextEntityConfig(srv.GetContext())
		tb           = ec.GetTable()
		tableName    = entity.GetTableName(srv.NewEntityPointer())
		isHistoryTbl = strings.HasSuffix(tableName, constants.TableHistorySubFix)
	)
	if so.UpdateDbRaContent {
		cols = append(cols, entity.RaContentDbName)
	}
	if so.UpdateDbEncrypt {
		for _, v := range tb.CryptFields {
			cols = utils.Merge(cols, v.Field)
		}
	}
	if so.UpdateDbSign {
		for _, v := range tb.SignFields {
			cols = utils.Merge(cols, v.Field)
		}
	}
	if isHistoryTbl {
		cols = slice.Map(cols, func(_ int, item string) string { return fmt.Sprintf("h_%s", item) })
	}

	// data 切片
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}
	if val.Kind() == reflect.Slice && val.Len() > 0 {
		dataChunks, _ := utils.Chunk(data, so.GetUpdateDbBatchSize())
		for _, dataChunk := range dataChunks {
			tx := srv.GetDB().Session(&gorm.Session{NewDB: true, Logger: slowLogger})
			err := dorm.BatchUpdate(tx, dataChunk, entity.IdDbName, cols...)
			if err != nil {
				return err
			}
		}
		sl, _ := utils.SliceLen(data)
		logger.Debugf("写入DB %d 条数据完成", sl)
		return nil
	}
	return nil
}
