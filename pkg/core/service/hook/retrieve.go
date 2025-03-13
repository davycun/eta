package hook

import "github.com/davycun/eta/pkg/core/iface"

func (cfg *SrvConfig) retrieveBefore(callbacks ...CallbackWrapper) error {
	var (
		err error
	)
	err = callAuthCallback(cfg, CallbackBefore)
	if err != nil {
		return err
	}
	return callCallback(cfg, iface.CurdRetrieve, cfg.Method, CallbackBefore, callbacks...)
}
func (cfg *SrvConfig) retrieveAfter(callbacks ...CallbackWrapper) error {
	return callCallback(cfg, iface.CurdRetrieve, cfg.Method, CallbackAfter, callbacks...)
}
