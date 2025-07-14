package iface

import (
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/duke-git/lancet/v2/slice"
	"reflect"
	"slices"
	"strings"
)

var (
	entityNameConfigMap = make(map[string]EntityConfig)
	tableNameConfigMap  = make(map[string]EntityConfig)
	baseUrlConfigMap    = make(map[string]EntityConfig)
)

func Registry(conf ...EntityConfig) {
	if len(conf) <= 0 {
		return
	}
	var (
		name, tableName, baseUrl string
	)
	for _, v := range conf {

		if v.RsDataType == nil {
			v.RsDataType = v.EntityType
		}
		//不要指针类型，要具体的结构体类型
		if v.ServiceType != nil && v.ServiceType.Kind() == reflect.Pointer {
			v.ServiceType = v.ServiceType.Elem()
		}
		//不要指针类型，要具体的结构体类型
		if v.ControllerType != nil && v.ControllerType.Kind() == reflect.Pointer {
			v.ControllerType = v.ControllerType.Elem()
		}

		name, tableName, baseUrl = v.Name, entity.GetTableName(v.NewEntityPointer()), v.BaseUrl

		if tableName == "" {
			logger.Warnf("EntityConfig[name:%s,base_url:%s] tableName will be set to name because it's  empty", tableName, baseUrl)
			tableName = name
		}

		if _, ok := entityNameConfigMap[name]; ok {
			logger.Errorf("EntityConfig repeated name %s", name)
			continue
		}
		if _, ok := tableNameConfigMap[tableName]; ok {
			logger.Errorf("EntityConfig repeated tableName %s", tableName)
			continue
		}

		if !v.DisableApi {
			if _, ok := baseUrlConfigMap[baseUrl]; ok {
				logger.Errorf("EntityConfig repeated baseUrl %s", baseUrl)
				continue
			}
		}
		//提前初始化
		v.GetTable()
		entityNameConfigMap[name] = v
		tableNameConfigMap[tableName] = v

		//有disableApi的情况
		if baseUrl != "" {
			baseUrlConfigMap[baseUrl] = v
		}
	}
}

func GetEntityConfigList() []EntityConfig {
	entityConfigList := make([]EntityConfig, 0, len(entityNameConfigMap))
	for _, v := range entityNameConfigMap {
		entityConfigList = append(entityConfigList, v)
	}
	return entityConfigList
}

func GetEntityConfigByName(name string) (EntityConfig, bool) {
	x, ok := entityNameConfigMap[name]
	if !ok {
		x, ok = tableNameConfigMap[name]
	}
	return x, ok
}
func GetEntityConfigByTableName(tbName string) (EntityConfig, bool) {

	x, ok := tableNameConfigMap[tbName]
	if !ok {
		x, ok = entityNameConfigMap[tbName]
	}
	return x, ok
}
func GetEntityConfigByUrl(fullUrl string) (EntityConfig, bool) {

	uls := strings.Split(fullUrl, "/")

	if len(uls) < 1 {
		return EntityConfig{}, false
	}
	for i := len(uls) - 1; i >= 0; i-- {
		ul := strings.Join(uls[:i], "/")
		if ec, ok := baseUrlConfigMap[ul]; ok {
			return ec, true
		}
		if ec, ok := baseUrlConfigMap[ul+"/"]; ok {
			return ec, true
		}
	}
	return EntityConfig{}, false
}

func GetMigrateEntityConfig(namespace ...string) []entity.Table {
	toList := make([]entity.Table, 0, len(entityNameConfigMap))
	for _, v := range entityNameConfigMap {
		if v.Migrate && (len(namespace) == 0 || slice.Contain(namespace, v.Namespace)) {
			toList = append(toList, v.Table)
		}
	}
	slices.SortFunc(toList, func(a, b entity.Table) int {
		return b.Order - a.Order
	})
	return toList
}

// GetMigrateLocalEntityConfig
// 返回需要再localDB中创建表的实体
func GetMigrateLocalEntityConfig(namespace ...string) []entity.Table {
	toList := make([]entity.Table, 0, len(entityNameConfigMap))
	for _, v := range entityNameConfigMap {
		if v.Migrate && (len(namespace) == 0 || slice.Contain(namespace, v.Namespace)) && v.LocatedLocal() {
			toList = append(toList, v.Table)
		}
	}
	slices.SortFunc(toList, func(a, b entity.Table) int {
		return b.Order - a.Order
	})
	return toList
}

// GetMigrateAppEntityConfig
// 返回需要再appDB中创建表的实体
func GetMigrateAppEntityConfig(namespace ...string) []entity.Table {
	toList := make([]entity.Table, 0, len(entityNameConfigMap))
	for _, v := range entityNameConfigMap {
		if v.Migrate && (len(namespace) == 0 || slice.Contain(namespace, v.Namespace)) && v.LocatedApp() {
			toList = append(toList, v.Table)
		}
	}
	slices.SortFunc(toList, func(a, b entity.Table) int {
		return b.Order - a.Order
	})
	return toList
}
