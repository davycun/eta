package iface

import (
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/core/entity"
	"gorm.io/gorm"
)

type SrvOptions struct {
	EC                    *EntityConfig
	OriginDB              *gorm.DB
	Ctx                   *ctx.Context
	DisableRetrieveWithES bool //是否禁用 ES 检索
	UseParamAuth          bool //默认是false，也就是需要权限，如果设置为true。那么就会根据参数（DisablePermFilter）决定是否需要权限
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
	s.GetEntityConfig().SetUseParamAuth(b)
	return s
}
func (s *SrvOptions) SetDisableRetrieveWithES(b bool) *SrvOptions {
	s.GetEntityConfig().SetDisableRetrieveWithES(b)
	return s
}
func (s *SrvOptions) GetTable() *entity.Table {
	return s.GetEntityConfig().GetTable()
}
func (s *SrvOptions) SetTable(tb *entity.Table) {
	s.GetEntityConfig().SetTable(tb)
}
func (s *SrvOptions) GetTableName() string {
	return s.GetTable().GetTableName()
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

func (s *SrvOptions) GetEntityConfig() *EntityConfig {
	if s.EC == nil {
		s.EC = GetContextEntityConfig(s.GetContext())
	}
	return s.EC
}
func (s *SrvOptions) SetEntityConfig(ec *EntityConfig) {
	s.EC = ec
}

// GetEsIndexName
// 返回es的索引名称，真正的索引是需要加上schema的前缀的
func (s *SrvOptions) GetEsIndexName() string {
	return s.GetEntityConfig().GetEsIndexName()
}
func (s *SrvOptions) NewEntityPointer() any {
	return s.GetEntityConfig().NewEntityPointer()
}
func (s *SrvOptions) NewEntitySlicePointer() any {
	return s.GetEntityConfig().NewEntitySlicePointer()
}
func (s *SrvOptions) NewResultPointer(method Method) any {
	return s.GetEntityConfig().NewResultPointer(method)
}
func (s *SrvOptions) NewResultSlicePointer(method Method) any {
	return s.GetEntityConfig().NewResultSlicePointer(method)
}

func (s *SrvOptions) Merge(src SrvOptions) {
	s.GetEntityConfig().Merge(src.GetTable())
}
