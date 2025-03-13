package global

import "github.com/davycun/eta/pkg/common/config"

// InitApplication 如果需要global包下的所有函数可用，需要先调用这个初始化方法
func InitApplication(cfg *config.Configuration) {
	globalApp = NewApplication(cfg)
}
