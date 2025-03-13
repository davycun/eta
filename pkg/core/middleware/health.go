package middleware

import (
	"context"
	"database/sql"
	"errors"
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/duke-git/lancet/v2/maputil"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/duke-git/lancet/v2/strutil"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

const (
	QMiddleware = "middleware"
	DATABASE    = "database"
	REDIS       = "redis"
	ES          = "es"
)

func Health(c *gin.Context) {
	u := strings.ToLower(c.Request.RequestURI)
	if u == "/favicon.ico" {
		c.AbortWithStatus(404)
		return
	}
	if !strings.HasPrefix(u, "/health") {
		return
	}
	var (
		dbInfo  = make(map[string]any)
		rdsInfo = make(map[string]any)
		esInfo  = make(map[string]any)
		dbErr   error
		rdsErr  error
		esErr   error
	)

	msg := ctype.Map{
		"message":      "ok",
		"current_time": time.Now(),
	}

	qmw := c.Query(QMiddleware)
	mid := strutil.SplitAndTrim(qmw, ",")

	// check database
	if slice.ContainBy(mid, func(item string) bool { return strings.ToLower(item) == DATABASE }) {
		dbInfo, dbErr = check(DATABASE, checkDb)
	}
	// check redis
	if slice.ContainBy(mid, func(item string) bool { return strings.ToLower(item) == REDIS }) {
		rdsInfo, rdsErr = check(REDIS, checkRedis)
	}
	// check es
	if slice.ContainBy(mid, func(item string) bool { return strings.ToLower(item) == ES }) {
		esInfo, esErr = check(ES, checkEs)
	}

	result := maputil.Merge(msg, dbInfo, rdsInfo, esInfo)
	code := 200
	if !slice.Every([]any{dbErr, rdsErr, esErr}, func(index int, item any) bool { return item == nil }) {
		code = 500
	}
	c.AbortWithStatusJSON(code, result)
}

func check(name string, f func() (err error)) (info map[string]any, err error) {
	err = f()
	msg, health := connected(err)
	return gin.H{
		name: gin.H{
			"message": msg,
			"health":  health,
		},
	}, err
}

func connected(err error) (string, bool) {
	if err == nil {
		return "connected", true
	}
	return "disconnected", false
}

// check database
func checkDb() (err error) {
	var d *sql.DB
	return caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			d, err = global.GetLocalGorm().DB()
			return err
		}).Call(func(cl *caller.Caller) error {
		return d.Ping()
	}).Err
}

// check redis
func checkRedis() (err error) {
	cli := global.GetRedis()
	r := cli.Ping(context.Background())
	if r.Err() != nil {
		err = r.Err()
	}
	return
}

// check es
func checkEs() (err error) {
	es := global.GetES()
	if es == nil {
		return errors.New("获取的 ES 是 nil")
	}
	health, err := es.EsApi.Cat.Health()
	if err != nil {
		return err
	}
	if health == nil || health.StatusCode != http.StatusOK {
		return errors.New("获取的 ES 状态码不是 200")
	}
	return
}
