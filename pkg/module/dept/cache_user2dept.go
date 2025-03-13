package dept

import (
	"github.com/davycun/eta/pkg/common/cache"
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/loader"
	"github.com/davycun/eta/pkg/eta/constants"
)

func DelUser2DeptCache(userId ...string) {
	for _, v := range userId {
		err, _ := cache.Del(constants.RedisKey(constants.User2DeptCacheKey, v))
		if err != nil {
			logger.Errorf("del user2dept cache err %s", err)
		}
	}
}

func LoadUser2DeptByUserId(c *ctx.Context, userId string) ([]RelationDept, error) {

	var (
		u2d []RelationDept
		db  = c.GetAppGorm()
	)

	err, b := cache.Get(constants.RedisKey(constants.User2DeptCacheKey, userId), &u2d)

	if err != nil {
		return u2d, err
	}
	if b {
		if len(u2d) < 1 {
			u2d = append(u2d, GetDefaultUser2Dept(userId, c.GetContextUserName()))
		} else if c.GetContextIsManager() {
			u2d = append(u2d, GetDefaultUser2Dept(userId, c.GetContextUserName()))
		}
		return u2d, err
	}

	ld := loader.NewRelationEntityLoader[Department, RelationDept](db, constants.TableUser2Dept, constants.TableDept)
	ld.AddRelationColumns(DefaultRelationDeptColumns...).AddEntityColumns(DefaultColumns...)
	toMap, err := ld.LoadToMap(userId)
	if err != nil || len(toMap) < 1 {
		return u2d, err
	}
	u2d = toMap[userId]

	err = cache.Set(constants.RedisKey(constants.User2DeptCacheKey, userId), &u2d)

	if len(u2d) < 1 {
		u2d = append(u2d, GetDefaultUser2Dept(userId, c.GetContextUserName()))
	} else if c.GetContextIsManager() {
		u2d = append(u2d, GetDefaultUser2Dept(userId, c.GetContextUserName()))
	}

	return u2d, err

}

func LoadUser2DeptByUserIdDeptId(c *ctx.Context, userId string, deptId string) (RelationDept, error) {
	u2dList, err := LoadUser2DeptByUserId(c, userId)

	if err != nil {
		return RelationDept{}, err
	}
	for _, v := range u2dList {
		if v.ToId == deptId {
			return v, nil
		}
	}
	return RelationDept{}, errs.NotFound
}
