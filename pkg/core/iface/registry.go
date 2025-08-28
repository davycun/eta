package iface

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/entity"
	"strings"
)

var (
	entityNameConfigMap = make(map[string]EntityConfig)
	tableNameConfigMap  = make(map[string]EntityConfig)
	baseUrlConfigMap    = make(map[string]EntityConfig)
	allEntityConfigMap  = make(map[string]EntityConfig) //key = name + tableName + baseUrl
)

func Registry(conf ...EntityConfig) {
	if len(conf) <= 0 {
		return
	}
	var (
		name, tableName, baseUrl string
	)
	for _, v := range conf {
		//不要指针类型，要具体的结构体类型
		v.ServiceType = utils.GetRealType(v.ServiceType)
		//不要指针类型，要具体的结构体类型
		v.ControllerType = utils.GetRealType(v.ControllerType)

		name, tableName, baseUrl = v.Name, entity.GetTableName(v.NewEntityPointer()), v.BaseUrl

		if err := check(v); err != nil {
			logger.Errorf("%s", err)
			continue
		}

		if name != "" {
			entityNameConfigMap[name] = v
		}

		if tableName != "" {
			//提前初始化
			v.GetTable()
			tableNameConfigMap[tableName] = v
		}
		if baseUrl != "" {
			baseUrlConfigMap[baseUrl] = v
		}

		if x := v.GetKey(); x != "" {
			allEntityConfigMap[x] = v
		}

	}
}

func check(ec EntityConfig) error {
	var (
		name    = ec.Name
		baseUrl = ec.BaseUrl
		tbName  = ec.GetTableName()
	)
	if _, ok := allEntityConfigMap[ec.GetKey()]; ok {
		return fmt.Errorf("EntityConfig[name:%s,table_name:%s,base_url:%s] had Exists", name, tbName, baseUrl)
	}
	if tbName == "" && ec.Migrate {
		return fmt.Errorf("EntityConfig[name:%s,table_name:%s,base_url:%s] if migrate is true, you should set tableName", name, tbName, baseUrl)
	}
	if ec.BaseUrl == "" && !ec.DisableApi {
		return fmt.Errorf("EntityConfig[name:%s,table_name:%s,base_url:%s] if disableApi is false,you should set baseUrl", name, tbName, baseUrl)
	}
	return nil
}

func GetEntityConfigList() []EntityConfig {
	entityConfigList := make([]EntityConfig, 0, len(allEntityConfigMap))
	//这里需要用baseUrlConfig,比较全
	for _, v := range allEntityConfigMap {
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
	for i := len(uls); i > 0; i-- {
		ul := strings.Join(uls[:i], "/")
		if ec, ok := baseUrlConfigMap[ul]; ok {
			return ec, true
		}
		if ec, ok := baseUrlConfigMap[fmt.Sprintf("%s/", ul)]; ok {
			return ec, true
		}
		if ec, ok := baseUrlConfigMap[fmt.Sprintf("/%s", ul)]; ok {
			return ec, true
		}
		if ec, ok := baseUrlConfigMap[fmt.Sprintf("/%s/", ul)]; ok {
			return ec, true
		}
	}
	return EntityConfig{}, false
}
func GetEntityConfigByKey(key string) (EntityConfig, bool) {
	ec, b := GetEntityConfigByTableName(key)
	if b {
		return ec, true
	}
	ec, b = GetEntityConfigByName(key)
	if b {
		return ec, true
	}
	ec, b = GetEntityConfigByUrl(key)
	return ec, b
}
