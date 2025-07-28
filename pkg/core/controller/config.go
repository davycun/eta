package controller

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/core/service"
	"github.com/gin-gonic/gin"
	"io"
)

const (
	BindTypeBody  = "body"
	BindTypeUri   = "uri"
	BindTypeQuery = "query"
)

// ApiConfig
// P代表绑定的参数，R代表响应的结果
type ApiConfig struct {
	Handler   Handler
	GetParam  func() any
	GetResult func() any
	Binding   func(c *gin.Context, param any) error
	Validate  func(c *gin.Context, param any) error
	BindType  string // uri_path,body,query
}
type Handler func(srv iface.Service, args any, rs any) error

// NewApi 构建接口的
// paramExtra 函数返回值需要时指针，表示dto.ModifyParam.Extra 应该接收的结构
// paramData 函数返回值需要是指针，表示dto.ModifyParam.Data应该接收的结构
// paramExtra默认是map[string]nil
// paramData 默认是对应的EntityService的GetEntitySlicePointer() 或者GetEntityPointer返回值
// P代表请求body绑定的结构体，R代表响应返回的结构体
func NewApi(tableName string, cfg ApiConfig) func(c *gin.Context) {
	return func(c *gin.Context) {
		var (
			err    error
			param  any
			result any
			ct     = ctx.GetContext(c)
		)

		if cfg.GetParam != nil {
			param = cfg.GetParam()
		} else {
			param = &dto.Param{}
		}
		if cfg.GetResult != nil {
			result = cfg.GetResult()
		} else {
			result = &dto.Result{}
		}

		ec, b := iface.GetEntityConfigByTableName(tableName)
		if !b {
			ProcessResult(c, nil, errs.NewServerError(fmt.Sprintf("can not found EntityConfig for %s", tableName)))
			return
		}

		if ec.NewService == nil {
			ec.NewService = service.NewServiceFactory(ec)
		}
		srv := ec.NewService(ct, ct.GetContextGorm(), &ec)

		if x, ok := param.(*dto.Param); ok {
			if x.Extra == nil {
				x.Extra = dto.DefaultParamExtra()
			}
			if x.Data == nil {
				x.Data = srv.NewEntitySlicePointer()
			}
		}

		if cfg.Binding != nil {
			err = cfg.Binding(c, param)
		} else {
			switch cfg.BindType {
			case BindTypeBody:
				err = BindBody(c, param)
			case BindTypeUri:
				err = BindUri(c, param)
			case BindTypeQuery:
				err = BindQuery(c, param)
			default:
				err = BindBody(c, param)
			}
		}

		if cfg.Validate != nil {
			err = cfg.Validate(c, param)
		}

		if err != nil && err != io.EOF {
			ProcessResult(c, result, err)
			return
		}
		//if x, ok := param.(*dto.Param); ok {
		//	if x.Extra != nil {
		//		err = ValidateStructFields(x.Extra, false)
		//		if err != nil {
		//			ProcessResult(c, result, err)
		//			return
		//		}
		//	}
		//}
		err = cfg.Handler(srv, param, result)
		ProcessResult(c, result, err)
	}
}
