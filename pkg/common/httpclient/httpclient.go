package httpclient

import (
	"io"
	"net/http"
	"sync"
)

var (
	DefaultHttpClient = defaultHttpClient()
)

const (
	MIMEJSON              = "application/json"
	MIMEHTML              = "text/html"
	MIMEXML               = "application/xml"
	MIMEXML2              = "text/xml"
	MIMEPlain             = "text/plain"
	MIMEPOSTForm          = "application/x-www-form-urlencoded"
	MIMEMultipartPOSTForm = "multipart/form-data"
	MIMEPROTOBUF          = "application/x-protobuf"
	MIMEMSGPACK           = "application/x-msgpack"
	MIMEMSGPACK2          = "application/msgpack"
	MIMEYAML              = "application/x-yaml"
	MIMETOML              = "application/toml"
)

type HttpClient struct {
	client *http.Client
	Error  error
	pool   *clientPool
	clt    *client
	doOne  sync.Once
}

func (h *HttpClient) Method(mtd string) *HttpClient {
	hc := h.getInstance()
	hc.Error = hc.clt.method(mtd).err
	return hc
}
func (h *HttpClient) Url(url string) *HttpClient {
	hc := h.getInstance()
	hc.Error = hc.clt.url(url).err
	return hc
}
func (h *HttpClient) ContentType(contentType string) *HttpClient {
	hc := h.getInstance()
	hc.Error = hc.clt.contentType(contentType).err
	return hc
}
func (h *HttpClient) Body(contentType string, reqBody any) *HttpClient {
	hc := h.getInstance()
	hc.Error = hc.clt.body(contentType, reqBody).err
	return hc
}
func (h *HttpClient) BodyReader(contentType string, reqBody io.Reader) *HttpClient {
	hc := h.getInstance()
	hc.Error = hc.clt.bodyReader(contentType, reqBody).err
	return hc
}
func (h *HttpClient) AddHeader(key, value string) *HttpClient {
	hc := h.getInstance()
	hc.Error = hc.clt.addHeader(key, value).err
	return hc
}
func (h *HttpClient) SetHeader(key, value string) *HttpClient {
	hc := h.getInstance()
	hc.Error = hc.clt.setHeader(key, value).err
	return hc
}
func (h *HttpClient) SetHeaders(headers map[string]string) *HttpClient {
	hc := h.getInstance()
	for k, v := range headers {
		hc.Error = hc.clt.setHeader(k, v).err
	}
	return hc
}
func (h *HttpClient) Bind(bind Binding) *HttpClient {
	hc := h.getInstance()
	hc.Error = hc.clt.bind(bind).err
	return hc
}
func (h *HttpClient) Build() *HttpClient {
	hc := h.getInstance()
	hc.Error = hc.clt.build().err
	return hc
}
func (h *HttpClient) Do(respBody any) *HttpClient {
	hc := h.getInstance()
	hc.clt.do(respBody, nil)
	cl := hc.clt
	hc.clt = nil
	hc.Error = hc.pool.Put(cl)
	return hc
}

func (h *HttpClient) Get(url string, dst any) error {
	err := h.
		Method("GET").
		Url(url).
		Build().
		Do(dst).Error
	return err
}
func (h *HttpClient) Post(url, contentType string, reqBody, respBody any) error {
	err := h.
		Method("POST").
		Url(url).
		Body(contentType, reqBody).
		Build().
		Do(respBody).Error
	return err
}
func (h *HttpClient) Put(url, contentType string, reqBody, respBody any) error {
	err := h.
		Method("PUT").
		Url(url).
		Body(contentType, reqBody).
		Build().
		Do(respBody).Error
	return err
}

func (h *HttpClient) DoHandle(handler ResponseHandler) *HttpClient {
	hc := h.getInstance()
	hc.clt.doWithHandler(handler)
	cl := hc.clt
	hc.clt = nil
	hc.Error = hc.pool.Put(cl)
	return hc
}

func (h *HttpClient) getInstance() *HttpClient {
	if h.clt == nil {
		return &HttpClient{
			client: h.client,
			pool:   h.pool,
			clt:    h.pool.Get(),
		}
	}
	return h
}

func defaultHttpClient() *HttpClient {
	// can config http.Transport
	ct := &http.Client{
		Transport: http.DefaultTransport,
	}
	return NewHttpClient(ct)
}

func NewHttpClient(clt *http.Client) *HttpClient {
	hc := &HttpClient{}
	hc.pool = newClientPool(clt, func() any {
		return &client{}
	})
	return hc
}
