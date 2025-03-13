package dto

import "github.com/davycun/eta/pkg/common/dorm/ctype"

type Param struct {
	RetrieveParam
	AggregateParam
	PartitionParam
	ModifyParam
}

type Result struct {
	Total        int64 `json:"total,omitempty"`
	PageSize     int   `json:"page_size,omitempty"`
	PageNum      int   `json:"page_num,omitempty"`
	Data         any   `json:"data,omitempty"`
	RowsAffected int64 `json:"rows_affected,omitempty"`
}

type ControllerResponse struct {
	Code        string  `json:"code,omitempty"`
	Success     bool    `json:"success,omitempty"`
	Message     string  `json:"message,omitempty"`
	FullMessage *string `json:"full_message,omitempty"`
	Result      any     `json:"result,omitempty"`
}

// DefaultParamExtra
// 默认的DefaultParamExtra的Extra属性类型
func DefaultParamExtra() any {
	return &ctype.Map{}
}

func InitPage(args *Param) {
	if args.PageSize < 1 {
		args.PageSize = 10
	}
	if args.PageNum < 1 {
		args.PageNum = 1
	}
}
