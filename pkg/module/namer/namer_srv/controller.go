package namer_srv

import (
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/core/controller"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/module/namer"
	"github.com/gin-gonic/gin"
)

func HandlerIdName(c *gin.Context) {
	var (
		err    error
		param  namer.IdNameParam
		result dto.Result
		mnMap  map[string]namer.IdName
		ct     = ctx.GetContext(c)
		dt     = make([]namer.IdName, 0, 1000)
	)

	err = controller.BindBody(c, &param)
	if err != nil {
		// ignore. 允许不传参数
	}

	if param.Ids != nil && len(param.Ids) > 0 {
		mnMap, err = namer.LoadByIds(ct, param.Ids)
	} else {
		// Ids 为空，查询全部
		mnMap, err = namer.LoadAllIdName(ct)
	}

	if err != nil {
		controller.ProcessResult(c, nil, err)
		return
	}

	ns := param.Namespace
	for _, v := range mnMap {
		if ns != "" {
			if v.Namespace == ns {
				dt = append(dt, v)
			}
		} else {
			dt = append(dt, v)
		}
	}

	result.Data = dt
	controller.ProcessResult(c, &result, nil)
	return
}
