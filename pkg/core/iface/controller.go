package iface

import (
	"github.com/gin-gonic/gin"
)

type CreateController interface {
	Create(c *gin.Context)
}
type UpdaterController interface {
	Update(c *gin.Context)
	UpdateByFilters(c *gin.Context)
}
type DeleteController interface {
	Delete(c *gin.Context)
	DeleteByFilters(c *gin.Context)
}
type RetrieveController interface {
	Query(c *gin.Context)
	Detail(c *gin.Context)
}
type StatsController interface {
	Count(c *gin.Context)
	Aggregate(c *gin.Context)
	Partition(c *gin.Context)
}
type ExportController interface {
	Export(c *gin.Context)
}
type ImportController interface {
	Import(c *gin.Context)
}

type ServiceController interface {
	SetNewService(srv NewService)
	GetNewService() NewService
}

type Controller interface {
	CreateController
	UpdaterController
	DeleteController
	RetrieveController
	StatsController
	ImportController
	ExportController
	ServiceController
}

type NewController func(srv NewService) Controller

type ControllerFactory interface {
	NewController(srv NewService) Controller
}
