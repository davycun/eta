package iface

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/entity"
	"reflect"
)

type ServiceConfig struct {
	ServiceType  reflect.Type            `json:"service_type,omitempty"` //如果NewService没有，那么就通过类型直接创建
	ResultType   map[Method]reflect.Type `json:"result_type,omitempty"`  //返回的数据类型
	NewService   NewService              `json:"new_service,omitempty"`  //服务工厂，需要自定义初始化Service就可以提供这个函数
	UseParamAuth bool                    //默认是false，也就是需要权限，如果设置为true。那么就会根据参数（DisablePermFilter）决定是否需要权限
}

func (s *ServiceConfig) SetUseParamAuth(b bool) {
	s.UseParamAuth = b
}

type ControllerConfig struct {
	BaseUrl        string        `json:"base_url,omitempty"`        //当前实体的通用路径
	DisableApi     bool          `json:"disable_api,omitempty"`     //取消暴露API
	ControllerType reflect.Type  `json:"controller_type,omitempty"` //如果NewController没有，那么就通过类型直接创建
	NewController  NewController `json:"new_controller,omitempty"`  //控制器工厂，需要自定义初始化Controller的就可以提供这个函数
	DisableMethod  []Method      `json:"disable_method,omitempty"`  //取消掉的方法（API）
	EnableMethod   []Method      `json:"enable_method,omitempty"`   //当暴露接口的时候，配置只允许那些接口
}

type EntityConfig struct {
	entity.Table
	ControllerConfig
	ServiceConfig
	Namespace string `json:"namespace,omitempty"` //区分不同的定制系统或者产品或者模块
	Name      string `json:"name,omitempty"`      //实体的唯一名字，在事务集合接口中用来唯一标志一个唯一的操作的对象
	Migrate   bool   `json:"migrate,omitempty"`   //是否需要进行migrate
	Order     int    `json:"order,omitempty"`     //数据依赖顺序。数值越小表示对其他实体的依赖越小，越优先处理数据
}

func (ec *EntityConfig) GetKey() string {
	return fmt.Sprintf("%s@%s@%s", ec.Name, ec.GetTableName(), ec.BaseUrl)
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

	//if !reflect2.IsNil(tb.EntityType) && len(tb.Fields) < 1 {
	//	tb.Fields = entity.GetTableFields(reflect.New(tb.EntityType))
	//}
	//if len(tb.EsFields) < 1 && ctype.Bool(tb.EsEnable) && tb.EsEntityType != nil {
	//	cols := make([]string, 0, len(tb.Fields))
	//	for _, v := range tb.Fields {
	//		cols = append(cols, v.Name)
	//	}
	//	esFields := entity.GetTableFields(tb.EsEntityType)
	//	slices.DeleteFunc(esFields, func(field entity.TableField) bool {
	//		return slice.Contain(cols, field.Name)
	//	})
	//	//es的字段是包括entity的字段和指定的额外字段
	//	tb.EsFields = esFields
	//}
	return tb
}
func (ec *EntityConfig) SetTable(tb *entity.Table) {
	ec.Table = *tb
}

func (ec *EntityConfig) NewResultPointer(method Method) any {
	rsType := ec.GetResultType(method)
	if rsType != nil {
		return reflect.New(rsType).Interface()
	}
	return nil
}

func (ec *EntityConfig) NewResultSlicePointer(method Method) any {
	rsType := ec.GetResultType(method)
	if rsType != nil {
		return reflect.New(reflect.SliceOf(rsType)).Interface()
	}
	return nil
}

func (ec *EntityConfig) GetResultType(method Method) reflect.Type {
	var (
		rsType reflect.Type
	)
	if ec.ResultType != nil {
		rsType = ec.ResultType[method]
	}
	if rsType == nil {
		rsType = ec.ResultType[MethodAll]
	}
	if rsType == nil {
		if ec.EsRetrieveEnabled() {
			rsType = ec.GetEsEntityType()
		} else {
			rsType = ec.GetEntityType()
		}
	}
	return rsType
}
