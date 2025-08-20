package plugin_es

import (
	"github.com/davycun/eta/pkg/common/dorm/xa"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/eta/ecf"
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

func GetConvert(fromTable, toTable string) Convert {
	if mp, ok := convertMap[fromTable]; ok {
		return mp[toTable].convert
	}
	return nil
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

	if len(cvt) < 1 && cfg.GetTable().EsEnabled() {
		return Sync2Es(cfg.TxDB, cfg.GetTable(), txData, false)
	}

	for _, v := range cvt {
		ec, b := ecf.GetEntityConfig(appDb, v.toTable)
		if !b {
			continue
		}
		txData.EsIndexName = entity.GetEsIndexNameByDb(appDb, v.toTable)
		targetData, err := v.convert(cfg, txData)
		if err != nil {
			return err
		}
		//这里传入的Table是对应的toTable的
		if err = Sync2Es(cfg.TxDB, ec.GetTable(), targetData, false); err != nil {
			return err
		}
	}
	return nil
}
