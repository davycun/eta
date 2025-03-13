package auth

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/cache"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/eta/constants"
	"gorm.io/gorm"
)

func DelAuth2RoleCache(scm, roleId string) {
	err, _ := cache.Del(constants.RedisKey(constants.Auth2RoleKey, scm, roleId))
	if err != nil {
		logger.Errorf("clean auth2role cache err %s", roleId)
	}
}

func LoadAuth2RoleByRoleId(db *gorm.DB, roleId string) (rolePerm []Auth2Role, err error) {
	var (
		scm    = dorm.GetDbSchema(db)
		dbType = dorm.GetDbType(db)
		exists bool
		cols   = []string{"id", "from_id", "to_id", "auth_type", "from_table", "to_table", "auth_table"}
	)
	err, exists = cache.Get(constants.RedisKey(constants.Auth2RoleKey, scm, roleId), &rolePerm)
	if exists || err != nil {
		return
	}

	err = db.Model(&rolePerm).Select(dorm.JoinColumns(dbType, "", cols)).Where(fmt.Sprintf(`%s = ?`, dorm.Quote(dbType, "to_id")), roleId).Find(&rolePerm).Error
	if err != nil {
		return
	}

	return rolePerm, cache.Set(constants.RedisKey(constants.Auth2RoleKey, scm, roleId), &rolePerm)
}
func LoadAuth2RoleByRoleIds(db *gorm.DB, roleIds ...string) (a2r []Auth2Role, err error) {
	for _, v := range roleIds {
		perm, err1 := LoadAuth2RoleByRoleId(db, v)
		if err1 != nil {
			return
		}
		for i, _ := range perm {
			a2r = append(a2r, perm[i])
		}
	}
	return
}
