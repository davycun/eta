package dorm

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/cache"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/golang-module/dongle"
	"gorm.io/gorm"
	"time"
)

var (
	RawFetchCache = false
	expiration    = time.Second * time.Duration(3*60*60) // 3小时
)

func RawFetch(raw string, db *gorm.DB, result any) error {

	if raw == "" {
		return nil
	}
	if RawFetchCache {
		logger.Debugf("RawFetch with cahce.")
		return rawFetchWithCache(raw, db, result)
	}
	return db.Raw(raw).Find(result).Error
}
func RawFetch2(raw string, db *gorm.DB, result any, cc bool) error {
	//if global.GetConfig().Server.RawFetchCache {
	if cc && RawFetchCache {
		logger.Debugf("RawFetch with cahce.")
		return rawFetchWithCache(raw, db, result)
	}
	return db.Raw(raw).Find(result).Error
}

func rawFetchWithCache(raw string, db *gorm.DB, result any) error {
	rdsKey := fmt.Sprintf("eta:dorm:raw_fetch:%s", dongle.Encrypt.FromString(raw).ByMd5().ToHexString())
	err, exists := cache.Get(rdsKey, result)
	if err == nil && exists {
		logger.Debugf("RawFetch with cahce data loaded.")
		//c, b := ctx.GetCurrentContext()
		//if b && c != nil {
		//	c.Set(ctx.FetchCacheContextKey, rdsKey)
		//}
		return nil
	}
	if err != nil {
		return err
	}
	err = db.Raw(raw).Find(result).Error
	if err != nil {
		return err
	}
	err = cache.SetEx(rdsKey, &result, expiration)
	return err
}
