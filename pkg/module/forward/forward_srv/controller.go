package forward_srv

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/controller"
	"github.com/davycun/eta/pkg/module/forward"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

func Forward(c *gin.Context) {
	var (
		err error
		hd  = newHandler(c)
	)
	if hd.err != nil {
		controller.ProcessResult(c, nil, hd.err)
		return
	}
	err = hd.do(hd.newRequest()).Err()

	if err != nil {
		controller.ProcessResult(c, nil, err)
	}
	return
}

type handler struct {
	c          *gin.Context
	ct         *ctx.Context
	vendorName string
	vendor     forward.Vendor
	err        error
	client     *resty.Client
	cacheKey   string //请求缓存的hash key
	realPath   string
	needCache  bool
}

func (h *handler) Err() error {
	return h.err
}

func newHandler(c *gin.Context) *handler {
	hd := &handler{
		c: c,
	}
	hd.vendorName = c.Param(forward.PathVendor)
	hd.ct = ctx.GetContext(c)
	hd.vendor, hd.err = forward.GetVendor(hd.ct.GetAppGorm(), hd.vendorName)

	if hd.err != nil {
		return hd
	}

	hd.client = resty.New()
	if hd.vendor.ProxyUrl != "" {
		hd.client.SetProxy(hd.vendor.ProxyUrl)
	}
	hd.client.SetBaseURL(hd.vendor.RandomBaseUrl()).
		SetDebug(hd.vendor.Debug).
		SetTimeout(time.Second * time.Duration(hd.vendor.Timeout)). //A Timeout of zero means no timeout.
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: hd.vendor.InsecureSkipVerify})
	if hd.vendor.Debug {
		hd.client.EnableGenerateCurlOnDebug()
	}

	hd.realPath = c.Param(forward.PathParam)
	hd.realPath = fmt.Sprintf("/%s", strings.TrimLeft(hd.realPath, "/")) //保证绝对路径斜杠开头，也避免前端路径传入多个斜杠
	hd.needCache = hd.vendor.NeedCache(hd.c.Request.Method, hd.realPath)

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
			req.SetQueryParamsFromValues(h.c.Request.URL.Query())
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
			if fc, ok := forward.GetHandleRequest(h.vendorName); ok {
				realBody, err1 := fc(req, h.vendor, reqBody)
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
		}).
		Call(func(cl *caller.Caller) error {
			if h.needCache {
				h.cacheKey, h.err = forward.MakeCacheKey(h.c, reqBody, h.vendorName)
			}
			return h.err
		}).Err

	return req
}

func (h *handler) do(req *resty.Request) *handler {

	var (
		resp      *resty.Response
		cacheData = forward.CacheData{}
	)
	if h.needCache && h.err == nil {
		cacheData, h.err = forward.LoadCacheData(h.cacheKey, h.vendor)
		//说明cache有效
		if cacheData.IsValid() {
			h.c.Header("X-Cache-Key", h.cacheKey)
			h.err = h.writeResponse(cacheData)
			return h
		}
	}
	if h.err != nil {
		return h
	}

	for _, v := range h.vendor.ExcludeRequestHeader {
		req.Header.Del(v)
	}

	//调用实际请求
	resp, h.err = req.Execute(h.c.Request.Method, h.realPath)
	if resp == nil {
		return h
	}
	cacheData, h.err = forward.MakeCacheData(h.vendor, resp)
	if h.err != nil {
		return h
	}
	h.err = errs.Cover(h.err, h.writeResponse(cacheData))
	if h.needCache && h.err == nil {
		err := forward.SaveCacheData(h.cacheKey, h.vendor, cacheData)
		if err != nil {
			logger.Errorf("save cache data err %s", err)
		}
	}
	return h
}

func (h *handler) writeResponse(data forward.CacheData) error {
	if h.err != nil {
		return h.err
	}
	h.c.Writer.WriteHeader(data.Status)
	for k := range data.Header {
		h.c.Header(k, data.Header.Get(k))
	}
	if len(data.Body) > 0 {
		_, err := h.c.Writer.Write(data.Body)
		if err != nil {
			_ = h.c.Error(err)
		}
	}
	h.c.Abort()
	return h.err
}
