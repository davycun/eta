package auth

import (
	"github.com/davycun/eta/pkg/eta/constants"
	"gorm.io/gorm"
)

func FetchAuth2RolePermissionByRoleIds(db *gorm.DB, roleIds ...string) (rs []Auth2Role, err error) {

	r, err := LoadAuth2RoleByRoleIds(db, roleIds...)
	for i, v := range r {
		if v.FromTable == constants.TablePermission {
			rs = append(rs, r[i])
		}
	}
	return
}
func FetchUserAuth2Role(db *gorm.DB, fromTable, userId, authTable string, authType Type) (a2r []Auth2Role, err error) {

	var (
		roleIds     []string
		auth2roleRs []Auth2Role
	)

	roleIds, err = LoadUserAllRoleIds(db, userId)
	if len(roleIds) < 1 || err != nil {
		return
	}
	auth2roleRs, err = LoadAuth2RoleByRoleIds(db, roleIds...)
	if err != nil || len(auth2roleRs) < 1 {
		return
	}

	for i, v := range auth2roleRs {
		if v.FromTable == fromTable && v.AuthTable == authTable && v.AuthType&authType == authType {
			a2r = append(a2r, auth2roleRs[i])
		}
	}
	return
}
func FetchUserAuth2RoleAll(db *gorm.DB, fromTable, userId, authTable string) (a2r []Auth2Role, err error) {

	var (
		roleIds     []string
		auth2roleRs []Auth2Role
	)

	roleIds, err = LoadUserAllRoleIds(db, userId)
	if len(roleIds) < 1 || err != nil {
		return
	}
	auth2roleRs, err = LoadAuth2RoleByRoleIds(db, roleIds...)
	if err != nil || len(auth2roleRs) < 1 {
		return
	}

	for i, v := range auth2roleRs {
		if v.FromTable == fromTable && v.AuthTable == authTable {
			a2r = append(a2r, auth2roleRs[i])
		}
	}
	return
}
