package cache

import (
	"github.com/davycun/eta/pkg/common/cache"
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/core/controller"
	"github.com/duke-git/lancet/v2/pointer"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/duke-git/lancet/v2/strutil"
	"github.com/gin-gonic/gin"
	"time"
)

type Controller struct {
	controller.DefaultController
}

func (handler Controller) scan(c *gin.Context) {
	var (
		err    error
		param  ScanParam
		result = ScanResult{}
	)

	err = caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			return controller.BindBody(c, &param)
		}).
		Call(func(cl *caller.Caller) error {
			keys, cur, err1 := cache.Scan(param.Cursor, param.Match, param.Count)
			if err1 == nil {
				result.Cursor = cur
				result.Keys = keys
			}
			return err1
		}).Err
	handler.ProcessResult(c, &result, err)
}

func (handler Controller) detail(c *gin.Context) {
	var (
		key struct {
			Key string `json:"key,omitempty" uri:"key" binding:"required"`
		}
		result = DetailResult{}
		err    error
	)

	err = caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			return controller.BindUri(c, &key)
		}).
		Call(func(cl *caller.Caller) error {
			err1, val, dur := cache.Detail(key.Key)
			if err1 != nil {
				return err1
			}
			if dur != nil {
				result.Ttl = pointer.Of(int64(dur.Seconds()))
			}
			result.Value = val
			return err1
		}).Err
	handler.ProcessResult(c, &result, err)
}

func (handler Controller) set(c *gin.Context) {
	var (
		err   error
		param SetParam
	)
	err = caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			return controller.BindBody(c, &param)
		}).
		Call(func(cl *caller.Caller) error {
			if param.Expiration <= 0 {
				err1 := cache.Set(param.Key, param.Value)
				if err1 != nil {
					return err1
				}
			} else {
				err1 := cache.SetEx(param.Key, param.Value, time.Second*time.Duration(param.Expiration))
				if err1 != nil {
					return err1
				}
			}
			return nil
		}).Err
	handler.ProcessResult(c, nil, err)
}

func (handler Controller) del(c *gin.Context) {
	var (
		err   error
		param DelParam
		keys  = make([]string, 0)
	)
	err = caller.NewCaller().
		// 绑定参数
		Call(func(cl *caller.Caller) error {
			return controller.BindBody(c, &param)
		}).
		// 查询 key
		Call(func(cl *caller.Caller) error {
			keys = append(keys, param.Keys...)
			for _, key := range param.Keys {
				if strutil.ContainsAny(key, []string{"*", "?", "[", "]"}) {
					var cursor uint64 = 0
					for {
						ks, cursor1, err1 := cache.Scan(cursor, key, 10000)
						if err1 != nil {
							return err1
						}
						keys = append(keys, ks...)
						cursor = cursor1
						if cursor1 == 0 {
							break
						}
					}
				} else {
					keys = append(keys, key)
				}
			}
			keys = slice.Unique(keys)
			return nil
		}).
		// 异步删除
		Call(func(cl *caller.Caller) error {
			return cache.Unlink(keys...)
		}).Err
	handler.ProcessResult(c, nil, err)
}
