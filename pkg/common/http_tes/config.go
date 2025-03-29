package http_tes

import (
	"github.com/davycun/eta/pkg/common/config"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/eta/middleware"
	"github.com/davycun/eta/pkg/eta/router"
	"github.com/davycun/eta/pkg/eta/validator"
	"github.com/gin-gonic/gin/binding"
	"os"
)

var (
	LoginToken  = constants.DefaultOpenApiFixedToken
	TransferKey = "8eadb267efd6e860"
)

func initServer() error {
	confFile := findFile("config_local.yml")
	destConfig := config.LoadConfig(confFile, nil)
	//binding.EnableDecoderUseNumber = true
	binding.EnableDecoderDisallowUnknownFields = true
	global.InitApplication(destConfig)
	middleware.InitMiddleware()
	router.InitRouter()
	validator.AddValidate()
	return nil
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
