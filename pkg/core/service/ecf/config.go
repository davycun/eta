package ecf

import (
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/modern-go/reflect2"
	"reflect"
	"slices"
	"strings"
)

var (
	entityNameConfigMap = make(map[string]EntityConfig)
	tableNameConfigMap  = make(map[string]EntityConfig)
	baseUrlConfigMap    = make(map[string]EntityConfig)
)

type EntityConfig struct {
	entity.Table
	Namespace      string              `json:"namespace,omitempty"`       //区分不同的定制系统或者产品或者模块
	Name           string              `json:"name,omitempty"`            //实体的唯一名字，在事务集合接口中用来唯一标志一个唯一的操作的对象
	Migrate        bool                `json:"migrate,omitempty"`         //是否需要进行migrate
	DisableApi     bool                `json:"disable_api,omitempty"`     //取消暴露API
	ServiceType    reflect.Type        `json:"service_type,omitempty"`    //如果NewService没有，那么就通过类型直接创建
	ControllerType reflect.Type        `json:"controller_type,omitempty"` //如果NewController没有，那么就通过类型直接创建
	NewService     iface.NewService    `json:"new_service,omitempty"`     //服务工厂，需要自定义初始化Service就可以提供这个函数
	NewController  iface.NewController `json:"new_controller,omitempty"`  //控制器工厂，需要自定义初始化Controller的就可以提供这个函数
	DisableMethod  []iface.Method      `json:"disable_method,omitempty"`  //取消掉的方法（API）
	EnableMethod   []iface.Method      `json:"enable_method,omitempty"`   //当暴露接口的时候，配置只允许那些接口
	BaseUrl        string              `json:"base_url,omitempty"`        //当前实体的通用路径
	Order          int                 `json:"order,omitempty"`           //数据依赖顺序。数值越小表示对其他实体的依赖越小，越优先处理数据
}

func (ec *EntityConfig) GetTable() *entity.Table {

	var (
		tb = &ec.Table
	)

	//这个情况一帮是对应的delta_data模块对应的EntityConfig
	//也就是EntityType和TableName都是动态的
	if tb.EntityType == nil && tb.TableName == "" {
		logger.Warnf("The EntityConfig's EntityType and TableName is empty which base url is [%s] ", ec.BaseUrl)
		tb.TableName = ec.Name
	}

	if !reflect2.IsNil(tb.EntityType) && len(tb.Fields) < 1 {
		tb.Fields = entity.GetTableFields(reflect.New(tb.EntityType))
	}
	if len(tb.EsExtraFields) < 1 && ctype.Bool(tb.EsEnable) && tb.EsEntityType != nil {
		cols := make([]string, 0, len(tb.Fields))
		for _, v := range tb.Fields {
			cols = append(cols, v.Name)
		}
		esFields := entity.GetTableFields(tb.EsEntityType)
		slices.DeleteFunc(esFields, func(field entity.TableField) bool {
			return slice.Contain(cols, field.Name)
		})
		//es的字段是包括entity的字段和指定的额外字段
		tb.EsExtraFields = esFields
	}
	return tb
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
	return toList
}
