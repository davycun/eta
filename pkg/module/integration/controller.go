package integration

import (
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/core/controller"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/core/service"
	"github.com/davycun/eta/pkg/eta/ecf"
	"github.com/davycun/eta/pkg/module/setting"
	"github.com/davycun/eta/pkg/module/template"
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

	//首次json反解析，主要是先获取各实体查询相关的参数，比如entityCode等
	err = jsoniter.Unmarshal(body, param)
	if err != nil {
		controller.ProcessResult(c, result, err)
		return
	}

	//绑定实际的实体到对应的参数，主要是先获取entityCode并且设置Data到对应的实体切片
	for i, item := range param.Items {

		cx, newSrv, err1 := parseTable(ct, item.EntityCode)
		if err1 != nil {
			controller.ProcessResult(c, result, err)
			return
		}
		//这个要在parseTable之后，应该parseTable 之后不报错，代表返回的cx已经存储了EntityConfig
		ec := ecf.GetContextEntityConfig(cx)
		if item.Command == iface.MethodUpdateByFilters.String() || iface.MethodDeleteByFilters.String() == item.Command {
			param.Items[i].Param.Data = ec.NewEntityPointer()
		} else {
			param.Items[i].Param.Data = ec.NewEntitySlicePointer()
		}
		srvList = append(srvList, txService{NewSrv: newSrv, C: cx, EC: ec, Command: item.Command, EntityCode: item.EntityCode})
	}
	//二次json反解析，主要是为了获取重新获取Param.Data具体的对应的实体
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
		ts.Param = item.Param
		ts.Result = &dto.Result{}
	}

	err = transactionCall(ct, param, srvList, result)
	controller.ProcessResult(c, result, err)
}

func parseTable(c *ctx.Context, entityCode string) (*ctx.Context, iface.NewService, error) {
	var (
		err   error
		cx    = c.Clone()
		ec, b = iface.GetEntityConfigByKey(entityCode)
		temp  = template.Template{}
		appDb = c.GetAppGorm()
		srv   iface.NewService
	)

	if b {
		ecTb := ec.GetTable()
		if bcTb, ok := setting.GetTableConfig(appDb, ecTb.GetTableName()); ok {
			ecTb.Merge(&bcTb)
			ec.SetTable(ecTb)
		}
		srv = ec.NewService
		if ec.NewService == nil {
			srv = service.NewServiceFactory(ec)
		}
	} else {
		temp, err = template.LoadByCode(appDb, entityCode)
		if err != nil {
			return cx, nil, err
		}
		srv = service.NewDefaultService
		ec.SetTable(temp.GetTable())
		template.SetContextTemplate(cx, &temp)
	}
	ecf.SetContextEntityConfig(cx, &ec)
	return cx, srv, nil
}
