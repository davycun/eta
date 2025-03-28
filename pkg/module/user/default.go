package user

import (
	"github.com/davycun/eta/pkg/common/dorm/ctype"
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
		UserKey: []userkey.UserKey{
			{
				AccessKey:    ctype.NewStringPrt(constants.DefaultOpenApiAk),
				AccessSecure: ctype.NewStringPrt(constants.DefaultOpenApiSk),
				FixedToken:   ctype.NewStringPrt(constants.DefaultOpenApiFixedToken),
			},
		},
	}
)

func GetRootUser() User {
	return rootUser
}
