package user2app

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/loader"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/app"
	"gorm.io/gorm"
)

var (
	//存储的内容是appId -> []Setting
	allData = loader.NewCacheLoader[User2App, []User2App](constants.TableUser2App, constants.CacheAllDataUser2App).SetKeyName(entity.FromIdDbName)
)

func CleanCache(db *gorm.DB) {
	allData.DeleteAll(db)
}

// LoadDefaultAppByUserId
// TODO 添加缓存
// 如果用户登录没有指定app，那么就需要让用户登录默认的app
func LoadDefaultAppByUserId(db *gorm.DB, userId string) (ap app.App, err error) {

	dt, err := allData.LoadData(db, userId)
	if err != nil {
		return
	}
	u2aList := dt[userId]
	if len(u2aList) < 1 {
		err = errs.NewClientError("用户没有分配任何应用")
		return
	}
	for _, v := range u2aList {
		if ctype.Bool(v.IsDefault) {
			ap, err = app.LoadAppById(db, v.ToId)
			return
		}
	}
	ap, err = app.LoadAppById(db, u2aList[0].ToId)
	return
}

func LoadUser2App(db *gorm.DB, userId string, appId string) (User2App, error) {
	u2aMap, err := allData.LoadData(db, userId)
	if err != nil {
		return User2App{}, err
	}
	u2aList := u2aMap[userId]
	for _, v := range u2aList {
		if v.ToId == appId {
			return v, nil
		}
	}
	return User2App{}, errs.NewClientError(fmt.Sprintf("用户[%s]没有分配应用[%s]", userId, appId))
}

// UserIsManagerForApp
// 这个需要缓存，因为每次auth都会调
func UserIsManagerForApp(userId, appId string) bool {
	u2a, err := LoadUser2App(global.GetLocalGorm(), userId, appId)
	if err != nil {
		logger.Errorf("user is manager for app err %s", err)
	}
	return ctype.Bool(u2a.IsManager)
}

// LoadUser2AppByUserId
// 调用不频繁，可以不用缓存
func LoadUser2AppByUserId(db *gorm.DB, userIdList ...string) ([]User2App, error) {
	u2aList, err := loadUser2App(db, map[string]any{entity.FromIdDbName: userIdList})
	return u2aList, err
}
func LoadUser2AppByAppId(db *gorm.DB, appIdList ...string) ([]User2App, error) {
	u2aList, err := loadUser2App(db, map[string]any{entity.ToIdDbName: appIdList})
	return u2aList, err
}

func loadUser2App(db *gorm.DB, filters map[string]any) ([]User2App, error) {
	var (
		u2aList = make([]User2App, 0, 1)
	)
	err := dorm.Table(db, constants.TableUser2App).Where(filters).Find(&u2aList).Error
	return u2aList, err
}
