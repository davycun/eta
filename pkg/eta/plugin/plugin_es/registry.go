package plugin_es

import (
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/dorm/xa"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/module/data/template"
)

var (
	convertMap = map[string]map[string]converter{} //tableName -> ConvertFunction
)

type converter struct {
	fromTable string
	toTable   string
	convert   Convert
}

// Convert
// 把实体列表转换成es的实体列表，如果没有指定
type Convert func(cfg *hook.SrvConfig, txData *xa.TxData) (*xa.TxData, error)

func RegisterConvert(fromTable, toTable string, convert Convert) {

	mp := convertMap[fromTable]
	if mp == nil {
		mp = make(map[string]converter)
	}
	mp[toTable] = converter{
		fromTable: fromTable,
		toTable:   toTable,
		convert:   convert,
	}
	convertMap[fromTable] = mp
}
func RemoveConvert(fromTable, toTable string) {
	if mp, ok := convertMap[fromTable]; ok {
		delete(mp, toTable)
	}
}

func convertAndSync2Es(cfg *hook.SrvConfig, txData *xa.TxData) error {
	var (
		tbName  = cfg.GetTableName()
		cvt, ok = convertMap[tbName]
		appDb   = cfg.GetContext().GetAppGorm()
	)
	if !ok {
		return nil
	}

	if len(cvt) < 1 && ctype.Bool(cfg.GetTable().EsEnable) {
		return Sync2Es(cfg.TxDB, cfg.GetTable(), txData, false)
	}

	for _, v := range cvt {
		tb, exists := iface.GetTableByTableName(v.toTable)
		//常规表名不存在的话，从模板表找
		if !exists && appDb != nil {
			temp, err := template.LoadByCode(appDb, v.toTable)
			if err != nil {
				logger.Errorf("convert to es find template err %s", err)
				continue
			}
			tb = temp.GetTable()
		}
		if tb == nil {
			continue
		}

		targetData, err := v.convert(cfg, txData)
		if err != nil {
			return err
		}

		//这里传入的Table是对应的toTable的
		if err = Sync2Es(cfg.TxDB, tb, targetData, false); err != nil {
			return err
		}
	}
	return nil
}
