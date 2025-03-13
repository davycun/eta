package dict

import (
	"github.com/davycun/eta/pkg/core/loader"
	"github.com/davycun/eta/pkg/eta/constants"
	"gorm.io/gorm"
)

var (
	// DataCache
	//存储的内容是appId -> []Setting
	DataCache = loader.NewCacheLoader[Dictionary, Dictionary](constants.TableDictionary, constants.CacheAllDataDictionary).SetKeyName("id")
)

func DelCache(db *gorm.DB, dataList ...Dictionary) {
	for _, v := range dataList {
		if v.ID == "" {
			continue
		}
		DataCache.Delete(db, v.ID)
	}
}

func LoadAllDictionary(db *gorm.DB) (dtMap map[string]Dictionary, err error) {
	return DataCache.LoadAll(db)
}
