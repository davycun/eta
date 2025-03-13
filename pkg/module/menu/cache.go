package menu

import (
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/loader"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/authorize/auth"
	"gorm.io/gorm"
)

var (
	//存储的内容是appId -> []Menu
	DataCache = loader.NewCacheLoader[Menu, Menu](constants.TableMenu, constants.CacheAllDataMenu).SetKeyName("id")
)

func LoadAllMenu(db *gorm.DB) (dtMap map[string]Menu, err error) {
	return DataCache.LoadAll(db)
}
func LoadMenuById(db *gorm.DB, id string) (Menu, error) {
	data, err := DataCache.LoadData(db, id)
	lb := data[id]
	return lb, err
}

func DelCache(db *gorm.DB, dataList ...Menu) {
	for _, v := range dataList {
		if v.ID == "" {
			continue
		}
		DataCache.Delete(db, v.ID)
	}
}

func LoadMenuByUserId(c *ctx.Context) (mn []Menu, err error) {
	mns, err := LoadAllMenu(c.GetAppGorm())
	if err != nil {
		return mn, err
	}
	for _, v := range mns {
		mn = append(mn, v)
	}
	return FilterMenuByUserId(c, c.GetContextUserId(), mn)
}
func FilterMenuByUserId(c *ctx.Context, userId string, menuList []Menu) (mn []Menu, err error) {
	auth2r, err := auth.FetchUserAuth2Role(c.GetAppGorm(), constants.TableMenu, userId, constants.TableMenu, auth.Read)
	if err != nil {
		return mn, err
	}
	if len(auth2r) < 1 {
		return mn, err
	}
	mns := make(map[string]Menu)
	for _, v := range menuList {
		mns[v.ID] = v
	}
	tmp := make(map[string]string)
	for _, v := range auth2r {
		s := tmp[v.FromId]
		if s != "" {
			continue
		}
		m, ok := mns[v.FromId]
		if ok {
			mn = append(mn, m)
			tmp[v.FromId] = v.FromId
		}
	}
	return mn, err
}
func FilterMenuByRoleIds(c *ctx.Context, menuList []Menu, roleIds ...string) (mn []Menu, err error) {

	auth2r, err := auth.LoadAuth2RoleByRoleIds(c.GetAppGorm(), roleIds...)
	if err != nil {
		return mn, err
	}
	if len(auth2r) < 1 {
		return mn, err
	}
	mns := make(map[string]Menu)
	for _, v := range menuList {
		mns[v.ID] = v
	}
	//临时保存已经加入过的菜单
	tmp := make(map[string]string)

	for _, v := range auth2r {
		if v.FromTable != constants.TableMenu || v.AuthTable != constants.TableMenu || v.AuthType != auth.Read {
			continue
		}
		s := tmp[v.FromId]
		if s != "" {
			continue
		}
		m, ok := mns[v.FromId]
		if ok {
			mn = append(mn, m)
			tmp[v.FromId] = v.FromId
		}
	}
	return mn, err
}

func GetParentIds(db *gorm.DB, id string) []string {
	ids := make([]string, 0, 1)
	if id == "" {
		return ids
	}
	dt, err := LoadMenuById(db, id)
	if err != nil {
		logger.Infof("menu get parentId err %s", err)
		return ids
	}
	if dt.ID != "" && dt.ParentId != "" {
		ids = append(ids, dt.ParentId)
		pIds := GetParentIds(db, dt.ParentId)
		ids = utils.AppendNoEmpty(ids, pIds...)
	}
	return ids
}
