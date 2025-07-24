package migrate

import (
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/migrate/mig_hook"
	"github.com/davycun/eta/pkg/core/ra"
)

func init() {
	mig_hook.AddCallback(mig_hook.CallbackForAll, afterMigrator)
}

// 这个Callback主要是处理创建表之后，去调用那些实现了MigratorAfter接口的实体以及创建带有ra特性的触发器
func afterMigrator(mc *mig_hook.MigConfig, pos mig_hook.CallbackPosition) error {
	if pos != mig_hook.CallbackAfter {
		return nil
	}
	var (
		t   = mc.TbOption
		tx  = mc.TxDB
		c   = mc.C
		val = t.NewEntityPointer()
	)
	if val == nil {
		return nil
	}
	//TODO 这里其实需要考虑dst中有重复的情况
	if ma, ok := val.(MigratorAfter); ok {
		if err := ma.AfterMigrator(tx, c); err != nil {
			return err
		}
	}

	if ma, ok := val.(entity.RaInterface); ok {
		raFields := ma.RaDbFields()
		if err := ra.CreateTrigger(tx, entity.GetTableName(val), raFields); err != nil {
			return err
		}
	}
	return nil
}
