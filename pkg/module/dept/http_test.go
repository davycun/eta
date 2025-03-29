package dept_test

import (
	"github.com/davycun/eta/pkg/common/dorm/filter"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/http_tes"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/authorize/auth"
	"github.com/davycun/eta/pkg/module/dept"
	"github.com/davycun/eta/pkg/module/menu"
	"github.com/davycun/eta/pkg/module/user"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNotPermissionNoErr(t *testing.T) {

	//新增一个新用户
	us := user.NewTestData()
	usRs := http_tes.Create(t, "/user/create", []user.User{us})
	assert.Equal(t, 1, len(usRs))

	deptList := createDept(t)

	//获取新用户token
	_, userId, token, err := http_tes.Login(us.Account.Data, us.Password)
	assert.Nil(t, err)
	assert.NotEqual(t, "", token)
	assert.Equal(t, us.ID, userId)

	//给新用户分配菜单权限，避免403
	newMenuAuth(t, userId)

	//新用户执行
	queryDeptNewUserOrAdmin(t, deptList, token)
	//admin执行
	queryDeptNewUserOrAdmin(t, deptList, "")

}

func createDept(t *testing.T) []dept.Department {
	//创建测试部门
	deptList := dept.NewTestData()
	ds := http_tes.Create(t, "/dept/create", deptList)
	assert.Equal(t, 4, len(ds))
	return deptList
}

func newMenuAuth(t *testing.T, userId string) {

	menuList := menu.NewTestData()
	ms := http_tes.Create(t, "/menu/create", menuList)
	assert.Equal(t, 4, len(ms))

	a2rList := []auth.Auth2Role{
		{
			BaseEdgeEntity: entity.BaseEdgeEntity{
				BaseEntity: entity.BaseEntity{
					ID: global.GenerateIDStr(),
				},
				FromId: menuList[0].ID,
				ToId:   userId,
			},
			FromTable: constants.TableMenu,
			ToTable:   constants.TableUser,
			AuthTable: constants.TableMenu,
			AuthType:  auth.Read,
		},
	}

	http_tes.Create(t, "/auth2role/create", a2rList)

}

func queryDeptNewUserOrAdmin(t *testing.T, deptList []dept.Department, token string) {
	//新用户去查询部门，应该为空
	deptList, total := http_tes.Query[dept.Department](t, "/dept/list", dto.RetrieveParam{
		Filters: []filter.Filter{
			{
				LogicalOperator: filter.And,
				Column:          entity.IdDbName,
				Operator:        filter.IN,
				Value:           []string{deptList[0].ID, deptList[1].ID, deptList[2].ID, deptList[3].ID},
			},
		},
	}, func(hc *http_tes.HttpCase) {
		if token != "" {
			if hc.Headers == nil {
				hc.Headers = make(map[string]string)
			}
			hc.Headers[http_tes.HeaderAuthorization] = token
		}
	})

	//表示新用户
	if token != "" {
		assert.Equal(t, 0, len(deptList))
		assert.Equal(t, 0, int(total))
	} else {
		assert.Equal(t, 4, len(deptList))
	}

}

type CustomStruct struct {
	Name string
}

func TestSlice(t *testing.T) {
	cs := []CustomStruct{
		{Name: "a"},
		{Name: "b"},
		{Name: "c"},
	}
	logger.Infof("cs: %v", cs)

	slice.ForEach(slice.Chunk(cs, 1), func(i int, v []CustomStruct) {
		modifyData1(v)
	})
	logger.Infof("cs1: %v", cs)
	assert.Equal(t, "a", cs[0].Name)

	slice.ForEach(slice.Chunk(cs, 1), func(i int, v []CustomStruct) {
		modifyData2(v)
	})
	logger.Infof("cs2: %v", cs)
	assert.Equal(t, "a", cs[0].Name)

	modifyData1(cs)
	logger.Infof("cs2_1: %v", cs)
	assert.Equal(t, "a", cs[0].Name)

	modifyData2(cs)
	logger.Infof("cs2_2: %v", cs)
	assert.Equal(t, "modify", cs[0].Name)
}
func modifyData1(cs []CustomStruct) {
	slice.ForEach(cs, func(i int, v CustomStruct) {
		v.Name = "modify"
	})
}
func modifyData2(cs []CustomStruct) {
	slice.ForEach(cs, func(i int, v CustomStruct) {
		cs[i].Name = "modify"
	})
}
