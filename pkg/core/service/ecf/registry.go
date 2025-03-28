package ecf

import (
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/entity"
	"reflect"
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
