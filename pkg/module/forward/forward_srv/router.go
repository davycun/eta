package forward_srv

import (
	"github.com/davycun/eta/pkg/common/global"
)

func Router() {
	global.GetGin().Any("/forward/:vendor/*path", Forward)
}
