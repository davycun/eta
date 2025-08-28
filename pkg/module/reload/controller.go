package reload

import (
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/dsync"
	"github.com/davycun/eta/pkg/core/controller"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"io"
)

func NewReloadController(rdType RdType) gin.HandlerFunc {

	return func(c *gin.Context) {
		var (
			param     = &RdParamList{}
			result    = &RdResultList{}
			body, err = io.ReadAll(c.Request.Body)
		)

		//首次json反解析，主要是先获取各实体查询相关的参数，比如entityCode等
		err = jsoniter.Unmarshal(body, param)
		if err != nil {
			controller.ProcessResult(c, result, err)
			return
		}

		//绑定实际的实体到对应的参数，主要是先获取entityCode并且设置Data到对应的实体切片
		for _, item := range param.Items {
			item.Param.Extra = &dsync.SyncOption{}
		}
		//二次json反解析，主要是为了获取重新获取Param.Data具体的对应的实体
		err = jsoniter.Unmarshal(body, param)
		if err != nil {
			controller.ProcessResult(c, result, err)
			return
		}

		err = processReload(rdType, ctx.GetContext(c), param, result)

		controller.ProcessResult(c, result, err)
	}
}
