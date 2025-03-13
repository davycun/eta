package plugin_es

import (
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/service/hook"
	"sync"
)

var (
	convertMap = sync.Map{} //tableName -> ConvertFunction
)

// Convert
// 把实体列表转换成es的实体列表，如果没有指定
type Convert func(cfg *hook.SrvConfig, entityList any) (any, error)

func RegisterConvert(tableName string, convert Convert) {
	convertMap.Store(tableName, convert)
}
func RemoveConvert(tableName string) {
	convertMap.Delete(tableName)
}

func ConvertEsEntity(cfg *hook.SrvConfig, entityList any) (any, error) {
	tbName := cfg.GetTableName()
	if cvt, ok := convertMap.Load(tbName); ok {
		if fc, ok1 := cvt.(Convert); ok1 {
			return fc(cfg, entityList)
		}
	}
	logger.Warnf("can not found the convert function of entity[%s] to esEntity", tbName)
	return entityList, nil
}
