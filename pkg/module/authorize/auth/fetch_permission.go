package auth

import (
	"github.com/davycun/eta/pkg/eta/constants"
	"gorm.io/gorm"
)

func FetchUserPermission(db *gorm.DB, userId string, authTable string, authType Type) (perms []Permission, err error) {

	var (
		roleIds     []string
		auth2roleRs []Auth2Role
	)

	roleIds, err = LoadUserAllRoleIds(db, userId)
	if len(roleIds) < 1 || err != nil {
		return
	}
	auth2roleRs, err = FetchAuth2RolePermissionByRoleIds(db, roleIds...)
	if err != nil || len(auth2roleRs) < 1 {
		return
	}

	for _, v := range auth2roleRs {
		if v.FromTable == constants.TablePermission && v.AuthTable == authTable && v.AuthType&authType == authType {
			pm, err1 := LoadPermissionById(db, v.FromId)
			if err1 != nil {
				err = err1
				return
			}
			if pm.ID != "" {
				perms = append(perms, pm)
			}
		}
	}
	return
}

func FetchRolePermission(db *gorm.DB, roleId string, authTable string, authType Type) (perms []Permission, err error) {

	var (
		auth2roleRs []Auth2Role
	)

	auth2roleRs, err = FetchAuth2RolePermissionByRoleIds(db, roleId)
	if err != nil || len(auth2roleRs) < 1 {
		return
	}

	for _, v := range auth2roleRs {
		if v.FromTable == constants.TablePermission && v.AuthTable == authTable && v.AuthType&authType == authType {
			pm, err1 := LoadPermissionById(db, v.FromId)
			if err1 != nil {
				err = err1
				return
			}
			if pm.ID != "" {
				perms = append(perms, pm)
			}
		}
	}
	return
}
