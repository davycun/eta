package iface

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/eta/constants"
	"gorm.io/gorm"
)

type (
	SrvOptions struct {
		EC       *EntityConfig
		OriginDB *gorm.DB
		Ctx      *ctx.Context
	}

	SrvOptionsFunc func(o *SrvOptions)
)

func NewSrvOptions(c *ctx.Context, db *gorm.DB, ec *EntityConfig, optionsFunc ...SrvOptionsFunc) SrvOptions {
	so := SrvOptions{}
	for _, fc := range optionsFunc {
		fc(&so)
	}
	return so
}

func NewSrvOptionsFromService(srv Service) SrvOptions {
	if srv == nil {
		return SrvOptions{}
	}

	return SrvOptions{
		EC:       srv.GetEntityConfig(),
		Ctx:      srv.GetContext(),
		OriginDB: srv.GetDB(),
	}
}

func (s *SrvOptions) UseParamAuth() bool {
	return s.GetTable().UseParamAuth()
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
	if s.EC == nil && s.Ctx != nil {
		//避免循环依赖，不用ecf包
		value, exists := s.Ctx.Get(constants.EntityConfigContextKey)
		if exists {
			s.EC = value.(*EntityConfig)
		}
	}
	return s.EC
}
func (s *SrvOptions) SetEntityConfig(ec *EntityConfig) {
	s.EC = ec
}

// GetEsIndexName
// 返回es的索引名称，真正的索引是需要加上schema的前缀的
func (s *SrvOptions) GetEsIndexName() string {
	return fmt.Sprintf("%s_%s", dorm.GetDbSchema(s.GetDB()), s.GetTableName())
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
