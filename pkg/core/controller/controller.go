package controller

import (
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/dorm/filter"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/core/service"
	"github.com/gin-gonic/gin"
	"io"
)

var (
	//ValidateFieldExcept = []string{"Data.ID", "Data.UpdatedAt", "Data.BaseEntity.ID", "Data.BaseEntity.UpdatedAt", "Data.BaseEdgeEntity.BaseEntity.ID", "Data.BaseEdgeEntity.BaseEntity.UpdatedAt"}
	ValidateFieldExcept = []string{"*.BaseEntity", "*.ID", "*.UpdatedAt"}
	//ValidateFieldPartial = []string{"Data.BaseEntity.ID", "Data.BaseEntity.UpdatedAt", "Data.BaseEdgeEntity.BaseEntity.ID", "Data.BaseEdgeEntity.BaseEntity.UpdatedAt"}
	ValidateFieldPartial = []string{"*.ID", "*.UpdatedAt"}
)

type DefaultController struct {
	NewService iface.NewService
}

func NewDefaultController(srvFactory iface.NewService) iface.Controller {
	if srvFactory == nil {
		srvFactory = service.NewDefaultService
	}
	h := &DefaultController{NewService: srvFactory}
	return h
}

func (handler *DefaultController) SetNewService(srv iface.NewService) {
	handler.NewService = srv
}
func (handler *DefaultController) GetNewService() iface.NewService {
	return handler.NewService
}
func (handler *DefaultController) Create(c *gin.Context) {
	var (
		param  = &dto.Param{}
		result = &dto.Result{}
	)
	srv := handler.GetService(c)
	err := bindModifyParam(srv, c, param, iface.MethodCreate)
	if err != nil {
		ProcessResult(c, result, err)
		return
	}
	err = srv.Create(param, result)
	ProcessResult(c, result, err)
}
func (handler *DefaultController) Update(c *gin.Context) {
	var (
		err    error
		param  = &dto.Param{}
		result = &dto.Result{}
	)
	srv := handler.GetService(c)
	err = bindModifyParam(srv, c, param, iface.MethodUpdate)
	if err != nil {
		ProcessResult(c, result, err)
		return
	}
	err = srv.Update(param, result)
	ProcessResult(c, result, err)
}
func (handler *DefaultController) UpdateByFilters(c *gin.Context) {
	var (
		err    error
		param  = &dto.Param{}
		result = &dto.Result{}
	)
	srv := handler.GetService(c)
	err = bindModifyParam(srv, c, param, iface.MethodUpdateByFilters)
	if err != nil {
		ProcessResult(c, result, err)
		return
	}
	err = srv.UpdateByFilters(param, result)
	ProcessResult(c, result, err)
}
func (handler *DefaultController) Delete(c *gin.Context) {

	var (
		err    error
		param  = &dto.Param{}
		result = &dto.Result{}
	)
	srv := handler.GetService(c)
	err = bindModifyParam(srv, c, param, iface.MethodDelete)
	if err != nil {
		ProcessResult(c, result, err)
		return
	}
	err = srv.Delete(param, result)
	ProcessResult(c, result, err)
}
func (handler *DefaultController) DeleteByFilters(c *gin.Context) {

	var (
		err    error
		param  = &dto.Param{}
		result = &dto.Result{}
	)
	srv := handler.GetService(c)
	err = bindModifyParam(srv, c, param, iface.MethodDeleteByFilters)
	if err != nil {
		ProcessResult(c, result, err)
		return
	}
	err = srv.DeleteByFilters(param, result)
	ProcessResult(c, result, err)
}
func (handler *DefaultController) Query(c *gin.Context) {
	var (
		err    error
		param  = &dto.Param{}
		result = &dto.Result{}
	)
	srv := handler.GetService(c)
	param.Extra = dto.DefaultParamExtra()
	err = BindBody(c, param)
	if err != nil && err != io.EOF {
		ProcessResult(c, result, err)
		return
	}
	dto.InitPage(param)
	err = srv.Query(param, result)
	result.PageSize = param.GetPageSize()
	result.PageNum = param.GetPageNum()
	ProcessResult(c, result, err)
}
func (handler *DefaultController) Detail(c *gin.Context) {
	var (
		id struct {
			ID string `json:"id,omitempty" uri:"id" binding:"required"`
		}
		err    error
		param  = &dto.Param{}
		result = &dto.Result{}
	)
	err = BindUri(c, &id)
	if err != nil {
		ProcessResult(c, result, err)
		return
	}
	err = BindBody(c, param)
	if err != nil && err != io.EOF {
		ProcessResult(c, result, err)
		return
	}

	srv := handler.GetService(c)
	param.Filters = append(param.Filters, filter.Filter{LogicalOperator: filter.And, Column: entity.IdDbName, Operator: filter.Eq, Value: id.ID})
	dto.InitPage(param)
	err = srv.Detail(param, result)
	ProcessResult(c, result, err)
}
func (handler *DefaultController) Count(c *gin.Context) {
	var (
		err    error
		param  = &dto.Param{}
		result = &dto.Result{}
	)

	err = BindBody(c, param)
	if err != nil && err != io.EOF {
		ProcessResult(c, result, err)
		return
	}
	srv := handler.GetService(c)
	dto.InitPage(param)
	err = srv.Count(param, result)
	ProcessResult(c, result, err)
}
func (handler *DefaultController) Aggregate(c *gin.Context) {
	var (
		err    error
		param  = &dto.Param{}
		result = &dto.Result{}
	)
	err = BindBody(c, param)
	if err != nil {
		ProcessResult(c, result, err)
		return
	}
	srv := handler.GetService(c)
	dto.InitPage(param)
	err = srv.Aggregate(param, result)
	result.PageSize = param.GetPageSize()
	result.PageNum = param.GetPageNum()
	ProcessResult(c, result, err)
}
func (handler *DefaultController) Partition(c *gin.Context) {
	var (
		err    error
		param  = &dto.Param{}
		result = &dto.Result{}
	)
	err = BindBody(c, param)
	if err != nil {
		ProcessResult(c, result, err)
		return
	}
	srv := handler.GetService(c)
	dto.InitPage(param)
	err = srv.Partition(param, result)
	result.PageSize = param.GetPageSize()
	result.PageNum = param.GetPageNum()
	ProcessResult(c, result, err)
}
func (handler *DefaultController) Import(c *gin.Context) {

}
func (handler *DefaultController) Export(c *gin.Context) {

}
func (handler *DefaultController) ProcessResult(c *gin.Context, data interface{}, err error) {
	ProcessResult(c, data, err)
}
func (handler *DefaultController) GetService(c *gin.Context) iface.Service {

	var (
		ct = ctx.GetContext(c)
	)
	if handler.NewService == nil {
		return service.NewDefaultService(ct, ct.GetContextGorm(), iface.GetContextEntityConfig(ct))
	}
	return handler.NewService(ct, ct.GetContextGorm(), iface.GetContextEntityConfig(ct))
}

func bindModifyParam(srv iface.Service, c *gin.Context, param *dto.Param, cdt iface.Method) error {
	var (
		err error
	)
	switch cdt {
	case iface.MethodCreate:
		param.Data = srv.NewEntitySlicePointer()
		err = BindBodyExcept(c, param, ValidateFieldExcept...)
	case iface.MethodUpdate, iface.MethodDelete:
		param.Data = srv.NewEntitySlicePointer()
		err = BindBodyPartial(c, param, ValidateFieldPartial...)
	case iface.MethodUpdateByFilters, iface.MethodDeleteByFilters:
		param.Data = srv.NewEntityPointer()
		err = BindBodyPartial(c, param, "filters")
	}
	return err
}
