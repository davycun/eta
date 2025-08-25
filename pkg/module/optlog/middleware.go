package optlog

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/run"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/controller"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/service"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/setting"
	"github.com/davycun/eta/pkg/module/user"
	"github.com/duke-git/lancet/v2/strutil"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	CreateData    = "新增数据"
	UpdateData    = "更新数据"
	DeleteData    = "删除数据"
	QueryData     = "查询数据"
	AggregateData = "统计数据"
	ImportData    = "导入数据"
	ExportData    = "导出数据"
	Login         = "登录"
	Logout        = "退出"
)

func Log(c *gin.Context) {
	var (
		uri    = c.Request.RequestURI
		method = c.Request.Method
		lg     = &OptLog{}
		ct     = ctx.GetContext(c).Clone()
	)
	if setting.IsIgnoreLogUri(nil, method, uri) {
		return
	}

	lg.OptTime = ctype.NewLocalTime(time.Now())
	c.Next()

	if c.Writer.Status() == http.StatusNotFound {
		return
	}

	lg.ReqId = ct.GetRequestId()
	lg.ReqUri = uri
	lg.ClientIp = c.ClientIP()
	lg.ClientType, lg.ClientTrigger = getOptClientType(c, lg)
	lg.OptUserId = ct.GetContextUserId()
	lg.OptDeptId = ct.GetContextCurrentDeptId()
	lg.OptType, lg.OptTarget, lg.OptContent = getOpt(c, lg)

	lg.Latency = time.Now().Sub(lg.OptTime.Data).Milliseconds() //响应时长
	lg.RsStatus = fmt.Sprintf("%d", c.Writer.Status())
	if c.GetHeader(constants.HeaderCryptSymmetryAlgorithm) != "" {
		lg.RsRemark = "返回结果加密"
	} else {
		lg.RsRemark = "返回结果非加密"
	}

	db := ct.GetAppGorm()
	if db == nil {
		return
	}
	//因为后面是异步调用，有可能CurrentContext已经被回收，所以先调用BeforeCreate，提前补全字段
	err := entity.BeforeCreate(&lg.BaseEntity, ct)
	if err != nil {
		logger.Errorf("日志BeforeCreate函数出错%s", err)
	}

	run.Go(func() {
		dt := []OptLog{*lg}
		ct.SetContextGorm(db)
		err = service.NewSrvWrapper(constants.TableOperateLog, ct, db).SetData(dt).Create()
		if err != nil {
			logger.Errorf("新增日志出错%s", err)
		}
	})
}

func getOptClientType(c *gin.Context, lg *OptLog) (clientType, clientTrigger string) {
	clientType = getHeaderDecode(c, constants.HeaderOptClientType)
	clientTrigger = getHeaderDecode(c, constants.HeaderOptClientTrigger)
	if clientType == "" {
		agent := strings.ToLower(c.Request.UserAgent())
		if strings.Contains(agent, "android") || strings.Contains(agent, "iphone") {
			clientType = "移动端"
		} else {
			clientType = "PC端"
		}
	}
	return
}

func getOpt(c *gin.Context, lg *OptLog) (optType, optTarget, optContent string) {

	var (
		uri = utils.GetUrlPath(c)
	)
	optType = getHeaderDecode(c, constants.HeaderOptType)
	optContent = getHeaderDecode(c, constants.HeaderOptContent)
	optTarget = getHeaderDecode(c, constants.HeaderOptTarget)

	if optContent != "" {
		lg.OptCategory = ManualCollect
	} else {
		lg.OptCategory = AutoCollect
	}

	if optTarget != "" && optContent != "" && optType != "" {
		return
	}

	s1, s2, s3 := getOptFromUrl(uri, ctx.GetContext(c))
	if optType == "" {
		optType = s1
	}
	if optContent == "" {
		optContent = s2
	}
	if optTarget == "" {
		optTarget = s3
	}
	return
}

func getOptFromUrl(uri string, c *ctx.Context) (optType, optContent, optTarget string) {
	var (
		tmp        = strings.Split(uri, "/")
		optTypeTmp = ""
		us, _      = user.GetContextUser(c)
	)

	ui, b := strings.CutSuffix(uri, "?")
	if b {
		uri = ui
	}

	switch len(tmp) {
	case 3:
		optTarget = tmp[1]
		optTypeTmp = tmp[2]
	case 4, 5:
		if strings.Contains(uri, "/detail") && strings.Contains(uri, "/data") {
			optTarget = tmp[1]
			optTypeTmp = tmp[2]
		} else {
			optTarget = tmp[2]
			optTypeTmp = tmp[3]
		}
	default:
		optTarget = strutil.Before(uri, "?")
		optTypeTmp = uri

	}

	switch optTypeTmp {
	case controller.ApiPathCreate:
		optType = CreateData
		optContent = fmt.Sprintf("%s新增了[%s]数据", us.Name, optTarget)
	case controller.ApiPathUpdate, controller.ApiPathUpdateByFilters:
		optType = UpdateData
		optContent = fmt.Sprintf("%s更新了[%s]的一些数据", us.Name, optTarget)
	case controller.ApiPathDelete, controller.ApiPathDeleteByFilters:
		optType = DeleteData
		optContent = fmt.Sprintf("%s删除了[%s]的一些数据", us.Name, optTarget)
	case controller.ApiPathQuery, controller.ApiPathPartition, "list":
		optType = QueryData
		optContent = fmt.Sprintf("%s查询了[%s]的数据", us.Name, optTarget)
	case controller.ApiPathCount, controller.ApiPathAggregate:
		optType = AggregateData
		optContent = fmt.Sprintf("%s查询了[%s]的统计数据", us.Name, optTarget)
	default:
		optType = strutil.Before(uri, "?")
		optContent = fmt.Sprintf("%s通过地址[%s]访问了系统", us.Name, uri)
	}
	return
}

func getHeaderDecode(c *gin.Context, key string) string {
	var (
		err error
		ct  = c.GetHeader(key)
	)
	ct, err = url.QueryUnescape(ct)
	if err != nil {
		logger.Errorf("UrlDecode Err %s", err)
	}
	return ct
}
