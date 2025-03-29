package auth_test

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/http_tes"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/authorize/auth"
	"github.com/davycun/eta/pkg/module/authorize/role"
	"github.com/davycun/eta/pkg/module/dept"
	"github.com/davycun/eta/pkg/module/menu"
	"github.com/davycun/eta/pkg/module/user"
	"github.com/davycun/eta/pkg/module/user/user2dept"
	"github.com/davycun/eta/pkg/module/user/user2role"
	"github.com/stretchr/testify/assert"
	"slices"
	"testing"
)

func TestAuth2Role(t *testing.T) {

	deptList := dept.NewTestData()
	roleList := role.NewTestData()
	us := user.NewTestData()
	menuList := menu.NewTestData()
	http_tes.Create(t, "/dept/create", deptList)
	http_tes.Create(t, "/role/create", roleList)
	http_tes.Create(t, "/user/create", []user.User{us})
	http_tes.Create(t, "/menu/create", menuList)

	user2RoleList := getUser2Role(us, roleList)
	user2DeptList := getUser2Dept(us, deptList, roleList) //deptList[0] 是主部门，默认会选择
	auth2RoleList := getAuth2Role(menuList, roleList)
	http_tes.Create(t, "/user2role/create", user2RoleList)
	http_tes.Create(t, "/user2dept/create", user2DeptList)
	http_tes.Create(t, "/auth2role/create", auth2RoleList)

	c := ctx.NewContextWithUserId(us.ID)
	//这里多了一个用户的默认基础角色 153035491560656898
	roleRsIds, err := auth.LoadUserRoleIdsByDeptId(c.GetAppGorm(), us.ID, deptList[2].ID)
	assert.Nil(t, err)
	assert.True(t, utils.ContainAny(roleRsIds, roleList[0].ID))
	assert.False(t, utils.ContainAny(roleRsIds, roleList[1].ID))

	roleRsIds, err = auth.LoadUserRoleIdsByDeptId(c.GetAppGorm(), us.ID, deptList[3].ID)
	assert.Nil(t, err)
	assert.False(t, utils.ContainAny(roleRsIds, roleList[0].ID))
	assert.True(t, utils.ContainAny(roleRsIds, roleList[1].ID))

	_, userId, token, err := http_tes.Login(us.Account.Data, us.Password)
	assert.Nil(t, err)
	assert.Equal(t, us.ID, userId)

	//当前部门下的角色菜单，默认主部门是deptList[0]
	menuRs, _ := http_tes.Query[menu.Menu](t, "/menu/my_menu", dto.RetrieveParam{UseCurDeptAuth: true}, func(hc *http_tes.HttpCase) {
		hc.Headers = map[string]string{
			http_tes.HeaderAuthorization: token,
		}
	})
	idx := slices.IndexFunc(menuRs, func(m menu.Menu) bool {
		if m.ID == menuList[3].ID {
			return true
		}
		return false
	})
	assert.Equal(t, -1, idx)

	//切换部门后，my_menu会变化
	http_tes.Modify[ctype.Map](t, fmt.Sprintf("/user/set_current_dept?current_dept=%s", deptList[3].ID), dto.ModifyParam{}, func(hc *http_tes.HttpCase) {
		hc.Headers = map[string]string{
			http_tes.HeaderAuthorization: token,
		}
	})
	menuRs, _ = http_tes.Query[menu.Menu](t, "/menu/my_menu", dto.RetrieveParam{UseCurDeptAuth: true}, func(hc *http_tes.HttpCase) {
		hc.Headers = map[string]string{
			http_tes.HeaderAuthorization: token,
		}
	})
	idx = slices.IndexFunc(menuRs, func(m menu.Menu) bool {
		if utils.ContainAny([]string{menuList[0].ID, menuList[1].ID, menuList[2].ID}, m.ID) {
			return true
		}
		return false
	})
	assert.Equal(t, -1, idx)

}

func getUser2Dept(us user.User, deptList []dept.Department, roleList []role.Role) []user2dept.User2Dept {
	return []user2dept.User2Dept{
		{
			BaseEdgeEntity: entity.BaseEdgeEntity{
				BaseEntity: entity.BaseEntity{
					ID: global.GenerateIDStr(),
				},
				FromId: us.ID,
				ToId:   deptList[2].ID,
			},
			IsMain:  ctype.NewBoolean(true, true),
			RoleIds: ctype.NewStringArrayPrt(roleList[0].ID),
		},
		{
			BaseEdgeEntity: entity.BaseEdgeEntity{
				BaseEntity: entity.BaseEntity{
					ID: global.GenerateIDStr(),
				},
				FromId: us.ID,
				ToId:   deptList[3].ID,
			},
			RoleIds: ctype.NewStringArrayPrt(roleList[1].ID),
		},
	}
}
func getUser2Role(us user.User, roleList []role.Role) []user2role.User2Role {
	return []user2role.User2Role{
		{
			BaseEdgeEntity: entity.BaseEdgeEntity{
				BaseEntity: entity.BaseEntity{
					ID: global.GenerateIDStr(),
				},
				FromId: us.ID,
				ToId:   roleList[0].ID,
			},
		},
		{
			BaseEdgeEntity: entity.BaseEdgeEntity{
				BaseEntity: entity.BaseEntity{
					ID: global.GenerateIDStr(),
				},
				FromId: us.ID,
				ToId:   roleList[1].ID,
			},
		},
	}
}

func getAuth2Role(menuList []menu.Menu, roleList []role.Role) []auth.Auth2Role {
	return []auth.Auth2Role{
		{
			BaseEdgeEntity: entity.BaseEdgeEntity{
				BaseEntity: entity.BaseEntity{
					ID: global.GenerateIDStr(),
				},
				FromId: menuList[0].ID,
				ToId:   roleList[0].ID,
			},
			FromTable: constants.TableMenu,
			ToTable:   constants.TableRole,
			AuthTable: constants.TableMenu,
			AuthType:  auth.Read,
		},
		{
			BaseEdgeEntity: entity.BaseEdgeEntity{
				BaseEntity: entity.BaseEntity{
					ID: global.GenerateIDStr(),
				},
				FromId: menuList[1].ID,
				ToId:   roleList[0].ID,
			},
			FromTable: constants.TableMenu,
			ToTable:   constants.TableRole,
			AuthTable: constants.TableMenu,
			AuthType:  auth.Read,
		},
		{
			BaseEdgeEntity: entity.BaseEdgeEntity{
				BaseEntity: entity.BaseEntity{
					ID: global.GenerateIDStr(),
				},
				FromId: menuList[2].ID,
				ToId:   roleList[0].ID,
			},
			FromTable: constants.TableMenu,
			ToTable:   constants.TableRole,
			AuthTable: constants.TableMenu,
			AuthType:  auth.Read,
		},
		{
			BaseEdgeEntity: entity.BaseEdgeEntity{
				BaseEntity: entity.BaseEntity{
					ID: global.GenerateIDStr(),
				},
				FromId: menuList[3].ID,
				ToId:   roleList[1].ID,
			},
			FromTable: constants.TableMenu,
			ToTable:   constants.TableRole,
			AuthTable: constants.TableMenu,
			AuthType:  auth.Read,
		},
	}
}
