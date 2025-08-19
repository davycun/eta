package plugin_crypt

import (
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/core/service/hook"
	"reflect"
)

// StoreEncrypt 在数据插入和修改前，加密
func StoreEncrypt(cfg *hook.SrvConfig, pos hook.CallbackPosition) (err error) {
	if pos != hook.CallbackBefore || (cfg.Method != iface.MethodCreate && cfg.Method != iface.MethodUpdate && cfg.Method != iface.MethodUpdateByFilters) {
		return
	}
	table := cfg.GetTable()
	if len(table.CryptFields) < 1 {
		return
	}
	for _, cryptInfo := range table.CryptFields {
		err = encryptValue(cfg.TxDB, table.GetTableName(), cryptInfo, cfg.Values...)
		if err != nil {
			return err
		}
	}
	return
}

// StoreDecrypt 在数据查询后，解密
func StoreDecrypt(cfg *hook.SrvConfig, pos hook.CallbackPosition) error {
	return hook.AfterRetrieve(cfg, pos, func(cfg *hook.SrvConfig, rs []reflect.Value) error {
		table := cfg.GetTable()
		if len(table.CryptFields) < 1 {
			return nil
		}

		for _, cryptInfo := range table.CryptFields {
			if !cryptInfo.Enable || cryptInfo.Field == "" {
				continue
			}
			err := decryptValue(cfg.TxDB, cfg.GetTableName(), cryptInfo, rs...)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
