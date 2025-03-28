package plugin_push

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/httpclient"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/run"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/core/service"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/broker/publish"
	"github.com/davycun/eta/pkg/module/broker/subscribe"
	jsoniter "github.com/json-iterator/go"
	"io"
	"net/http"
)

// PublishModifyCallbacks
// 如果订阅了某张表数据的接口，通过这个插件进行推送
func PublishModifyCallbacks(cfg *hook.SrvConfig, pos hook.CallbackPosition) error {

	switch pos {
	case hook.CallbackBefore:
	case hook.CallbackAfter:
		switch cfg.Method {
		case iface.MethodCreate, iface.MethodUpdate, iface.MethodUpdateByFilters:
			run.Go(func() {
				err := afterModify(cfg, cfg.NewValues)
				if err != nil {
					logger.Errorf("回调接口发生错误%s", err)
				}
			})
		case iface.MethodDelete, iface.MethodDeleteByFilters:
			run.Go(func() {
				err := afterModify(cfg, cfg.OldValues)
				if err != nil {
					logger.Errorf("回调接口发生错误%s", err)
				}
			})
		default:

		}
	default:
	}
	return nil
}
func afterModify(cfg *hook.SrvConfig, data any) error {
	target, err := subscribe.LoadSubscriberByTarget(cfg.Ctx.GetAppGorm(), cfg.GetTableName())
	if err != nil {
		return err
	}

	var (
		rds = make([]publish.Record, 0, len(target))
	)

	for _, v := range target {
		var (
			param = dto.ModifyParam{}
			rs    = publish.Record{}
			clt   *httpclient.HttpClient
		)
		param.Data = data
		switch v.Method {
		case http.MethodPost:
			clt = httpclient.DefaultHttpClient.Method(http.MethodPost)
		default:
			clt = httpclient.DefaultHttpClient.Method(http.MethodPost)
		}

		rs.Request, err = jsoniter.MarshalToString(&param)

		if len(v.Header) > 0 {
			for _, val := range v.Header {
				for _, vs := range val.Values {
					clt.AddHeader(val.Key, vs)
				}
			}
		}

		clt.Url(v.Url).
			Body(httpclient.MIMEJSON, &param).
			AddHeader("Content-Type", httpclient.MIMEJSON).
			DoHandle(func(resp *http.Response) error {
				rs.SubId = v.ID
				rs.Status = fmt.Sprintf("%d", resp.StatusCode)
				rs.Count = 1
				if resp.StatusCode != 200 {
					logger.Errorf("publish err for {subId:%s,url:%s}", v.ID, v.Url)
					dt, err2 := io.ReadAll(resp.Body)
					if err2 != nil {
						logger.Errorf("read publish call err %s", err2)
					}
					if len(dt) > 0 {
						rs.Response = utils.BytesToString(dt)
					}
				}
				return nil
			})
		rds = append(rds, rs)
	}
	if len(rds) < 1 {
		return nil
	}
	//因为是异步，所以不能cfg.TxDB(接口返回前就被提交了）
	return service.NewSrvWrapper(constants.TablePublishRecord, cfg.Ctx, cfg.Ctx.GetAppGorm()).SetData(&rds).Create()
}
