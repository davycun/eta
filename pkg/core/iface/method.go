package iface

const (
	CallbackBefore CallbackPosition = 1
	CallbackAfter  CallbackPosition = 2
)
const (
	CurdModify   CurdType = "modify"   //修改数据，包括新增、修改、删除
	CurdRetrieve CurdType = "retrieve" //查询数据，包括查询及统计
	CurdAll      CurdType = "all"
)
const (
	MethodCreate          Method = "create"
	MethodUpdate          Method = "update"
	MethodUpdateByFilters Method = "update_by_filters"
	MethodDelete          Method = "delete"
	MethodDeleteByFilters Method = "delete_by_filters"
	MethodQuery           Method = "query"
	MethodDetail          Method = "detail"
	MethodCount           Method = "count"
	MethodAggregate       Method = "aggregate"
	MethodPartition       Method = "partition"
	MethodExport          Method = "export"
	MethodImport          Method = "import"
	MethodAll             Method = "method_all"
	MethodList            Method = "list" //预留定制的方法
)

var (
	AllModifyMethod = []Method{MethodCreate, MethodUpdate, MethodUpdateByFilters, MethodDelete, MethodDeleteByFilters}
)

func GetAllModifyMethod() []Method {
	return AllModifyMethod
}

type Method string

func (m Method) String() string {
	return string(m)
}

type CallbackPosition int

type CurdType string
