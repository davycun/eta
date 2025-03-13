package loader

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/cache"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/eta/constants"
	"gorm.io/gorm"
	"reflect"
	"strings"
	"sync"
)

const (
	allKey = "all"
)

// Cache
// T代表实体，V可能是实体，也可以是切片
type Cache[T, V any] struct {
	hasAll    *sync.Map //key:appId ,values is appId
	store     *sync.Map //key: appId,value is sync.Map
	cacheKey  string    //因为这个Cache是本地的，为了解决多节点数据不一致问题，所以通过redis的cacheKey来通知所有节点进行本地缓存更新
	mutex     *sync.Mutex
	cols      []string
	ldCfg     EntityLoaderConfig
	valueType reflect.Type
}

// NewCacheLoader
// cacheKey 表示的是redis记录的一个缓存key，主要是多节点应用的时候用于通知本地缓存更新
// 使用方式，请搜索已经在使用的代码
func NewCacheLoader[T, V any](tbName string, cacheKey string, cols ...string) *Cache[T, V] {
	c := &Cache[T, V]{
		hasAll:    &sync.Map{},
		valueType: reflect.TypeFor[V](),
		store:     &sync.Map{},
		mutex:     &sync.Mutex{},
		cacheKey:  cacheKey,
		ldCfg: EntityLoaderConfig{
			TableName:            tbName,
			DefaultEntityColumns: cols,
		},
	}
	cache.AddAfterDel(c.receiveDelKeyFromRedis)
	return c
}

// no publish
func (c *Cache[T, V]) receiveDelKeyFromRedis(keys ...string) {
	for _, v := range keys {
		if strings.HasPrefix(v, strings.ReplaceAll(c.cacheKey, "%s", "")) {
			logger.Infof("delete cacheLoader key %s", v)
			appId, dataId := splitPublishDelKey(v)
			c.getStore(appId).Delete(dataId)
		}
	}
}

func (c *Cache[T, V]) SetKeyName(key string) *Cache[T, V] {
	c.ldCfg.IdColumn = key
	return c
}

func (c *Cache[T, V]) LoadData(db *gorm.DB, keys ...string) (map[string]V, error) {

	var (
		scm = dorm.GetDbSchema(db)
	)
	mp, notExistsIds := c.loadExists(scm, keys...)
	if len(notExistsIds) < 1 {
		return mp, nil
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	mp, notExistsIds = c.loadExists(scm, keys...)
	if len(notExistsIds) < 1 {
		return mp, nil
	}

	var (
		dataList []T
		err      error
	)
	c.ldCfg.Ids = notExistsIds
	if c.ldCfg.IdColumn == "" {
		c.ldCfg.IdColumn = "id"
	}

	err = NewEntityLoader(db, c.ldCfg).Load(&dataList)
	if len(dataList) > 0 {
		c.addCache(scm, dataList...)
	}
	mp, _ = c.loadExists(scm, keys...)
	return mp, err
}
func (c *Cache[T, V]) LoadAll(db *gorm.DB) (map[string]V, error) {

	var (
		scm = dorm.GetDbSchema(db)
	)
	if _, ok := c.hasAll.Load(scm); ok {
		return c.getAll(scm), nil
	}
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if _, ok := c.hasAll.Load(scm); ok {
		return c.getAll(scm), nil
	}

	var (
		dataList []T
	)

	err := dorm.Table(db, c.ldCfg.TableName).Select(c.ldCfg.DefaultEntityColumns).Find(&dataList).Error
	if len(dataList) > 0 {
		c.addCache(scm, dataList...)
		c.hasAll.Store(scm, scm)
	}
	return c.getAll(scm), err
}
func (c *Cache[T, V]) Delete(db *gorm.DB, keys ...string) {
	var (
		scm = dorm.GetDbSchema(db)
	)
	for _, v := range keys {
		c.getStore(scm).Delete(v)
		k := concatPublishDelKey(c.cacheKey, scm, v)
		err := cache.PublishDelKey(k)
		if err != nil {
			logger.Errorf("cacheLoader publish delete key [%s] error %s", k, err)
		}
		if _, ok := c.hasAll.Load(scm); ok {
			c.hasAll.Delete(scm)
		}
	}
}
func (c *Cache[T, V]) HasAll(db *gorm.DB) bool {
	var (
		scm = dorm.GetDbSchema(db)
	)
	_, ok := c.hasAll.Load(scm)
	return ok
}
func (c *Cache[T, V]) SetHasAll(db *gorm.DB, hasAll bool) {
	var (
		scm = dorm.GetDbSchema(db)
	)
	if hasAll {
		c.hasAll.Store(scm, scm)
	} else {
		c.hasAll.Delete(scm)
	}
}
func (c *Cache[T, V]) DeleteAll(db *gorm.DB) {
	var (
		scm = dorm.GetDbSchema(db)
	)
	c.store.Delete(scm)
}
func (c *Cache[T, V]) DeleteAllBySchema(scm string) {
	c.store.Delete(scm)
}

func (c *Cache[T, V]) getAll(appId string) map[string]V {
	mp := make(map[string]V)
	c.getStore(appId).Range(func(key, value any) bool {
		mp[key.(string)] = value.(V)
		return true
	})
	return mp
}
func (c *Cache[T, V]) loadExists(appId string, keys ...string) (exists map[string]V, notExistsKeys []string) {
	exists = make(map[string]V)
	for _, v := range keys {
		ad, ok := c.getStore(appId).Load(v)
		if ok {
			exists[v] = ad.(V)
			continue
		}
		notExistsKeys = utils.Merge(notExistsKeys, v)
	}
	return exists, notExistsKeys
}
func (c *Cache[T, V]) getStore(appId string) *sync.Map {
	if x, ok := c.store.Load(appId); ok {
		return x.(*sync.Map)
	}
	st := &sync.Map{}
	c.store.Store(appId, st)
	return st
}

func (c *Cache[T, V]) addCache(appId string, dataList ...T) {
	if len(dataList) < 1 {
		return
	}
	for _, v := range dataList {
		key := getKey(c.ldCfg.IdColumn, v)
		if key == "" {
			continue
		}

		if isSlice(c.valueType) {
			ds, ok := c.getStore(appId).Load(key)
			if ok {
				ds = append(ds.([]T), v)
			} else {
				tmp := make([]T, 0, 1)
				tmp = append(tmp, v)
				ds = tmp
			}
			c.getStore(appId).Store(key, ds)
		} else {
			c.getStore(appId).Store(key, v)
		}
	}
}

func getKey(key string, v any) string {
	return entity.GetString(v, key)
}

func isSlice(vType reflect.Type) bool {
	switch vType.Kind() {
	case reflect.Slice:
		return true
	case reflect.Pointer:
		return isSlice(vType.Elem())
	default:

	}
	return false
}

func concatPublishDelKey(redisKey, appId, dataId string) string {
	return constants.RedisKey(redisKey, fmt.Sprintf("%s@%s", appId, dataId))
}
func splitPublishDelKey(publishKey string) (string, string) {
	split := strings.Split(publishKey, ":")
	if len(split) < 1 {
		logger.Warnf("split publish key [%s] error", publishKey)
		return "", ""
	}
	pk := split[len(split)-1]
	appIdAndDataId := strings.Split(pk, "@")
	if len(appIdAndDataId) != 2 {
		logger.Warnf("split publish key [%s] error", publishKey)
		return "", ""
	}
	return appIdAndDataId[0], appIdAndDataId[1]
}

func getSlicePointer(tp reflect.Type) any {
	switch tp.Kind() {
	case reflect.Struct, reflect.Map:
		return reflect.New(reflect.SliceOf(tp)).Interface()
	case reflect.Slice:
		return reflect.New(tp).Interface()
	case reflect.Pointer:
		return getSlicePointer(tp.Elem())
	default:

	}
	return nil
}
