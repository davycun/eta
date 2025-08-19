package setting

import (
	"errors"
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/loader"
	"github.com/davycun/eta/pkg/eta/constants"
	jsoniter "github.com/json-iterator/go"
	"gorm.io/gorm"
	"reflect"
)

var (
	//存储的内容是appId -> []Setting
	localDbAllData = loader.NewCacheLoader[Setting, Setting](constants.TableSetting, constants.CacheAllDataConfig).SetKeyName(entity.IdDbName)
	appDbAllData   = loader.NewCacheLoader[Setting, Setting](constants.TableSetting, constants.CacheAllDataConfigSetting).SetKeyName(entity.IdDbName)
)

func loadAllConfig(db *gorm.DB) (dtMap map[string]Setting, err error) {
	if db == nil {
		db = global.GetLocalGorm()
	}
	if isAppDb(db) {
		return appDbAllData.LoadAll(db)
	}
	return localDbAllData.LoadAll(db)
}

func HasCacheAll(db *gorm.DB, all bool) {
	if isAppDb(db) {
		appDbAllData.SetHasAll(db, all)
	}
	appDbAllData.SetHasAll(db, all)
}

func DelCache(db *gorm.DB, dataList ...reflect.Value) {
	for _, v := range dataList {
		id := entity.GetString(v, entity.IdDbName)
		if id == "" {
			continue
		}
		if isAppDb(db) {
			appDbAllData.Delete(db, id)
		}
		localDbAllData.Delete(db, id)
	}
}

func isAppDb(db *gorm.DB) bool {
	if db == nil {
		return false
	}
	var (
		lcDb = global.GetLocalGorm()
	)

	return !(dorm.GetDbHost(db) == dorm.GetDbHost(lcDb) &&
		dorm.GetDbType(db) == dorm.GetDbType(lcDb) &&
		dorm.GetDbPort(db) == dorm.GetDbPort(lcDb) &&
		dorm.GetDbSchema(db) == dorm.GetDbSchema(lcDb))
}

func unmarshal(db *gorm.DB, category, name string, rs any) (bool, error) {
	cfg, exists := GetSetting(db, category, name)
	if !exists || !ctype.IsValid(cfg.Content) {
		return false, nil
	}
	ms, err := jsoniter.Marshal(cfg.Content)
	if err != nil {
		return false, err
	}
	err = jsoniter.Unmarshal(ms, rs)
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetSetting(db *gorm.DB, category, name string) (st Setting, exists bool) {
	alSet, err := loadAllConfig(db)
	if err != nil {
		logger.Errorf("load all config err %s", err)
		return
	}
	for _, v := range alSet {
		if v.Category == category && name == v.Name {
			return v, true
		}
	}
	return
}

// GetConfig
// 根据category和name获取配置信息，
// 1. 根据传入的db从对应的schema获取配置，如果db是appDB并且获取不到配置，那么会从localDB获取配置
// 2. 如果db是localDB并且获取不到配置，那么会从defaultSettingMap中获取配置
func GetConfig[T any](db *gorm.DB, category, name string) (T, error) {
	var (
		cfg     T
		errList = make([]error, 0)
	)
	b, err := unmarshal(db, category, name, &cfg)
	if err != nil {
		errList = append(errList, errors.New(fmt.Sprintf("load config[category:%s,name:%s] err %s", category, name, err)))
	}
	if b {
		return cfg, errors.Join(errList...)
	}
	//如果从appDb里面找不到，那就从localDB里面找
	if isAppDb(db) {
		b, err = unmarshal(global.GetLocalGorm(), category, name, &cfg)
		if err != nil {
			errList = append(errList, errors.New(fmt.Sprintf("load config[category:%s,name:%s] err %s", category, name, err)))
		}
		if b {
			return cfg, errors.Join(errList...)
		}
	}

	cfg = GetDefault[T](category, name)
	return cfg, errors.Join(errList...)
}
