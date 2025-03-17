package dept

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/cache"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/eta/constants"
	"gorm.io/gorm"
)

func LoadDeptById(db *gorm.DB, id string) Department {
	var dp Department

	b, _ := cache.Get(constants.RedisKey(constants.DeptCacheKey, id), &dp)
	if b {
		return dp
	}

	err := db.Model(&dp).Where(fmt.Sprintf(`"id" = ?`), id).Find(&dp).Error
	if err != nil {
		logger.Errorf("can not find department by id %s", id)
		return dp
	}
	err = cache.Set(constants.RedisKey(constants.DeptCacheKey, id), &dp)
	if err != nil {
		logger.Errorf("set department to cache err %s", id)
		return dp
	}
	return dp
}

func DelDeptCache(deptIds ...string) {
	for _, v := range deptIds {
		_, err := cache.Del(constants.RedisKey(constants.DeptCacheKey, v))
		if err != nil {
			logger.Errorf("clean department cache err %s", v)
		}
	}

}

func DelDeptAndUser2DeptCache(db *gorm.DB, deptIds ...string) error {
	if len(deptIds) < 1 {
		return nil
	}
	for _, v := range deptIds {
		DelDeptCache(v)
	}
	var (
		uId []string
	)
	err := dorm.Table(db, constants.TableUser2Dept).
		Select(`"from_id"`).
		Where(`"to_id" in ?`, deptIds).
		Find(&uId).Error
	if err != nil {
		return err
	}

	DelUser2DeptCache(uId...)

	return nil
}
