package iface

import (
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/core/entity"
	"gorm.io/gorm"
)

type SrvOptions struct {
	DisableRetrieveWithES bool          //是否禁用 ES 检索
	UseParamAuth          bool          //默认是false，也就是需要权限，如果设置为true。那么就会根据参数（DisablePermFilter）决定是否需要权限
	Table                 *entity.Table //服务配置
	OriginDB              *gorm.DB
	Ctx                   *ctx.Context
}
type SrvOptionsFunc func(o *SrvOptions)

func NewSrvOptions(optionsFunc ...SrvOptionsFunc) SrvOptions {
	so := SrvOptions{}
	for _, fc := range optionsFunc {
		fc(&so)
	}
	return so
}

func (s *SrvOptions) SetUseParamAuth(b bool) *SrvOptions {
	s.UseParamAuth = b
	return s
}
func (s *SrvOptions) SetDisableRetrieveWithES(b bool) *SrvOptions {
	s.DisableRetrieveWithES = b
	return s
}
func (s *SrvOptions) SetTable(tb *entity.Table) {
	s.Table = tb
}
func (s *SrvOptions) GetTable() *entity.Table {
	if s.Table == nil {
		s.Table = entity.GetContextTable(s.GetContext())
	}
	return s.Table
}
func (s *SrvOptions) GetContext() *ctx.Context {
	return s.Ctx
}
func (s *SrvOptions) SetContext(c *ctx.Context) {
	s.Ctx = c
}
func (s *SrvOptions) SetDB(orm *gorm.DB) {
	s.OriginDB = orm
}
func (s *SrvOptions) GetDB() *gorm.DB {
	return s.OriginDB
}
func (s *SrvOptions) GetDbType() dorm.DbType {
	return dorm.GetDbType(s.GetDB())
}
func (s *SrvOptions) GetTableName() string {
	return s.GetTable().GetTableName()
}

// GetEsIndexName
// 返回es的索引名称，真正的索引是需要加上schema的前缀的
func (s *SrvOptions) GetEsIndexName() string {
	return s.GetTable().GetEsIndexName()
}
func (s *SrvOptions) NewEntityPointer() any {
	return s.GetTable().NewEntityPointer()
}
func (s *SrvOptions) NewRsDataPointer() any {
	return s.GetTable().NewRsDataPointer()
}
func (s *SrvOptions) NewEntitySlicePointer() any {
	return s.GetTable().NewEntitySlicePointer()
}
func (s *SrvOptions) NewRsDataSlicePointer() any {
	return s.GetTable().NewRsDataSlicePointer()
}
