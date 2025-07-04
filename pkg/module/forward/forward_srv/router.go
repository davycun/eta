package forward_srv

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/module/forward"
)

func Router() {
	global.GetGin().Any(fmt.Sprintf("/forward/:%s/*%s", forward.PathVendor, forward.PathParam), Forward)
}
