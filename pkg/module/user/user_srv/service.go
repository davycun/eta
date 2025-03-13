package user_srv

import (
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/service"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/dept"
	"github.com/davycun/eta/pkg/module/setting"
	"github.com/davycun/eta/pkg/module/user"
	"github.com/davycun/eta/pkg/module/user/login/captcha"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"time"
)

type Service struct {
	service.DefaultService
}

func (s *Service) ChangePassword(param *user.ModifyPasswordParam, result *dto.Result) error {

	var (
		err    error
		db     = global.GetLocalGorm()
		userId = s.GetContext().GetContextUserId()
		us     = user.User{}
		cfg, _ = setting.GetLoginConfig(s.GetContext().GetAppGorm())
	)

	err = caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			if !cfg.PasswordValid(param.NewPassword) {
				return errs.NewClientError("密码格式不正确")
			}
			return nil
		}).
		Call(func(cl *caller.Caller) error {
			us, err = user.LoadUserById(db, userId)
			return err
		}).
		Call(func(cl *caller.Caller) error {
			if us.ID == "" {
				return errs.NewClientError("用户不存在")
			}
			return nil
		}).
		Call(func(cl *caller.Caller) error {
			if bcrypt.CompareHashAndPassword(utils.StringToBytes(us.Password), utils.StringToBytes(param.CurrentPassword)) != nil {
				return errs.NewClientError("原密码错误")
			}
			return nil
		}).
		Call(func(cl *caller.Caller) error {
			// 校验密码
			pattern := cfg.PwdValidateReg
			if pattern != "" {
				if match, _ := regexp.MatchString(pattern, param.NewPassword); !match {
					return errs.NewClientError("密码格式不正确")
				}
			}
			return nil
		}).
		Call(func(cl *caller.Caller) error {

			passwd, err1 := bcrypt.GenerateFromPassword(utils.StringToBytes(param.NewPassword), bcrypt.DefaultCost)
			if err1 != nil {
				return err1
			}
			return global.GetLocalGorm().Model(&us).Where(map[string]any{"id": us.ID}).Updates(map[string]interface{}{
				"password":        utils.BytesToString(passwd),
				"last_update_pwd": &ctype.LocalTime{Data: time.Now(), Valid: true},
			}).Error
		}).Err

	return err
}
func (s *Service) SetCurrentDept(args *dto.Param, result *dto.Result) error {

	var (
		c         = s.GetContext()
		userId    = c.GetContextUserId()
		curDeptId = c.GetGinContext().Query("current_dept")
		tk, err   = user.LoadTokenByToken(ctx.GetToken(s.GetContext()))
	)
	if err != nil {
		return err
	}

	if curDeptId == "" {
		return errs.NewClientError("query param current_dept is empty")
	}
	if userId == "" {
		return user.NotLogin
	}
	relDept, err := dept.LoadUser2DeptByUserId(c, userId)
	if err != nil {
		return err
	}

	for _, v := range relDept {
		if v.ToId == curDeptId {
			tk.DeptId = curDeptId
			return user.StoreToken(tk)
		}
	}
	//可以设置自己的ID为当前部门
	if curDeptId == userId {
		tk.DeptId = curDeptId
		return user.StoreToken(tk)
	}

	return errs.NewClientError("当前用户不在 指定的 current_dept 部门")
}
func (s *Service) Current(args *dto.Param, result *dto.Result) error {

	var (
		c       = s.GetContext()
		tk      user.TokenInfo
		u2dList []dept.RelationDept
		curDept dept.RelationDept
	)

	tk, ok := user.GetContextToken(c)

	if !ok || tk.Token == "" {
		return user.NotLogin
	}
	u, ok := user.GetContextUser(c)
	if !ok {
		return user.NotLogin
	}

	curDept, err := dept.LoadUser2DeptByUserIdDeptId(c, tk.UserId, tk.DeptId)
	if err != nil {
		return err
	}

	u2dList, err = dept.LoadUser2DeptByUserId(c, tk.UserId)
	if err != nil {
		return err
	}

	u.User2Dept = u2dList
	u.CurrentDept = curDept
	result.Data = GetUserProp(*u)
	return err
}

// ChangePhone
// 如果以前的用户手机号为空表示是设置手机号
// 如果用户以前的手机号不为空，那么验证码就是验证老手机的
func (s *Service) ChangePhone(args *user.ModifyPhoneParam, result *dto.Result) error {

	var (
		db       = global.GetLocalGorm()
		us, err  = user.LoadUserById(db, s.GetContext().GetContextUserId())
		oldPhone = ctype.ToString(us.Phone)
		userId   = s.GetContext().GetContextUserId()
	)

	if oldPhone == "" {
		if !captcha.Verify(captcha.Captcha{Code: args.Code, Phone: args.NewPhone}) {
			return errs.NewClientError("验证码错误")
		}
	} else {
		if !captcha.Verify(captcha.Captcha{Code: args.Code, Phone: oldPhone}) {
			return errs.NewClientError("验证码错误")
		}
	}

	err = dorm.Table(db, constants.TableUser).Where(map[string]any{"id": userId}).Updates(map[string]any{"phone": args.NewPhone}).Error
	if err != nil {
		return err
	}
	user.DelUserCache(userId)
	return nil
}
func (s *Service) ResetPassword(param *user.ResetPasswordParam, result *dto.Result) error {
	var (
		c  = s.GetContext()
		db = global.GetLocalGorm()
	)

	if !c.GetContextIsManager() {
		return errs.NewServerError("只有管理员才能重置密码")
	}

	// 校验密码
	cfg, _ := setting.GetLoginConfig(s.GetContext().GetAppGorm())
	if !cfg.PasswordValid(param.NewPassword) {
		return errs.NewClientError("密码格式不正确")
	}

	// 更新新密码
	passwd, err := bcrypt.GenerateFromPassword(utils.StringToBytes(param.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = dorm.Table(db, constants.TableUser).
		Where(map[string]any{"id": param.UserIdList}).
		Updates(map[string]any{"password": utils.BytesToString(passwd), "last_update_pwd": &ctype.LocalTime{Data: time.Now(), Valid: true}}).Error
	if err != nil {
		return err
	}
	user.DelUserCache(param.UserIdList...)
	return nil
}
