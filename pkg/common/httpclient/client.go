package httpclient

import (
	"bytes"
	"context"
	"errors"
	"github.com/davycun/eta/pkg/common/utils"
	"io"
	"net/http"
	"strings"
	"sync"
)

type ResponseHandler func(resp *http.Response) error

type client struct {
	forceBind     bool
	rqMethod      string
	rqUrl         string
	rqContentType string
	rqBody        io.Reader
	rqHeader      http.Header

	req     *http.Request
	binding Binding
	client  *http.Client
	resp    *http.Response

	called bool
	closed bool
	err    error
}

func (c *client) url(url string) *client {
	c.rqUrl = url
	return c
}
func (c *client) method(method string) *client {
	c.rqMethod = method
	return c
}
func (c *client) contentType(contentType string) *client {
	c.rqContentType = contentType
	return c
}
func (c *client) body(contentType string, reqBody any) *client {
	c.rqContentType = contentType
	switch reqBody.(type) {
	case string:
		c.rqBody = strings.NewReader(reqBody.(string))
	case []byte:
		c.rqBody = bytes.NewReader(reqBody.([]byte))
	default:
		bd := GetBinding(contentType)
		c.rqBody, c.err = bd.UnBind(reqBody)
	}
	return c
}
func (c *client) bodyReader(contentType string, reqBody io.Reader) *client {
	c.rqContentType = contentType
	c.rqBody = reqBody
	return c
}
func (c *client) addHeader(key, value string) *client {
	c.rqHeader.Add(key, value)
	return c
}
func (c *client) setHeader(key, value string) *client {
	c.rqHeader.Set(key, value)
	return c
}

// 用来解析响应的Bingding，如果不设置会根据Response的Content-Type进行设置
func (c *client) bind(bind Binding) *client {
	c.binding = bind
	return c
}
func (c *client) build() *client {
	if c.err != nil {
		return c
	}
	if c.rqUrl == "" {
		c.err = errors.Join(c.err, errors.New("you should set url before build"))
	}
	if c.err == nil && c.req == nil {
		c.req, c.err = http.NewRequestWithContext(context.Background(), c.rqMethod, c.rqUrl, c.rqBody)
		if len(c.rqHeader) > 0 {
			c.req.Header = c.rqHeader
		}
		if c.rqContentType != "" {
			c.req.Header.Set("Content-Type", c.rqContentType)
		}
	}
	return c
}

func (c *client) do(respBody any, handle ResponseHandler) *client {
	if c.err != nil || c.called {
		return c
	}
	c.build()
	if c.req == nil || c.err != nil {
		return c
	}
	c.resp, c.err = c.client.Do(c.req)
	c.called = true
	defer func() {
		c.closeBody()
	}()
	if c.err != nil {
		return c
	}
	rs := readBody(c, c.resp)
	//debug c.resp.Write(logger.Logger.Writer())
	//ss := utils.BytesToString(rs)
	//logger.Debugf("the request url[%s] response body is: %s", c.rqUrl, ss)
	switch respBody.(type) {
	case *string:
		*(respBody.(*string)) = utils.BytesToString(rs)
	default:
		if len(rs) > 0 {
			c.err = c.initBind().binding.Bind(rs, respBody)
		}
	}
	if c.resp.StatusCode != http.StatusOK {
		if len(rs) > 0 {
			c.err = errors.Join(c.err, errors.New(utils.BytesToString(rs)))
		}
	}

	return c
}
func (c *client) doWithHandler(handle ResponseHandler) *client {
	if c.err != nil || c.called {
		return c
	}
	c.build()
	if c.req == nil || c.err != nil {
		return c
	}
	c.resp, c.err = c.client.Do(c.req)
	c.called = true
	defer func() {
		c.closeBody()
	}()
	if c.err != nil {
		return c
	}
	c.err = handle(c.resp)
	return c
}

func readBody(c *client, resp *http.Response) []byte {
	dt, err := io.ReadAll(resp.Body)
	if err != nil && err != io.EOF {
		c.err = errors.Join(c.err, err)
	}
	return dt
}

// 根据Response的Content-Type来决定采用什么格式解析body
func (c *client) initBind() *client {
	if c.binding == nil {
		if c.resp != nil {
			ct := c.resp.Header.Get("Content-Type")
			if i := strings.Index(ct, ";"); i > 0 {
				ct = ct[:i]
			}
			if ct == "" {
				ct = c.rqContentType
			}
			c.binding = GetBinding(ct)
		} else {
			c.binding = GetBinding(c.rqContentType)
		}
	}
	return c
}
func (c *client) joinError(err ...error) *client {
	c.err = errors.Join(err...)
	return c
}

func (c *client) closeBody() {
	if c.resp != nil && c.resp.Body != nil {
		err := c.resp.Body.Close()
		c.err = errors.Join(c.err, err)
		c.closed = true
	}
}

func (c *client) reset() *client {
	c.rqMethod = "GET"
	c.rqUrl = ""
	//c.rqContentType = MIMEJSON
	c.rqContentType = ""
	c.rqBody = nil
	c.rqHeader = make(http.Header)
	c.req = nil
	c.binding = nil
	c.client = nil
	c.resp = nil
	c.called = false
	c.closed = false
	c.err = nil
	c.forceBind = true
	return c
}

func (c *client) clear() {
	c.rqMethod = ""
	c.rqUrl = ""
	c.rqContentType = ""
	c.rqBody = nil
	c.rqHeader = nil
	c.req = nil
	c.binding = nil
	c.client = nil
	c.resp = nil
	c.called = false
	c.closed = false
	c.err = nil
}

type clientPool struct {
	client *http.Client
	pool   sync.Pool
}

func (p *clientPool) Get() *client {
	c := p.pool.Get().(*client)
	c.reset()
	c.client = p.client
	return c
}
func (p *clientPool) Put(c *client) error {
	c.closeBody()
	err := c.err
	c.clear()
	p.pool.Put(c)
	return err
}

func newClientPool(clt *http.Client, newFunc func() any) *clientPool {
	pool := clientPool{}
	pool.client = clt
	pool.pool.New = newFunc
	return &pool
}
