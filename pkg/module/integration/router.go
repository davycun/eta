package integration

import (
	"github.com/davycun/eta/pkg/common/global"
)

func Router() {
	group := global.GetGin().Group("/integration")
	group.POST("/transaction", Transaction) // 同一事务处理多个操作

}
