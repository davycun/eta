package auth

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/cache"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/eta/constants"
	"gorm.io/gorm"
)

func DelPermissionCache(id string) {
	err, _ := cache.Del(constants.RedisKey(constants.PermissionCacheKey, id))
	if err != nil {
		logger.Errorf("clean permission cache err %s", id)
	}
}
func LoadPermissionById(db *gorm.DB, id string) (perm Permission, err error) {
	//var dp Permission

	err, b := cache.Get(constants.RedisKey(constants.PermissionCacheKey, id), &perm)
	if err != nil {
		return
	}
	if b {
		return
	}

	err = db.Model(&perm).Where(fmt.Sprintf(`"id" = ?`), id).Find(&perm).Error
	if err != nil {
		return
	}

	return perm, cache.Set(constants.RedisKey(constants.PermissionCacheKey, id), &perm)
}
