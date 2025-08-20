package loader

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/cache"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/builder"
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
	hasAll       *sync.Map //key:appId ,values is appId
	store        *sync.Map //key: appId,value is sync.Map
	cacheKey     string    //因为这个Cache是本地的，为了解决多节点数据不一致问题，所以通过redis的cacheKey来通知所有节点进行本地缓存更新
	mutex        *sync.Mutex
	ldCfg        EntityLoaderConfig
	valueType    reflect.Type
	extraKey     []string
	extraKey2Key map[string]string
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
			columns:   cols,
			tableName: tbName,
			idColumn:  entity.IdDbName,
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
	c.ldCfg.idColumn = key
	return c
}

// AddExtraKeyName
// 比如有一个表有id、code、phone、email三个唯一字段，主key是id，那么这个时候就可以添加code、phone、email这些额外的key
// 这样通过LoadKey的时候，如果找不到主key，那么会去查找额外的key对应的主key，然后再通过找到的主key再找数据
func (c *Cache[T, V]) AddExtraKeyName(key ...string) *Cache[T, V] {
	c.extraKey = utils.Merge(c.extraKey, key...)
	return c
}

func (c *Cache[T, V]) LoadData(db *gorm.DB, keyValues ...string) (map[string]V, error) {

	var (
		appId       = dorm.GetAppId(db)
		keyNameList = []string{c.ldCfg.idColumn} //主key放第一位，理论上可以通过keyValues来初步判断这个值是主key的值还是副key的值
	)
	keyNameList = append(keyNameList, c.extraKey...)

	//只要有一个类型的key的数据加载到了就直接返回
	for _, keyName := range keyNameList {
		dt, err := c.loadData(db, keyName, keyValues...)
		if err != nil || len(dt) > 0 {
			return dt, err
		}
	}

	mp, _ := c.loadExists(appId, keyValues...)
	return mp, nil
}

func (c *Cache[T, V]) LoadAll(db *gorm.DB) (map[string]V, error) {

	var (
		appId = dorm.GetAppId(db)
	)
	if _, ok := c.hasAll.Load(appId); ok {
		return c.getAll(appId), nil
	}
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if _, ok := c.hasAll.Load(appId); ok {
		return c.getAll(appId), nil
	}

	dataList, err := c.selectData(db, c.ldCfg.idColumn)
	if len(dataList) > 0 {
		c.addCache(appId, dataList...)
		c.hasAll.Store(appId, appId)
	}
	return c.getAll(appId), err
}

func (c *Cache[T, V]) loadData(db *gorm.DB, keyName string, keyValues ...string) (map[string]V, error) {

	var (
		appId = dorm.GetAppId(db)
	)
	mp, notExistsIds := c.loadExists(appId, keyValues...)
	if len(notExistsIds) < 1 {
		return mp, nil
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	mp, notExistsIds = c.loadExists(appId, keyValues...)
	if len(notExistsIds) < 1 {
		return mp, nil
	}

	dataList, err := c.selectData(db, keyName, notExistsIds...)
	if len(dataList) > 0 {
		c.addCache(appId, dataList...)
	}
	mp, _ = c.loadExists(appId, keyValues...)
	return mp, err
}

func (c *Cache[T, V]) selectData(db *gorm.DB, keyName string, keyValues ...string) ([]T, error) {
	var (
		dataList []T
		err      error
	)
	var (
		tableName = c.ldCfg.tableName
		cols      = c.ldCfg.columns
	)
	if len(cols) > 0 {
		cols = utils.Merge(cols, keyName, entity.IdDbName, c.ldCfg.idColumn)
		cols = utils.Merge(cols, c.extraKey...)
	}

	if len(keyValues) > 0 {
		err = NewEntityLoader(db, func(opt *EntityLoaderConfig) {
			opt.SetIds(keyValues...).SetColumns(cols...).SetTableName(tableName).SetIdColumn(keyName)
		}).Load(&dataList)
	} else {
		bd := builder.NewSqlBuilder(dorm.GetDbType(db), dorm.GetDbSchema(db), tableName).AddColumn(cols...)
		listSql, _, err1 := bd.Build()
		if err1 != nil {
			return dataList, err1
		}
		err = dorm.RawFetch(listSql, db, &dataList)
	}

	return dataList, err
}

func (c *Cache[T, V]) Delete(db *gorm.DB, keys ...string) {
	var (
		appId = dorm.GetAppId(db)
	)
	for _, v := range keys {
		c.getStore(appId).Delete(v)
		k := concatPublishDelKey(c.cacheKey, appId, v)
		err := cache.PublishDelKey(k)
		if err != nil {
			logger.Errorf("cacheLoader publish delete key [%s] error %s", k, err)
		}
		if _, ok := c.hasAll.Load(appId); ok {
			c.hasAll.Delete(appId)
		}
	}
}
func (c *Cache[T, V]) HasAll(db *gorm.DB) bool {
	var (
		appId = dorm.GetAppId(db)
	)
	_, ok := c.hasAll.Load(appId)
	return ok
}
func (c *Cache[T, V]) SetHasAll(db *gorm.DB, hasAll bool) {
	var (
		appId = dorm.GetAppId(db)
	)
	if hasAll {
		c.hasAll.Store(appId, appId)
	} else {
		c.hasAll.Delete(appId)
	}
}
func (c *Cache[T, V]) DeleteAll(db *gorm.DB) {
	c.deleteAll(dorm.GetAppId(db))
}
func (c *Cache[T, V]) deleteAll(appIds ...string) {
	if len(appIds) < 1 {
		c.store.Range(func(key, value any) bool {
			appIds = append(appIds, key.(string))
			return false
		})
	}
	for _, v := range appIds {
		if v == "" {
			continue
		}
		c.store.Delete(v)
	}
}

func (c *Cache[T, V]) DeleteAllByAppId(appId string) {
	c.deleteAll(appId)
}

func (c *Cache[T, V]) getAll(appId string) map[string]V {
	mp := make(map[string]V)
	c.getStore(appId).Range(func(key, value any) bool {
		mp[key.(string)] = value.(V)
		return true
	})
	return mp
}
func (c *Cache[T, V]) loadExists(appId string, keyValues ...string) (exists map[string]V, notExistsKeys []string) {
	exists = make(map[string]V)
	for _, v := range keyValues {
		ad, ok := c.getStore(appId).Load(v)
		//如果传入的keyValues是副key的值，那么应该是找不到，所以通过副keyValue找主keyValue
		if !ok {
			rk := c.extraKey2Key[v]
			if rk != "" {
				ad, ok = c.getStore(appId).Load(rk)
			}
		}
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
		key := getKey(c.ldCfg.idColumn, v)
		if key == "" {
			continue
		}
		//存储主key与副key的关系
		for _, ek := range c.extraKey {
			ekVal := getKey(ek, v)
			if ekVal != "" {
				c.extraKey2Key[ekVal] = key
			}
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
