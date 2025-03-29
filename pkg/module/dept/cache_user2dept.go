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
		_, err := cache.Del(constants.RedisKey(constants.User2DeptCacheKey, v))
		if err != nil {
			logger.Errorf("del user2dept cache err %s", err)
		}
	}
}

func LoadUser2DeptByUserId(c *ctx.Context, userId string) ([]RelationDept, error) {

	var (
		u2d = make([]RelationDept, 0, 1)
		db  = c.GetAppGorm()
	)

	b, err := cache.Get(constants.RedisKey(constants.User2DeptCacheKey, userId), &u2d)

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
	if err != nil {
		return u2d, err
	}

	if x, ok := toMap[userId]; ok {
		u2d = x
	}
	//有可能设置一个空的
	//缓存的设置必须放在后面这个if之前，因为不需要缓存虚拟部门（采用用户ID）的情况
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
