package controller

import (
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/dorm/filter"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/core/service"
	"github.com/davycun/eta/pkg/eta/ecf"
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
	srvList := handler.GetModifyService(c)
	err := bindModifyParam(c, param, iface.MethodCreate)
	if err != nil {
		ProcessResult(c, result, err)
		return
	}
	for _, s := range srvList {
		err = s.Create(param, result)
		if err != nil {
			break
		}
	}
	ProcessResult(c, result, err)
}
func (handler *DefaultController) Update(c *gin.Context) {
	var (
		err    error
		param  = &dto.Param{}
		result = &dto.Result{}
	)
	srvList := handler.GetModifyService(c)
	err = bindModifyParam(c, param, iface.MethodUpdate)
	if err != nil {
		ProcessResult(c, result, err)
		return
	}
	for _, s := range srvList {
		err = s.Update(param, result)
		if err != nil {
			break
		}
	}
	ProcessResult(c, result, err)
}
func (handler *DefaultController) UpdateByFilters(c *gin.Context) {
	var (
		err    error
		param  = &dto.Param{}
		result = &dto.Result{}
	)
	srvList := handler.GetModifyService(c)
	err = bindModifyParam(c, param, iface.MethodUpdateByFilters)
	if err != nil {
		ProcessResult(c, result, err)
		return
	}
	for _, s := range srvList {
		err = s.UpdateByFilters(param, result)
		if err != nil {
			break
		}
	}
	ProcessResult(c, result, err)
}
func (handler *DefaultController) Delete(c *gin.Context) {

	var (
		err    error
		param  = &dto.Param{}
		result = &dto.Result{}
	)
	srvList := handler.GetModifyService(c)
	err = bindModifyParam(c, param, iface.MethodDelete)
	if err != nil {
		ProcessResult(c, result, err)
		return
	}
	for _, s := range srvList {
		err = s.Delete(param, result)
		if err != nil {
			break
		}
	}
	ProcessResult(c, result, err)
}
func (handler *DefaultController) DeleteByFilters(c *gin.Context) {
	var (
		err    error
		param  = &dto.Param{}
		result = &dto.Result{}
	)

	srvList := handler.GetModifyService(c)
	err = bindModifyParam(c, param, iface.MethodDeleteByFilters)
	if err != nil {
		ProcessResult(c, result, err)
		return
	}
	for _, s := range srvList {
		err = s.DeleteByFilters(param, result)
		if err != nil {
			break
		}
	}
	ProcessResult(c, result, err)
}
func (handler *DefaultController) Query(c *gin.Context) {
	var (
		err    error
		param  = &dto.Param{}
		result = &dto.Result{}
	)
	srv := handler.GetRetrieveService(c)
	param.Extra = dto.DefaultParamExtra()
	err = BindBody(c, param)
	if err != nil && err != io.EOF {
		ProcessResult(c, result, err)
		return
	}
	clearAuth(param)
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
	clearAuth(param)
	err = BindBody(c, param)
	if err != nil && err != io.EOF {
		ProcessResult(c, result, err)
		return
	}

	srv := handler.GetRetrieveService(c)
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
	clearAuth(param)
	srv := handler.GetRetrieveService(c)
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
	srv := handler.GetRetrieveService(c)
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
	srv := handler.GetRetrieveService(c)
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
func (handler *DefaultController) GetModifyService(c *gin.Context) []iface.Service {

	var (
		ct      = ctx.GetContext(c)
		ec      = ecf.GetContextEntityConfig(ct)
		srvList = make([]iface.Service, 0, 2)
		ns      iface.NewService
	)

	if handler.NewService == nil {
		ns = service.NewDefaultService
	} else {
		ns = handler.NewService
	}
	//有多个service的原因是，有的实体表可能同时存在app库或者主库里面
	//srvList = append(srvList, ns(ct, ct.GetContextGorm(), ecf.GetContextEntityConfig(ct)))
	if ec.LocatedApp() {
		tmpCt := ct.Clone()
		tmpCt.SetContextGorm(ct.GetAppGorm())
		srvList = append(srvList, ns(tmpCt, tmpCt.GetContextGorm(), ec))
	}
	if ec.LocatedLocal() {
		tmpCt := ct.Clone()
		tmpCt.SetContextGorm(global.GetLocalGorm())
		srvList = append(srvList, ns(tmpCt, tmpCt.GetContextGorm(), ec))
	}
	return srvList
}
func (handler *DefaultController) GetRetrieveService(c *gin.Context) iface.Service {
	var (
		ct = ctx.GetContext(c)
		ec = ecf.GetContextEntityConfig(ct)
		ns iface.NewService
	)
	if handler.NewService == nil {
		ns = service.NewDefaultService
	} else {
		ns = handler.NewService
	}
	return ns(ct, ct.GetContextGorm(), ec)
}

func bindModifyParam(c *gin.Context, param *dto.Param, cdt iface.Method) error {
	var (
		err error
		ec  = ecf.GetContextEntityConfig(ctx.GetContext(c))
	)
	switch cdt {
	case iface.MethodCreate:
		param.Data = ec.NewEntitySlicePointer()
		err = BindBodyExcept(c, param, ValidateFieldExcept...)
	case iface.MethodUpdate, iface.MethodDelete:
		param.Data = ec.NewEntitySlicePointer()
		err = BindBodyPartial(c, param, ValidateFieldPartial...)
	case iface.MethodUpdateByFilters, iface.MethodDeleteByFilters:
		param.Data = ec.NewEntityPointer()
		err = BindBodyPartial(c, param, "filters")
	}
	clearAuth(param)
	return err
}

// 在controller层清除权限相关filters，避免客户端传入影响真实权限控制
func clearAuth(args *dto.Param) {
	clear(args.AuthFilters)
	clear(args.Auth2RoleFilters)
	clear(args.AuthRecursiveFilters)
}
