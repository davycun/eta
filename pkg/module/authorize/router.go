package authorize

import (
	"github.com/davycun/eta/pkg/common/global"
)

func Router() {
	group := global.GetGin().Group("/authorize")
	group.Any("/check", Authorization)
}
