package forward_srv

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/module/forward"
)

func InitModule() {
	global.GetGin().Any(fmt.Sprintf("/forward/:%s/*%s", forward.PathVendor, forward.PathParam), Forward)
}
