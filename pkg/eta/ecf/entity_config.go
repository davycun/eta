package ecf

import (
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/module/setting"
	"github.com/davycun/eta/pkg/module/template"
	"gorm.io/gorm"
)

func GetEntityConfig(appDb *gorm.DB, key string) (iface.EntityConfig, bool) {

	var (
		ec                 = iface.EntityConfig{}
		tbSet, tbSetExists = setting.GetTableConfig(appDb, key) //合同配置中心动态的配置
		cfg1, b1           = iface.GetEntityConfigByKey(key)
	)

	//优先级1：从配置中心获取
	if b1 {
		if tbSetExists {
			cfg1.Merge(&tbSet)
		} else {
			//应为传入的key可能是url或者tableName，而配置表里的key是tableName，所以这里需要获取tableName
			tbSet, tbSetExists = setting.GetTableConfig(appDb, cfg1.GetTableName())
			if tbSetExists {
				cfg1.Merge(&tbSet)
			}
		}
		return cfg1, true
	}

	//优先级2：从template获取
	tmp, err := template.LoadByCode(appDb, key)
	if err != nil {
		logger.Errorf("load template error: %v", err)
		return ec, false
	}

	ec.Table = *tmp.GetTable()
	if tbSetExists {
		ec.Merge(&tbSet)
	} else {
		//应为传入的key可能是url或者tableName，而配置表里的key是tableName，所以这里需要获取tableName
		tbSet, tbSetExists = setting.GetTableConfig(appDb, ec.GetTableName())
		if tbSetExists {
			ec.Merge(&tbSet)
		}
	}
	return ec, true
}
