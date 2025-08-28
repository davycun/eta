package ecf

import (
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/module/setting"
	"github.com/davycun/eta/pkg/module/template"
	"github.com/duke-git/lancet/v2/slice"
	"gorm.io/gorm"
	"slices"
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

	if appDb == nil {
		return cfg1, false
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

func GetEntityConfigCtxOrSetting(c *ctx.Context, appDb *gorm.DB, key string) (iface.EntityConfig, bool) {

	if ec := GetContextEntityConfig(c); ec != nil {
		return *ec, true
	}
	return GetEntityConfig(appDb, key)
}

func GetMigrateAppTable(db *gorm.DB, namespace ...string) []entity.Table {
	var (
		toList = make([]entity.Table, 0)
		ecList = getMergedEntityConfig(db)
	)

	for _, v := range ecList {
		if v.Migrate && v.LocatedApp() && len(namespace) == 0 || slice.Contain(namespace, v.Namespace) {
			toList = append(toList, v.Table)
		}
	}
	slices.SortFunc(toList, func(a, b entity.Table) int {
		return b.Order - a.Order
	})
	return toList
}

// GetMigrateLocalTable
// 返回需要再localDB中创建表的实体
func GetMigrateLocalTable(db *gorm.DB, namespace ...string) []entity.Table {

	var (
		toList = make([]entity.Table, 0)
		ecList = getMergedEntityConfig(db)
	)

	for _, v := range ecList {
		if v.Migrate && v.LocatedLocal() && len(namespace) == 0 || slice.Contain(namespace, v.Namespace) {
			toList = append(toList, v.Table)
		}
	}
	slices.SortFunc(toList, func(a, b entity.Table) int {
		return b.Order - a.Order
	})
	return toList
}

// GetEsTable
// 返回需要再appDB中创建表的实体
func GetEsTable(db *gorm.DB, namespace ...string) []entity.Table {

	var (
		toList = make([]entity.Table, 0)
		ecList = getMergedEntityConfig(db)
	)

	for _, v := range ecList {
		if v.EsEnabled() && len(namespace) == 0 || slice.Contain(namespace, v.Namespace) {
			toList = append(toList, v.Table)
		}
	}
	slices.SortFunc(toList, func(a, b entity.Table) int {
		return b.Order - a.Order
	})
	return toList
}

func getMergedEntityConfig(db *gorm.DB) []iface.EntityConfig {

	var (
		ecList = iface.GetEntityConfigList()
	)
	for i, _ := range ecList {
		ec := &ecList[i]
		tbSet, tbSetExists := setting.GetTableConfig(db, ec.GetTableName()) //合同配置中心动态的配置
		if tbSetExists {
			ec.Merge(&tbSet)
		}
	}

	slices.SortFunc(ecList, func(a, b iface.EntityConfig) int {
		return b.Order - a.Order
	})
	return ecList
}
