package user

import (
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/user/userkey"
)

const (
	RootUserAccount  = "root"
	RootUserPassword = "root@123"
)

var (
	rootUser = User{
		Account:  ctype.NewStringPrt(RootUserAccount),
		Name:     "管理员",
		Password: RootUserPassword,
		Valid:    ctype.Boolean{Valid: true, Data: true},
	}
	openApiUser = User{
		Account:  ctype.NewStringPrt("openapi_admin"),
		Name:     "OpenApi管理员",
		Password: "admin@123",
		Category: constants.UserTypeOpenApi,
		UserKey: []userkey.UserKey{
			{
				AccessKey:    ctype.NewStringPrt(constants.DefaultOpenApiAk),
				AccessSecure: ctype.NewStringPrt(constants.DefaultOpenApiSk),
				FixedToken:   ctype.NewStringPrt(constants.DefaultOpenApiFixedToken),
			},
		},
		Valid:      ctype.Boolean{Valid: true, Data: true},
		BaseEntity: entity.BaseEntity{ID: global.GenerateIDStr()},
	}
)

func GetRootUser() User {
	return rootUser
}
