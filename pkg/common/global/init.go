package global

import (
	"github.com/davycun/eta/pkg/common/config"
	"time"
)

// InitApplication 如果需要global包下的所有函数可用，需要先调用这个初始化方法
func InitApplication(cfg *config.Configuration) error {
	initFixedZone()
	globalApp = NewApplication(cfg)
	return nil
}

func initFixedZone() {
	switch time.Local.String() {
	case "Asia/Shanghai":
		time.Local = time.FixedZone("CST", 8*3600)
	}
}
