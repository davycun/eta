package reload

import (
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/iface"
)

var (
	hooks                           = map[RdType]map[string]Callback{} // rdType -> tableName ->  callback
	CallbackBefore CallbackPosition = 1
	CallbackAfter  CallbackPosition = 2
)

const (
	RdTypeDb2Es   RdType = "db2es"
	RdTypeFeature RdType = "feature"
)

type (
	RdType           string
	CallbackPosition int
	RdService        struct {
		srv    iface.Service
		param  *dto.Param
		result *dto.Result
	}

	Callback func(cfg *RdService, pos CallbackPosition) error
)

func (r RdService) GetService() iface.Service {
	return r.srv
}
func (r RdService) GetParam() *dto.Param {
	return r.param
}
func (r RdService) GetResult() *dto.Result {
	return r.result
}

func Registry(tableName string, rdType RdType, callback Callback) {
	mp, ok := hooks[rdType]
	if !ok {
		mp = map[string]Callback{}
	}
	mp[tableName] = callback
	hooks[rdType] = mp
}

func callbackBefore(rdType RdType, rds *RdService) error {
	_ = operateTrigger(rds.srv, false)
	if cf, ok := hooks[rdType]; ok {
		if fc, ok1 := cf[rds.GetService().GetTableName()]; ok1 {
			return fc(rds, CallbackBefore)
		}
	}
	return nil
}
func callbackAfter(rdType RdType, rds *RdService) error {
	_ = operateTrigger(rds.srv, true)
	if cf, ok := hooks[rdType]; ok {
		if fc, ok1 := cf[rds.GetService().GetTableName()]; ok1 {
			return fc(rds, CallbackAfter)
		}
	}
	return nil
}
