package reload

import "github.com/davycun/eta/pkg/core/dto"

type RdParamList struct {
	Concurrent bool      `json:"concurrent,omitempty"` //是否并发执行，默认是顺序执行
	Items      []RdParam `json:"items" binding:"required"`
}
type RdParam struct {
	TableName string     `json:"table_name" binding:"required"` // EntityConfig.Name, template.code
	Param     *dto.Param `json:"param" binding:"required"`
}

type RdResultList struct {
	Items []RdResult `json:"items" binding:"required"`
}
type RdResult struct {
	TableName string      `json:"table_name"`
	Result    *dto.Result `json:"result"`
}
