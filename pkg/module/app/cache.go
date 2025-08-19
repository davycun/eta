package app

import (
	"errors"
	"fmt"
	"github.com/davycun/eta/pkg/common/cache"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/dao"
	"github.com/davycun/eta/pkg/eta/constants"
	"gorm.io/gorm"
)

var (
	InvalidAppErr = errs.NewServerError("app不可用")
)

func LoadAppById(db *gorm.DB, id string) (ap App, err error) {
	if id == "" {
		return ap, errors.New("app id can not be empty")
	}
	b, err := cache.Get(constants.RedisKey(constants.AppKey, id), &ap)
	if b && ctype.Bool(ap.Valid) {
		return ap, nil
	}
	if err != nil {
		return ap, err
	}

	err = dao.FetchById(id, db, &ap)
	if err != nil {
		return ap, err
	}
	if !ctype.Bool(ap.Valid) {
		return ap, InvalidAppErr
	}

	return ap, cache.Set(constants.RedisKey(constants.AppKey, ap.ID), ap)
}
func DelAppCache(appId ...string) {
	for _, v := range appId {
		_, err := cache.Del(constants.RedisKey(constants.AppKey, v))
		if err != nil {
			logger.Errorf("del app cache err %s", err)
		}
	}
	DelDefaultAppCache()
}

func DelDefaultAppCache() {
	defaultAppId = ""
}

var (
	defaultAppId = ""
)

// LoadDefaultApp
// 需要db作为参数，避免第一次创建默认app之后，还没有提交事务，但是再某个环节又需要获取DefaultApp，这个时候就需要传入TxDB才能查询得到
func LoadDefaultApp(db *gorm.DB) (ap App, err error) {
	if db == nil {
		db = global.GetLocalGorm()
	}
	if defaultAppId != "" {
		//缓存一下，本地远程调试的时候，可以改下本地redis 配置而不必修改数据库配置
		ap, err = LoadAppById(global.GetLocalGorm(), defaultAppId)
		if ap.ID != "" && ctype.Bool(ap.Valid) {
			return
		}
	}

	var (
		appList = make([]App, 0, 1)
		dbType  = dorm.GetDbType(db)
	)
	err = db.Model(&appList).
		Where(map[string]any{"valid": true, "is_default": true}).
		Order(fmt.Sprintf(`%s asc`, dorm.Quote(dbType, "id"))).
		Limit(1).Find(&appList).Error

	if err != nil {
		return
	}

	if len(appList) < 1 {
		err = errs.NewServerError("没有可用的app")
		return
	}
	ap = appList[0]
	return
}

func LoadAllApp() (apps []App, err error) {
	db := global.GetLocalGorm()
	err = db.Model(&apps).Where(&App{Valid: ctype.Boolean{Data: true, Valid: true}}).Find(&apps).Error
	return
}

// LoadAppIdBySchema
// TODO 其实这个不一定完全准确的，因为可能跨DB的schema相同，当然如果系统自动生成的话是不会相同的，除非创建app的时候人为指定了
func LoadAppIdBySchema(scm string) string {
	apps, err := LoadAllApp()
	if err != nil {
		logger.Errorf("LoadAppIdBySchema err %s", err)
	}
	for _, v := range apps {
		if v.Database.Schema == scm {
			return v.ID
		}
	}
	return ""
}
