package ctx

import (
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"sync"
	"time"
)

const (
	DeltaGinContextKey      = "DeltaGinContextKey"
	GormContextKey          = "gormContextKey"
	GormAppKey              = "gormAppKey"
	DorisAppKey             = "dorisAppKey"
	UserIdContextKey        = "userIdContextKey"
	UserKeyContextKey       = "userKeyContextKey"
	UserNameContextKey      = "userNameContextKey"
	AppIdContextKey         = "appIdContextKey"
	CurrentDeptIdContextKey = "currentDeptIdContextKey"
	TokenContextKey         = "tokenContextKey"
	ManagerContextKey       = "managerContextKey"
	StorageContextKey       = "storageContextKey"
	NebulaContextKey        = "nebulaContextKey"
	FetchCacheContextKey    = "FetchCacheContextKey"
	OpFromEsContextKey      = "OpFromEsContextKey"
	TokenInfoContextKey     = "tokenInfoContextKey"
)

var (
	ctxInitOptions = make([]ContextOpt, 0, 10)
)

func AddNewContextOptions(f ...ContextOpt) {
	for i, _ := range f {
		ctxInitOptions = append(ctxInitOptions, f[i])
	}
}

type ContextOpt func(c *Context)

type Context struct {
	mu    sync.RWMutex
	keys  sync.Map
	reqId string
}

func (c *Context) Clone() *Context {
	cc := &Context{reqId: c.reqId}

	c.keys.Range(func(key, value any) bool {
		cc.keys.Store(key, value)
		return true
	})
	return cc
}

func GetContext(c *gin.Context) *Context {
	value, exists := c.Get(DeltaGinContextKey)
	if exists {
		return value.(*Context)
	}
	return NewContextWithGinContext(c)
}

func NewContext() *Context {
	c := &Context{
		reqId: ulid.Make().String(),
	}
	for _, v := range ctxInitOptions {
		v(c)
	}
	return c
}
func NewContextWithUserId(userId string) *Context {
	c := &Context{
		reqId: ulid.Make().String(),
	}
	c.SetContextUserId(userId)
	for _, v := range ctxInitOptions {
		v(c)
	}
	return c
}

func NewContextWithGinContext(c *gin.Context) *Context {
	value, exists := c.Get(DeltaGinContextKey)
	if exists {
		return value.(*Context)
	}
	ct := &Context{}
	BindGinContext(ct, c)
	for _, v := range ctxInitOptions {
		v(ct)
	}
	return ct
}
func BindGinContext(ct *Context, c *gin.Context) {
	ct.reqId = ulid.Make().String()
	ct.Set(DeltaGinContextKey, c)
	c.Set(DeltaGinContextKey, ct)
}

func (c *Context) GetRequestId() string {
	return c.reqId
}

func (c *Context) GetGinContext() *gin.Context {
	value, exists := c.Get(DeltaGinContextKey)
	if exists {
		return value.(*gin.Context)
	}
	return nil
}

func (c *Context) Get(key string) (value any, exists bool) {
	value, exists = c.keys.Load(key)
	if !exists {
		key = utils.HumpToUnderline(key)
		return c.keys.Load(key)
	}
	return
}
func (c *Context) Set(key, value any) {
	c.keys.Store(key, value)
}

func (c *Context) GetString(key string) (s string) {
	if val, ok := c.Get(key); ok && val != nil {
		s, _ = val.(string)
	}
	return
}

func (c *Context) GetBool(key string) (b bool) {
	if val, ok := c.Get(key); ok && val != nil {
		b, _ = val.(bool)
	}
	return
}

func (c *Context) GetInt(key string) (i int) {
	if val, ok := c.Get(key); ok && val != nil {
		i, _ = val.(int)
	}
	return
}

func (c *Context) GetInt64(key string) (i64 int64) {
	if val, ok := c.Get(key); ok && val != nil {
		i64, _ = val.(int64)
	}
	return
}

func (c *Context) GetUint(key string) (ui uint) {
	if val, ok := c.Get(key); ok && val != nil {
		ui, _ = val.(uint)
	}
	return
}

func (c *Context) GetUint64(key string) (ui64 uint64) {
	if val, ok := c.Get(key); ok && val != nil {
		ui64, _ = val.(uint64)
	}
	return
}

func (c *Context) GetFloat64(key string) (f64 float64) {
	if val, ok := c.Get(key); ok && val != nil {
		f64, _ = val.(float64)
	}
	return
}

func (c *Context) GetTime(key string) (t time.Time) {
	if val, ok := c.Get(key); ok && val != nil {
		t, _ = val.(time.Time)
	}
	return
}

func (c *Context) GetDuration(key string) (d time.Duration) {
	if val, ok := c.Get(key); ok && val != nil {
		d, _ = val.(time.Duration)
	}
	return
}

func (c *Context) GetStringSlice(key string) (ss []string) {
	if val, ok := c.Get(key); ok && val != nil {
		ss, _ = val.([]string)
	}
	return
}

func (c *Context) GetStringMap(key string) (sm map[string]any) {
	if val, ok := c.Get(key); ok && val != nil {
		sm, _ = val.(map[string]any)
	}
	return
}

func (c *Context) GetStringMapString(key string) (sms map[string]string) {
	if val, ok := c.Get(key); ok && val != nil {
		sms, _ = val.(map[string]string)
	}
	return
}

func (c *Context) GetStringMapStringSlice(key string) (smss map[string][]string) {
	if val, ok := c.Get(key); ok && val != nil {
		smss, _ = val.(map[string][]string)
	}
	return
}
