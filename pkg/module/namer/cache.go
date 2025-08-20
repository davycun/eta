package namer

import (
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/loader"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/duke-git/lancet/v2/maputil"
)

var (
	//存储的内容是appId -> []IdName

	deptCol = []string{"id", "name", "namespace", "parent_id"}
	userCol = []string{"id", "name", "account", "category"}

	userIdNameCache = loader.NewCacheLoader[IdName, IdName](constants.TableUser, constants.CacheAllDataIdNameUser, userCol...).SetKeyName(entity.IdDbName)
	deptIdNameCache = loader.NewCacheLoader[IdName, IdName](constants.TableDept, constants.CacheAllDataIdNameDept, deptCol...).SetKeyName(entity.IdDbName)
)

func LoadAllIdName(c *ctx.Context) (mnMap map[string]IdName, err error) {

	userAll, err := userIdNameCache.LoadAll(global.GetLocalGorm())
	if err != nil {
		return nil, err
	}

	deptAll, err := deptIdNameCache.LoadAll(c.GetAppGorm())

	for k, v := range deptAll {
		userAll[k] = v
	}
	return userAll, err
}

func LoadByIds(c *ctx.Context, ids []string) (mnMap map[string]IdName, err error) {
	mnMap, err = LoadAllIdName(c)
	if err != nil {
		return
	}
	mnMap = maputil.FilterByKeys(mnMap, ids)
	return
}

func DelIdNameCacheByContext(c *ctx.Context) {
	userIdNameCache.DeleteAllAppData(global.GetLocalGorm())
	if c.GetAppGorm() != nil {
		deptIdNameCache.DeleteAllAppData(c.GetAppGorm())
	}
}
