package entity

import (
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/eta/constants"
)

const (
	HistoryInsert Operator = 1
	HistoryUpdate Operator = 2
	HistoryDelete Operator = 3
)

type Operator int

func (o Operator) String() string {
	switch o {
	case HistoryInsert:
		return "insert"
	case HistoryUpdate:
		return "update"
	case HistoryDelete:
		return "delete"
	default:
		return "unknown"
	}
}
func TableName(originTableName string) string {
	return originTableName + constants.TableHistorySubFix
}

type History struct {
	ID        string           `json:"id,omitempty" gorm:"type:varchar(255);column:id"`
	CreatedAt *ctype.LocalTime `json:"created_at,omitempty" gorm:"<-:create;comment:创建时间;not null"`
	OpType    Operator         `json:"op_type,omitempty" gorm:"column:op_type;comment:操作类型"`
	OptUserId string           `json:"opt_user_id,omitempty" gorm:"column:opt_user_id;comment:操作人ID"`
	OptDeptId string           `json:"opt_dept_id,omitempty" gorm:"column:opt_dept_id;comment:操作人当前部门ID"`
}
