package forward_srv

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/core/controller"
	"github.com/davycun/eta/pkg/module/forward"
	"github.com/davycun/eta/pkg/module/setting"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"io"
	"strings"
	"time"
)

func Forward(c *gin.Context) {
	hd := newHandler(c)
	err := hd.do(hd.newRequest()).Err()
	if err != nil {
		controller.ProcessResult(c, nil, err)
	}
	return
}

type handler struct {
	c      *gin.Context
	ct     *ctx.Context
	vendor string
	cred   setting.BaseCredentials
	err    error
	client *resty.Client
}

func newHandler(c *gin.Context) *handler {
	hd := &handler{
		c: c,
	}
	hd.vendor = c.Param("vendor")
	hd.ct = ctx.GetContext(c)
	hd.cred, hd.err = forward.GetVendor(hd.ct.GetAppGorm(), hd.vendor)

	if hd.err != nil {
		return hd
	}

	hd.client = resty.New()
	if hd.cred.ProxyUrl != "" {
		hd.client.SetProxy(hd.cred.ProxyUrl)
	}
	hd.client.SetBaseURL(hd.cred.BaseUrl).
		SetDebug(hd.cred.Debug).
		SetTimeout(time.Second * time.Duration(hd.cred.Timeout)). //A Timeout of zero means no timeout.
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: hd.cred.InsecureSkipVerify})

	return hd
}

func (h *handler) newRequest() *resty.Request {

	var (
		err     error
		req     = h.client.R()
		reqBody []byte
	)

	h.err = caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			for k := range h.c.Request.Header {
				req.SetHeader(k, h.c.GetHeader(k))
			}
			return h.err
		}).
		Call(func(cl *caller.Caller) error {
			if h.c.Request.Body != nil {
				reqBody, err = io.ReadAll(h.c.Request.Body)
				if err != nil {
					h.err = errors.New(fmt.Sprintf("read request body error %s", err))
				}
			}
			return h.err
		}).
		Call(func(cl *caller.Caller) error {
			if fc, ok := forward.GetHandleRequest(h.vendor); ok {
				realBody, err1 := fc(req, h.cred, reqBody)
				if err1 != nil {
					h.err = errors.New(fmt.Sprintf("request handler error %s", err1))
				}
				if len(realBody) > 0 {
					reqBody = realBody
				}
			}
			return h.err
		}).
		Call(func(cl *caller.Caller) error {
			if len(reqBody) > 0 {
				req.SetBody(reqBody)
			}
			return h.err
		}).Err

	return req
}

func (h *handler) do(req *resty.Request) *handler {

	if h.err != nil {
		return h
	}

	var (
		realPath = h.c.Param("path") //注意前端不要调用的时候传入多个斜杠
		resp     *resty.Response
	)
	resp, h.err = req.Execute(h.c.Request.Method, realPath)

	if resp == nil {
		return h
	}

	//如果body为空可能是GET或者其他原因，如果之前发生了错误，就用delta的格式返回，如果没有就直接返回
	var (
		respBody = resp.Body()
	)
	writeHeader(h.c, resp)
	if len(respBody) < 1 {
		h.c.Abort()
		return h
	}

	if fc, ok := forward.GetHandleResponse(h.vendor); ok {
		realBody, err1 := fc(resp, h.cred, respBody)
		if err1 != nil {
			_ = h.c.Error(err1)
			//_ = h.c.Error(errors.New(fmt.Sprintf("response handler error %s", err1)))
		}
		if len(realBody) > 0 {
			respBody = realBody
		}
	}

	if len(respBody) > 0 {
		_, err := h.c.Writer.Write(respBody)
		if err != nil {
			_ = h.c.Error(err)
		}
	}
	h.c.Abort()
	return h
}

func (h *handler) Err() error {
	return h.err
}

func writeHeader(c *gin.Context, resp *resty.Response) {
	c.Writer.WriteHeader(resp.StatusCode())
	for k := range resp.Header() {
		//取消掉跨域相关的头信息和压缩相关的头信息
		//resty和net/http都会处理gzip压缩，如果响应头中为Content-Encoding：gzip，说明服务端已经做了gzip压缩
		//那么 net/http会做解压，可以定制http.Transport.DisableCompression为true来拒绝net/http解压
		//但是如果net/http不解压了，resty还是会解压，但是resty不支持配置不解压
		//所以从resty拿到的body是解压后的，直接往客户端写，但如果给客户端的响应头没有去除Content-Encoding: gzip，那么客户端会拿非压缩的数据进行gzip解压会报错
		if strings.HasPrefix(k, "Cross") || strings.HasPrefix(k, "Content-Encoding") {
			continue
		}
		c.Header(k, resp.Header().Get(k))
	}
}
