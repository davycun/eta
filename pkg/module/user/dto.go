package user

import (
	"github.com/davycun/eta/pkg/common/dorm/filter"
	"github.com/davycun/eta/pkg/module/authorize/role"
	"github.com/davycun/eta/pkg/module/dept"
	jsoniter "github.com/json-iterator/go"
)

type ListParam struct {
	WithDept         bool            `json:"with_dept"`
	WithRole         bool            `json:"with_role"`
	User2DeptFilters []filter.Filter `json:"user2dept_filters"`
	DeptFilters      []filter.Filter `json:"dept_filters"`
	User2RoleFilters []filter.Filter `json:"user2role_filters"`
}

type ListResult struct {
	User
	User2Dept []dept.RelationDept `json:"user2dept,omitempty" gorm:"-:all"`
	User2Role []role.RelationRole `json:"user2role,omitempty" gorm:"-:all"`
}

type SendSmsCodeParam struct {
	Phone string `json:"phone" binding:"required,mobile"`
	AppId string `json:"app_id" binding:"required"`
}

// === 批量用户导入 ===

type ImportUserDept struct {
	ID   string `json:"id" binding:"required"`
	Post string `json:"post"`
}
type ImportUserRole struct {
	ID string `json:"id" binding:"required"`
}
type ImportUserItem struct {
	Name      string           `json:"name"  binding:"required"`
	Account   string           `json:"account"  binding:"required"`
	Password  string           `json:"password" binding:"required"`
	Phone     string           `json:"phone"`
	User2Dept []ImportUserDept `json:"user2dept"`
	User2Role []ImportUserRole `json:"user2role"`
}
type ImportUserParam struct {
	Items []ImportUserItem `json:"items" binding:"required,min=1"`
}
type ImportUserResult struct {
	TaskID string `json:"task_id,omitempty"`
}
type ImportUserFailedItem struct {
	Account      string   `json:"account,omitempty"`
	FailedReason []string `json:"failed_reason,omitempty"`
}
type ImportUserWsResult struct {
	TaskID       string                 `json:"task_id,omitempty"`
	Status       string                 `json:"status,omitempty"`
	FinishedStep []string               `json:"finished_step,omitempty"` // 登录名校验/手机号校验/部门信息校验/用户角色校验
	CreatedCount int64                  `json:"created_count,omitempty"`
	Failed       []ImportUserFailedItem `json:"failed,omitempty"`
}

func (i *ImportUserWsResult) ToString() string {
	toString, err := jsoniter.MarshalToString(i)
	if err != nil {
		return ""
	}
	return toString
}

// ModifyPasswordParam 修改密码
type ModifyPasswordParam struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
}

// ModifyPhoneParam 修改手机号
type ModifyPhoneParam struct {
	NewPhone string `json:"new_phone" binding:"required"`
	Code     string `json:"code" binding:"required"`
}

// ResetPasswordParam 重置密码
type ResetPasswordParam struct {
	UserIdList  []string `json:"user_id_list" binding:"required,max=1,min=1"`
	NewPassword string   `json:"new_password" binding:"required"`
}

type IdNameParam struct {
	Namespace string   `json:"namespace"`
	Ids       []string `json:"ids"`
}

type SendWsMessageParam struct {
	SendWsMessage bool `json:"send_ws_message"`
}
