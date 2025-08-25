package plugin_crypt

import (
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/module/setting"
	"gorm.io/gorm"
)

type codecPlugin struct{}

type Option func(p *codecPlugin)

func NewEncryptPlugin(opts ...Option) gorm.Plugin {
	p := &codecPlugin{}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// Name return plugin name
func (p codecPlugin) Name() string {
	return "crypt_plugin"
}

// Initialize initialize plugin
func (p codecPlugin) Initialize(db *gorm.DB) error {
	// 1. 获取db的Callback对象
	cb := db.Callback()

	// 2. 将需要执行的操作以及对应的回调函数注册到Callback对象
	hooks := []struct {
		name   string
		action error
	}{
		// 创建操作前的回调函数
		{name: "before:create", action: cb.Create().Before("gorm:create").Register("crypt:before:create", storeEncrypt)},
		// 创建操作后的回调函数
		//{name: "after:create", action: cb.Create().After("gorm:create").Register("codec:after:create", callbackCodec("after:create", ProcessAfter))},
		// 更新操作前的回调函数
		{name: "before:update", action: cb.Update().Before("gorm:update").Register("crypt:before:update", storeEncrypt)},
		// 更新操作后的回调函数
		//{name: "after:update", action: cb.Update().After("gorm:update").Register("codec:after:update", callbackCodec("after:update", ProcessAfter))},
		// 查询操作后的回调函数
		//{name: "after:query", action: cb.Query().After("gorm:query").Register("codec:after:query", callbackCodec("after:query", ProcessAfter))},
	}

	for _, h := range hooks {
		err := h.action
		if err != nil {
			return err
		}
	}

	return nil
}

// 注意这个函数内不能直接用传入的db的原因是，这个db可能是create或者updater等传进来的
// 如果用当前db去做查询就会导致正在操作的create或者update出问题
func storeEncrypt(db *gorm.DB) {

	if db.Error != nil || db.Statement.Table == "" {
		return
	}

	var (
		tableName = db.Statement.Table
		val       = db.Statement.ReflectValue
		appDb, _  = global.LoadGormByAppId(dorm.GetAppId(db))
	)
	tb, b := setting.GetTableConfig(appDb, tableName)

	if !b || len(tb.CryptFields) < 1 {
		return
	}
	valSlice := utils.ConvertToValueArray(val)
	for _, v := range tb.CryptFields {
		err := encryptValue(appDb, tableName, v, valSlice...)
		if err != nil {
			logger.Errorf("storeEncrypt err %s", err)
			return
		}
	}
	return
}
