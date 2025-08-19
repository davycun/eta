package plugin_geo

import (
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/core/service/hook"
	"reflect"
)

func ProcessGeometryResult(cfg *hook.SrvConfig, pos hook.CallbackPosition) error {
	return hook.AfterRetrieve(cfg, pos, func(cfg *hook.SrvConfig, rs []reflect.Value) error {
		if cfg.Param.GeoGCSType != "" || cfg.Param.GeoFormat != "" {
			for _, v := range rs {
				if v.IsValid() && v.CanAddr() {
					eAddr := v.Addr()
					if eAddr.CanInterface() {
						e := eAddr.Interface()
						if gf, ok := e.(ctype.GeometryFormat); ok {
							gf.GeomFormat(cfg.Param.GeoGCSType, cfg.Param.GeoFormat)
						}
					}
				}
			}
		}
		return nil
	})
}
