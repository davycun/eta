package dto

import (
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/filter"
)

type AggregateParam struct {
	AggregateColumns []dorm.AggregateColumn `json:"aggregate_columns,omitempty" binding:"dive"`
	GroupColumns     []string               `json:"group_columns" binding:"required"`
	Having           []filter.Having        `json:"having"`
}

// PartitionParam
// 如果有distinct，那么如果有orderBy，那么orderBy 左边排序必须也要出现distinct的字段
type PartitionParam struct {
	Distinct         []string               `json:"distinct"`
	PartitionColumns []dorm.PartitionColumn `json:"partition_columns"`
}
