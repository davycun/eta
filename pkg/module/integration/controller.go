package integration

import (
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/core/controller"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/core/service"
	"github.com/davycun/eta/pkg/core/service/ecf"
	"github.com/davycun/eta/pkg/module/data"
	"github.com/davycun/eta/pkg/module/data/template"
	"github.com/davycun/eta/pkg/module/setting"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"io"
)

// Transaction
// 同一事务处理多个操作
func Transaction(c *gin.Context) {
	var (
		param     = &CommandParam{}
		result    = &CommandResult{}
		ct        = ctx.GetContext(c)
		body, err = io.ReadAll(c.Request.Body)
		srvList   = make([]txService, 0, 2)
	)

	err = jsoniter.Unmarshal(body, param)
	if err != nil {
		controller.ProcessResult(c, result, err)
		return
	}

	//绑定实际的实体到对应的参数，主要是先获取entityCode并且设置Data到对应的实体切片
	for i, item := range param.Items {
		var (
			cx     = ct.Clone()
			newSrv iface.NewService
			tb     *entity.Table
		)
		cx, newSrv, err = parseTable(ct, item.EntityCode)
		if err != nil {
			controller.ProcessResult(c, result, err)
			return
		}
		tb = entity.GetContextTable(cx)
		if item.Command == iface.MethodUpdateByFilters.String() || iface.MethodDeleteByFilters.String() == item.Command {
			param.Items[i].Param.Data = tb.NewEntityPointer()
		} else {
			param.Items[i].Param.Data = tb.NewEntitySlicePointer()
		}
		srvList = append(srvList, txService{NewSrv: newSrv, C: cx, TB: tb, Command: item.Command})
	}
	//二次绑定，主要是重新制定dto.Param.Data为具体的实体类型
	err = jsoniter.Unmarshal(body, param)
	if err != nil {
		controller.ProcessResult(c, result, err)
		return
	}
	for i, item := range param.Items {
		switch item.Command {
		case iface.MethodCreate.String():
			err = controller.ValidateStructFields(item, true, controller.ValidateFieldExcept...)
		case iface.MethodUpdate.String(), iface.MethodDelete.String():
			err = controller.ValidateStructFields(item, false, controller.ValidateFieldPartial...)
		case iface.MethodUpdateByFilters.String(), iface.MethodDeleteByFilters.String():
			err = controller.ValidateStructFields(item, false, "filters")
		}
		if err != nil {
			controller.ProcessResult(c, result, err)
			return
		}

		ts := &srvList[i]
		ts.Param = &item.Param
		ts.Result = &dto.Result{}
	}

	err = TransactionCall(ct, srvList, result)
	controller.ProcessResult(c, result, err)
}

func parseTable(c *ctx.Context, entityCode string) (*ctx.Context, iface.NewService, error) {
	var (
		err   error
		cx    = c.Clone()
		ec, b = ecf.GetEntityConfigByName(entityCode)
		temp  = template.Template{}
		appDb = c.GetAppGorm()
		srv   iface.NewService
	)

	if b {
		ecTb := ec.GetTable()
		bcTb, _ := setting.GetTableConfig(appDb, ecTb.GetTableName())
		bcTb.Merge(ecTb)
		entity.SetContextTable(cx, &bcTb)
		srv = ec.NewService
		if ec.NewService == nil {
			srv = service.NewDefaultService
		}
	} else {
		temp, err = template.LoadByCode(appDb, entityCode)
		if err != nil {
			return cx, nil, err
		}
		tb := temp.GetTable()
		srv = service.NewDefaultService
		entity.SetContextTable(cx, tb)
		data.SetContextTemplate(cx, &temp)
	}
	return cx, srv, nil
}
