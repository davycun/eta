package http_tes

import (
	"github.com/davycun/eta/cmd/server"
	"github.com/davycun/eta/pkg/common/config"
	"github.com/davycun/eta/pkg/eta/constants"
	"os"
)

var (
	LoginToken  = constants.DefaultOpenApiFixedToken
	TransferKey = "8eadb267efd6e860"
)

func initServer() error {
	confFile := findFile("config_local.yml")
	destConfig := config.LoadConfig(confFile, nil)
	return server.CallSpecialLifeCycleWithConfig(destConfig, server.InitConfig,
		server.InitPlugin, server.InitApplication,
		server.InitMiddleware, server.InitValidator,
		server.InitData, server.InitEntityConfigRouter,
		server.InitModules, server.InitMigrator, server.Migrate)
}

func findFile(filename string) string {
	// 读取配置文件, 解决跑测试的时候找不到配置文件的问题，最多往上找5层目录
	for i := 0; i < 20; i++ {
		if _, err := os.Stat(filename); err == nil {
			return filename
		} else {
			filename = "../" + filename
		}
	}
	return ""
}
