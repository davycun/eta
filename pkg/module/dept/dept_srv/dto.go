package dept_srv

import (
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	jsoniter "github.com/json-iterator/go"
)

type ImportDeptItem struct {
	DeptName            string        `json:"dept_name,omitempty" binding:"required"` //部门
	ManageLevel         ctype.Integer `json:"manage_level,omitempty"`                 //管理层级
	ManageAreaId        []string      `json:"manage_area_id,omitempty"`               //管理范围id
	BuildingUsage       []int         `json:"building_usage,omitempty"`               //楼宇类别
	ID                  string        `json:"-"`                                      //部门id，处理过程中填充
	ParentId            string        `json:"-"`                                      //父部门id，处理过程中填充
	BuildingUsageReduce []int         `json:"-"`                                      //拆分后的楼宇类别
}
type ImportDeptParamExtra struct {
	NameSpace string `json:"name_space" binding:"required"`
}
type ImportDeptParam struct {
	ImportDeptParamExtra
	Items []ImportDeptItem `json:"items,omitempty" binding:"required,min=1"`
}

type ImportDeptFailedItem struct {
	DeptName     string   `json:"dept_name,omitempty"`
	FailedReason []string `json:"failed_reason,omitempty"`
}

type ImportDeptResult struct {
	TaskID string `json:"task_id,omitempty"`
}

type ImportDeptWsResult struct {
	TaskID       string                 `json:"task_id,omitempty"`
	Status       string                 `json:"status,omitempty"`
	FinishedStep []string               `json:"finished_step,omitempty"` //管理层级校验/楼宇类别校验/管理范围校验
	CreatedCount int64                  `json:"created_count,omitempty"`
	UpdatedCount int64                  `json:"updated_count,omitempty"`
	Failed       []ImportDeptFailedItem `json:"failed,omitempty"`
}

func (i *ImportDeptWsResult) ToString() string {
	toString, err := jsoniter.MarshalToString(i)
	if err != nil {
		return ""
	}
	return toString
}
