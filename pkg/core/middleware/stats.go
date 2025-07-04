package middleware

import (
	"bufio"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/duke-git/lancet/v2/maputil"
	"github.com/gin-gonic/gin"
	"os"
	"strings"
	"time"
)

func Stats(c *gin.Context) {
	u := strings.ToLower(c.Request.URL.Path)
	u = strings.TrimSpace(strings.ToLower(u))
	if u != "/stats" && u != "/stats/" {
		return
	}
	var (
		codeInfo = getCodeInfo()
	)
	systemInfo := ctype.Map{
		"message": "ok",
		//"version":      global.Version,
		"current_time": time.Now(),
		"commit_tag":   os.Getenv("COMMIT_TAG"),
		"commit_hash":  os.Getenv("COMMIT_HASH"),
		"target_arch":  os.Getenv("TARGET_ARCH"),
	}
	result := maputil.Merge(systemInfo, codeInfo)
	c.AbortWithStatusJSON(200, result)
}

func getCodeInfo() ctype.Map {
	var (
		rs      = make(ctype.Map)
		fl, err = os.Open(".env")
	)
	if err != nil {
		logger.Infof("加载环境变量文件出错 %s", err)
	}
	sc := bufio.NewScanner(fl)
	for sc.Scan() {
		text := sc.Text()
		if text == "\n" {
			//空行
			continue
		}
		env := strings.Split(text, "=")
		if len(env) == 2 {
			rs[strings.ToLower(strings.TrimSpace(env[0]))] = strings.TrimSpace(env[1])
		}
	}

	return rs
}
