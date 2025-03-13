package task

import (
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/eta/constants"
	jsoniter "github.com/json-iterator/go"
	"gorm.io/gorm/schema"
)

const (
	DataTaskValidating = "校验中"
	DataTaskProcessing = "处理中"
	DataTaskFinished   = "已完成"
	DataTaskFailed     = "失败"
)

type DataTask struct {
	entity.BaseEntity
	Tbl            string      `json:"tbl" gorm:"column:tbl;not null;comment:表名/模板编码"`
	Name           string      `json:"name" gorm:"column:name;not null;comment:名称"`
	Status         string      `json:"status" gorm:"column:status;not null;comment:状态"`
	FilePath       string      `json:"file_path" gorm:"column:file_address;comment:导出/导入文件路径"`
	FailFilePath   string      `json:"fail_file_path,omitempty" gorm:"column:fail_file_path;comment:失败文件路径"`
	FailCount      int64       `json:"fail_count" gorm:"column:fail_count;comment:失败条数"`
	FailReason     *ctype.Json `json:"fail_reason,omitempty" gorm:"column:fail_reason;comment:失败原因"`
	Params         *ctype.Json `json:"params,omitempty" gorm:"column:params;comment:参数"`
	Total          int64       `json:"total" gorm:"column:total;comment:总条数"`
	ProcessedCount int64       `json:"processed_count"` // 已处理条数
}

func (t DataTask) TableName(n schema.Namer) string {
	if n == nil {
		return constants.TableDataTask
	}
	return n.TableName(constants.TableDataTask)
}
func (t DataTask) ToString(containParams bool) string {
	tmp := DataTask{
		BaseEntity:   entity.BaseEntity{ID: t.ID},
		Tbl:          t.Tbl,
		Name:         t.Name,
		Status:       t.Status,
		FilePath:     t.FilePath,
		FailFilePath: t.FailFilePath,
		FailCount:    t.FailCount,
		FailReason:   t.FailReason,
		//Params:         t.Params,
		Total:          t.Total,
		ProcessedCount: t.ProcessedCount,
	}

	if containParams {
		tmp.Params = t.Params
	}
	toString, err := jsoniter.MarshalToString(tmp)
	if err != nil {
		return ""
	}
	return toString
}
