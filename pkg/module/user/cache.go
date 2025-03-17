package user

import (
	"errors"
	"github.com/davycun/eta/pkg/common/cache"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/dao"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/dept"
	"github.com/duke-git/lancet/v2/maputil"
	"github.com/duke-git/lancet/v2/slice"
	"gorm.io/gorm"
)

var (
	NotLogin      = errors.New("user Not Login")
	InvalidAppErr = errs.NewServerError("app不可用")
)

func loadDeptIdByUserId(appDb *gorm.DB, userId string) (deptId string, err error) {
	var (
		u2dList []dept.RelationDept
	)
	deptId = userId

	err = dorm.Table(appDb, constants.TableUser2Dept).Where(map[string]any{"from_id": userId}).Find(&u2dList).Error
	if err != nil {
		return
	}
	for _, u2d := range u2dList {
		if u2d.IsMain {
			deptId = u2d.ToId
		}
	}
	return
}

func LoadUserById(db *gorm.DB, userId string) (us User, err error) {
	if userId == "" {
		return us, errors.New("user id can not be empty")
	}
	b, err := cache.Get(constants.RedisKey(constants.UserKey, userId), &us)
	if b {
		return
	}
	if err != nil {
		return us, err
	}

	err = dao.FetchById(userId, db, &us)
	if err != nil {
		return
	}

	return us, cache.Set(constants.RedisKey(constants.UserKey, us.ID), us)
}
func DelUserCache(userId ...string) {
	for _, v := range userId {
		_, err := cache.Del(constants.RedisKey(constants.UserKey, v))
		if err != nil {
			logger.Errorf("del user cache err %s", err)
		}
	}
}

func DelUserTokenByIdAndDeptId(userId string, deptId ...string) error {
	var (
		userTokenKey = constants.RedisKey(constants.UserTokenKey, userId)
		userToken    = make(map[string]TokenDept)
		toDel        = make(map[string]TokenDept)
		err          error
	)

	_, err = cache.Get(userTokenKey, &userToken)
	if err != nil {
		return err
	}
	if len(userToken) <= 0 {
		return nil
	}

	toDel = maputil.Filter(userToken, func(key string, value TokenDept) bool {
		return slice.Contain(deptId, value.DeptId)
	})
	if len(toDel) > 0 {
		for _, tk := range maputil.Keys(toDel) {
			err = DelUserToken(tk)
			if err != nil {
				return err
			}
		}
	}

	userToken = maputil.Filter(userToken, func(key string, value TokenDept) bool {
		return !slice.Contain(deptId, value.DeptId)
	})
	if len(userToken) <= 0 {
		_, err = cache.Del(userTokenKey)
		return err
	} else {
		ex, err1 := cache.TTL(userTokenKey)
		if err1 != nil {
			return err1
		}
		return cache.SetEx(userTokenKey, userToken, ex)
	}
}

// LoadDefaultUser
// 参数为什么需要db，避免默认用户第一次创建之后还在事务内，必须事务内才内查询得到
func LoadDefaultUser(db *gorm.DB) (User, error) {

	var (
		err      error
		userList = make([]User, 0, 1)
	)
	err = dorm.Table(db, constants.TableUser).
		Where(map[string]any{"account": GetRootUser().Account}).
		Limit(1).
		Find(&userList).Error
	if err != nil || len(userList) < 1 {
		return User{}, err
	}
	return userList[0], nil
}
