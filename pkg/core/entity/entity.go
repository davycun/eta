package entity

import (
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/utils"
	"strings"
)

const (
	IdDbName              = "id"
	CreatedAtDbName       = "created_at"
	UpdatedAtDbName       = "updated_at"
	CreatorIdDbName       = "creator_id"
	UpdaterIdDbName       = "updater_id"
	CreatorDeptIdDbName   = "creator_dept_id"
	UpdaterDeptIdDbName   = "updater_dept_id"
	FieldUpdaterDbName    = "field_updater"
	FieldUpdaterIdsDbName = "field_updater_ids"
	ExtraDbName           = "extra"
	EtlExtraDbName        = "etl_extra"
	FromIdDbName          = "from_id"
	ToIdDbName            = "to_id"
	RemarkDbName          = "remark"
	RaContentDbName       = "ra_content"
	RankDbName            = "rank"

	IdFieldName              = "ID"
	CreatedAtFieldName       = "CreatedAt"
	UpdatedAtFieldName       = "UpdatedAt"
	CreatorIdFieldName       = "CreatorId"
	UpdaterIdFieldName       = "UpdaterId"
	CreatorDeptIdFieldName   = "CreatorDeptId"
	UpdaterDeptIdFieldName   = "UpdaterDeptId"
	FieldUpdaterFieldName    = "FieldUpdater"
	FieldUpdaterIdsFieldName = "FieldUpdaterIds"
	ExtraFieldName           = "Extra"
	EtlExtraFieldName        = "EtlExtra"
	FromIdFieldName          = "FromId"
	ToIdFieldName            = "ToId"
	RemarkFieldName          = "Remark"
	RaContentFieldName       = "RaContent"
	RankFieldName            = "Rank"
)

var (
	DefaultVertexColumns  = []string{IdDbName, CreatedAtDbName, UpdatedAtDbName, CreatorIdDbName, UpdaterIdDbName, CreatorDeptIdDbName, UpdaterDeptIdDbName, FieldUpdaterDbName, FieldUpdaterIdsDbName, ExtraDbName, RemarkDbName}
	DefaultEdgeColumns    = append(DefaultVertexColumns, FromIdDbName, ToIdDbName)
	DefaultHistoryColumns = []string{IdDbName, CreatedAtDbName, "op_type"}
)

// MetaInfo
// 通过这个字段可以控制一些trigger是否需要执行
// 比如控制一些更新冗余字典的操作就不需要记录历史记录，ExecHistory设置为true，表示当前更新不需要记录历史记录，记得要在trigger上把他改为false
type MetaInfo struct {
	dorm.JsonType
	ExecHistory bool `json:"exec_history,omitempty"`
	ExecUpdater bool `json:"exec_updater,omitempty"`
}

// BaseEntity TODO 可以定义接口，然后来处理CreateUser、UpdateUser等
type BaseEntity struct {
	ID              string             `json:"id,omitempty" gorm:"type:varchar(255);column:id;<-:create;comment:ID" redis:"id"  binding:"required" es:"type:keyword"`
	CreatedAt       *ctype.LocalTime   `json:"created_at,omitempty" gorm:"<-:create;comment:创建时间;not null" redis:"created_at"  binding:"ignore" es:"type:date"`
	UpdatedAt       int64              `json:"updated_at,omitempty" gorm:"autoUpdateTime:milli;comment:更新时间戳;not null"  redis:"updated_at"  binding:"required" es:"type:long"`
	CreatorId       string             `json:"creator_id,omitempty" gorm:"type:varchar(255);column:creator_id;<-:create;comment:创建人;not null"  redis:"creator_id" es:"type:keyword"`
	UpdaterId       string             `json:"updater_id,omitempty" gorm:"type:varchar(255);column:updater_id;comment:最后更新人;not null"  redis:"updater_id" es:"type:keyword"`
	CreatorDeptId   string             `json:"creator_dept_id,omitempty" gorm:"type:varchar(255);column:creator_dept_id;<-:create;comment:创建的时候创建人所属部门;not null"  es:"type:keyword"`
	UpdaterDeptId   string             `json:"updater_dept_id,omitempty" gorm:"type:varchar(255);column:updater_dept_id;comment:更新的时候更新人所属部门;not null"  es:"type:keyword"`
	FieldUpdater    *ctype.Json        `json:"field_updater,omitempty" gorm:"column:field_updater;serializer:json;comment:实体中每个字段更新者"   binding:"ignore" es:"type:object"` //key是字段名称，value是更新用户的id
	FieldUpdaterIds *ctype.StringArray `json:"field_updater_ids,omitempty" gorm:"column:field_updater_ids;comment:存储每个字段更新人的ID"   binding:"ignore" es:"type:keyword"`
	Remark          string             `json:"remark,omitempty" gorm:"column:remark;comment:备注"   binding:"ignore" es:"type:text"`
	Extra           *ctype.Json        `json:"extra,omitempty" gorm:"column:extra;serializer:json;comment:扩展字段"  es:"type:object"`
	EtlExtra        *ctype.Json        `json:"etl_extra,omitempty" gorm:"column:etl_extra;serializer:json;comment:数据工程扩展字段" es:"type:object"`
	RaContent       *ctype.Text        `json:"ra_content,omitempty" gorm:"column:ra_content;comment:存储当前数据的描述" es:"type:text"`
	MetaInfo        MetaInfo           `json:"meta_info,omitempty" gorm:"column:meta_info;serializer:json;comment:扩展字段" es:"ignore"`
	Deleted         bool               `json:"deleted,omitempty" gorm:"-:all" es:"ignore"` //已经有history表，没必要存在软删除，否则如果作为筛选条件的话，影响性能
	CreatorName     string             `json:"creator_name,omitempty" gorm:"-:all" es:"type:keyword"`
	UpdaterName     string             `json:"updater_name,omitempty" gorm:"-:all" es:"type:keyword"`
	CreatorDeptName string             `json:"creator_dept_name,omitempty" gorm:"-:all" es:"type:keyword"`
	UpdaterDeptName string             `json:"updater_dept_name,omitempty" gorm:"-:all" es:"type:keyword"`
}

func (b BaseEntity) String() string {
	ss := make([]string, 0, 2)
	if b.ID != "" {
		ss = append(ss, b.ID)
	}
	ss = utils.Merge(ss, b.ID, b.CreatorId, b.UpdaterId, b.CreatorDeptId, b.UpdaterDeptId, b.Remark)
	return strings.Join(ss, " ")
}

func (b BaseEntity) DefaultColumns() []string {
	return []string{"*"}
}
func (b BaseEntity) MustColumns() []string {
	return []string{IdDbName, UpdatedAtDbName}
}
func (b BaseEntity) EmbeddedPrefix() string {
	return "emb_"
}

type BaseEdgeEntity struct {
	BaseEntity
	//TODO nebula property
	FromId string `json:"from_id,omitempty" redis:"from" gorm:"type:varchar(255);index;column:from_id;<-:create;not null;comment:关联实体起点ID" binding:"required" es:"type:keyword"`
	ToId   string `json:"to_id,omitempty" redis:"to" gorm:"type:varchar(255);index;column:to_id;<-create;not null;comment:关联实体终点ID"  binding:"required" es:"type:keyword"`
	Rank   int    `json:"rank,omitempty" gorm:"-:all" es:"ignore"`
}

func (b BaseEdgeEntity) DefaultColumns() []string {
	return []string{"*"}
}
func (b BaseEdgeEntity) MustColumns() []string {
	return []string{IdDbName, UpdatedAtDbName, FromIdDbName, ToIdDbName}
}
