package iface

import (
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/entity"
	"gorm.io/gorm"
)

// ContextService 主要是可以设置和获取GIN 引擎
type ContextService interface {
	GetContext() *ctx.Context
	SetContext(c *ctx.Context)
}

// OrmService 获取和设置gorm.DB结构体
type OrmService interface {
	SetDB(orm *gorm.DB)
	GetDB() *gorm.DB
}

type EntityService interface {
	NewEntityPointer() any
	NewEntitySlicePointer() any
}
type RsDataService interface {
	NewRsDataPointer() any
	NewRsDataSlicePointer() any
}

type TableService interface {
	GetTableName() string
	GetTable() *entity.Table
	SetTable(tb *entity.Table)
}

type CreateService interface {
	Create(args *dto.Param, result *dto.Result) error
}
type UpdateService interface {
	Update(args *dto.Param, result *dto.Result) error
	UpdateByFilters(args *dto.Param, result *dto.Result) error
}
type DeleteService interface {
	Delete(args *dto.Param, result *dto.Result) error
	DeleteByFilters(args *dto.Param, result *dto.Result) error
}
type RetrieveService interface {
	Query(args *dto.Param, result *dto.Result) error
	Detail(args *dto.Param, result *dto.Result) error
	DetailById(id int64, result *dto.Result) error
}
type AggregateService interface {
	Count(args *dto.Param, result *dto.Result) error
	Aggregate(args *dto.Param, result *dto.Result) error
}
type PartitionService interface {
	Partition(args *dto.Param, result *dto.Result) error
}
type ExportService interface {
	Export(args *dto.Param, result *dto.Result) error
}
type ImportService interface {
	Import(args *dto.Param, result *dto.Result) error
}
type InitService interface {
	Init() error
}

type OptionService interface {
	Options(fs ...SrvOptionsFunc)
}

type Service interface {
	TableService
	ContextService
	OrmService
	CreateService
	DeleteService
	UpdateService
	RetrieveService
	EntityService
	RsDataService
	AggregateService
	PartitionService
	ExportService
	ImportService
	InitService
}
type NewService func(c *ctx.Context, db *gorm.DB, tb *entity.Table) Service

type ServiceFactory interface {
	NewService(c *ctx.Context, db *gorm.DB, tb *entity.Table) Service
}
