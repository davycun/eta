package user_srv

import (
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/module/app"
	"github.com/davycun/eta/pkg/module/user"
)

func GetUserProp(u user.User) ctype.Map {
	return ctype.Map{
		"id":         u.ID,
		"updated_at": u.UpdatedAt,
		"account":    u.Account,
		"avatar":     u.Avatar,
		"category":   u.Category,
		"name":       u.Name,
		"sex":        u.Sex,
		"phone":      u.Phone,
		"email":      u.Email,
		//"is_manager":   u.IsManager,
		"valid":        u.Valid,
		"job_title":    u.JobTitle,
		"user2dept":    u.User2Dept,
		"current_dept": u.CurrentDept,
	}
}
func GetAppProp(ap app.App) ctype.Map {
	return ctype.Map{
		"id":      ap.ID,
		"name":    ap.Name,
		"company": ap.Company,
		"logo":    ap.Logo,
		"slogan":  ap.Slogan,
	}
}
