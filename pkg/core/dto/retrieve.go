package dto

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/expr"
	"github.com/davycun/eta/pkg/common/dorm/filter"
	"slices"
)

type OrderBy struct {
	Column string `json:"column"`
	Asc    bool   `json:"asc"`
}

type CommonParam struct {
}

type RetrieveParam struct {
	AutoCount    bool             `json:"auto_count"` //是否添加count的返回
	OnlyCount    bool             `json:"only_count"`
	GeoFormat    string           `json:"geo_format"` //支持wkb、ewkb、wkt、geojson
	GeoGCSType   string           `json:"geo_gcs_type"`
	PageSize     int              `json:"page_size"`
	PageNum      int              `json:"page_num"` //从1开始，不是从0开始
	OrderBy      []dorm.OrderBy   `json:"order_by"`
	ExtraColumns []expr.ExpColumn `json:"extra_columns"` //ExtraColumns表示一些额外的需要通过表达式来获取的字段

	//一些通用的
	Columns           []string        `json:"columns"`
	MustColumns       []string        `json:"must_columns"` //主要是在一些retrieve的callback中需要必须获取一些字段，比如签名验签
	Filters           []filter.Filter `json:"filters"`
	Extra             any             `json:"extra"`               //其他额外的请求参数由继承者可以自定义
	UseCurDeptAuth    bool            `json:"use_cur_dept_auth"`   //表示只是使用当前角色(部门)进行权限限制
	DisablePermFilter bool            `json:"disable_perm_filter"` //禁用权限过滤
	//*****当前从权限里面获取到了Recursive相关的filter就放在这个字段里面，后续Builder SQL的时候需要用到
	AuthFilters          []filter.Filter `json:"auth_filters"`
	AuthRecursiveFilters []filter.Filter `json:"auth_recursive_filters"`
	Auth2RoleFilters     []filter.Filter `json:"auth2role_filters"` //权限是通过auth2role表存储的关联数据

	//********针对树状实体********
	WithTree         bool            `json:"with_tree,omitempty"`         //返回结果是否是树形结构,这个只是对实现了ctype.TreeEntity接口的实体的接口有用
	WithParent       bool            `json:"with_parent,omitempty"`       //返回结果是否携带树状节点的父节点信息
	RecursiveFilters []filter.Filter `json:"recursive_filters,omitempty"` //查询树的顶点（如果IsUp是true就是底点）条件
	IsUp             bool            `json:"is_up,omitempty"`             //isUp 和WithParent
	LoadAll          bool            `json:"load_all"`                    //针对WithTree的情况，并且数据量小的时候才支持
	TreeDepth        int             `json:"tree_depth"`                  //树状结构深度

	SearchContent string `json:"search_content,omitempty"` // search 全文检索，暂时只在一些定制接口中支持

}

func (s RetrieveParam) GetOffset() int {
	return (s.GetPageNum() - 1) * s.GetPageSize()
}
func (s RetrieveParam) GetPageNum() int {
	if s.PageNum < 1 {
		s.PageNum = 1
	}
	return s.PageNum
}

func (s RetrieveParam) GetPageSize() int {
	if s.PageSize < 1 {
		s.PageSize = 10
	}
	return s.PageSize
}
func (s RetrieveParam) GetLimit() int {
	if s.PageSize < 1 {
		s.PageSize = 10
	}
	return s.PageSize
}

func (s RetrieveParam) ResolveOrderByString(tableName, defaultOrderBy string, addKeyWords bool) string {
	return dorm.ResolveOrderByString(s.OrderBy, tableName, defaultOrderBy, addKeyWords)
}
func (s RetrieveParam) ResolveEsOrderBy() []map[string]interface{} {
	return dorm.ResolveEsOrderBy(s.OrderBy...)
}

func (s RetrieveParam) ResolveLimitString() string {
	return fmt.Sprintf(` limit %d offset %d `, s.GetPageSize(), s.GetOffset())
}

func (s RetrieveParam) ResolveFilter(dbType dorm.DbType, addWhere bool, tableName string) string {
	if len(s.Filters) > 0 {
		where := filter.ResolveWhereTable(tableName, s.Filters, dbType)
		if addWhere {
			return " where " + where
		}
		return where
	}
	return ""
}
func (s RetrieveParam) Clone() RetrieveParam {
	rp := RetrieveParam{
		GeoFormat:         s.GeoFormat,
		GeoGCSType:        s.GeoGCSType,
		PageNum:           s.PageNum,
		PageSize:          s.PageSize,
		OrderBy:           slices.Clone(s.OrderBy),
		Filters:           slices.Clone(s.Filters),
		Columns:           slices.Clone(s.Columns),
		ExtraColumns:      slices.Clone(s.ExtraColumns),
		AutoCount:         s.AutoCount,
		OnlyCount:         s.OnlyCount,
		Extra:             s.Extra,
		UseCurDeptAuth:    s.UseCurDeptAuth,
		DisablePermFilter: s.DisablePermFilter,
		SearchContent:     s.SearchContent,
		WithTree:          s.WithTree,
		WithParent:        s.WithParent,
		RecursiveFilters:  slices.Clone(s.RecursiveFilters),
		IsUp:              s.IsUp,
		LoadAll:           s.LoadAll,
	}
	return rp
}

func GetExtra[T any](args *Param) *T {
	var (
		t T
	)
	switch args.Extra.(type) {
	case T:
		t = args.Extra.(T)
		return &t
	case *T:
		return args.Extra.(*T)
	}
	return &t
}

// NewParamWithExtra 示例：
//
//	controller.Publish(tableName, "/list", controller.ApiConfig{
//			GetParam: NewParamWithExtra[ExtraParam](),
//		})
//
// the param is &Param{RetrieveParam: RetrieveParam{Extra: &t}}
func NewParamWithExtra[T any]() func() any {
	return func() any {
		var t T
		return &Param{
			RetrieveParam: RetrieveParam{
				Extra: &t,
			},
		}
	}
}

// NewParamFunc 示例：
//
//	controller.Publish(tableName, "/list", controller.ApiConfig{
//			GetParam: NewParamFunc[MyParam](),
//		})
//
// the param is &MyParam
func NewParamFunc[T any]() func() any {
	return func() any {
		return new(T)
	}
}
