package user

import (
	"errors"
	"fmt"
	"github.com/davycun/eta/pkg/common/cache"
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/app"
	"github.com/davycun/eta/pkg/module/setting"
	"github.com/davycun/eta/pkg/module/user/user2app"
	"github.com/davycun/eta/pkg/module/user/userkey"
	"github.com/duke-git/lancet/v2/maputil"
	"github.com/duke-git/lancet/v2/slice"
	"gorm.io/gorm"
	"strings"
	"time"
)

type TokenInfo struct {
	Token     string    `json:"token,omitempty"`
	UserId    string    `json:"user_id,omitempty"`
	DeptId    string    `json:"dept_id,omitempty"`
	AppId     string    `json:"app_id,omitempty"`
	Key       string    `json:"key,omitempty"`       //传输加密中的Key
	OnlyOne   bool      `json:"only_one,omitempty"`  //是否只是允许一个用户一个token，也就是同一个用户返回同一个token
	ExpiredAt time.Time `json:"expire_at,omitempty"` //在什么时间点过期
}

func CheckTokenInfo(t *TokenInfo) bool {
	if t.UserId == "" || t.AppId == "" || t.Token == "" {
		return false
	}
	if t.DeptId == "" {
		t.DeptId = t.UserId
	}
	if t.ExpiredAt.IsZero() {
		cfg, _ := setting.GetLoginConfig(global.GetLocalGorm())
		dur := cfg.TokenExpireIn
		t.ExpiredAt = time.Now().Add(time.Second * time.Duration(dur))
	}
	return true
}
func (tk TokenInfo) ExpireIn() time.Duration {
	if tk.ExpiredAt.IsZero() {
		return GetTokenExpireIn()
	}
	return tk.ExpiredAt.Sub(time.Now())
}

type TokenDept struct {
	DeptId     string    `json:"dept_id,omitempty"`
	CreateTime time.Time `json:"create_time,omitempty"`
}

func GetTokenExpireIn() time.Duration {
	cfg, _ := setting.GetLoginConfig(global.GetLocalGorm())
	if cfg.TokenExpireIn <= 0 {
		return time.Second * time.Duration(setting.DefaultTokenExpireIn)
	}
	return time.Second * time.Duration(cfg.TokenExpireIn)
}

// StoreToken
// 支持新增或者更新
func StoreToken(tk TokenInfo) error {
	if !CheckTokenInfo(&tk) {
		return errors.New("the appId and userId can not be empty")
	}
	var (
		token2UserKey = constants.RedisKey(constants.TokenKey, tk.Token)
		expiration    = tk.ExpireIn()
	)
	return caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			return storeUser2Token(tk, expiration)
		}).
		Call(func(cl *caller.Caller) error {
			return cache.SetEx(token2UserKey, tk, expiration)
		}).Err
}
func storeUser2Token(tk TokenInfo, expiration time.Duration) error {
	var (
		user2TokenKey = constants.RedisKey(constants.UserTokenKey, tk.UserId)
		tkMap         = make(map[string]TokenInfo) //token -> tokenInfo
	)
	if tk.Token == "" || tk.UserId == "" {
		logger.Error("StoreUserToken token or userId is empty.")
		return nil
	}
	if tk.OnlyOne {
		return cache.SetEx(user2TokenKey, map[string]TokenInfo{tk.Token: tk}, expiration)
	}
	_, err := cache.Get(user2TokenKey, &tkMap)
	if err != nil {
		return err
	}
	tkMap = maputil.Filter(tkMap, func(key string, value TokenInfo) bool {
		return time.Now().Sub(value.ExpiredAt) > 0
	})

	tkMap[tk.Token] = tk
	return cache.SetEx(user2TokenKey, tkMap, expiration)
}

func LoadTokenByUserId(userId string) (tkList []TokenInfo, err error) {
	var (
		user2TokenKey = constants.RedisKey(constants.UserTokenKey, userId)
	)
	_, err = cache.Get(user2TokenKey, &tkList)
	return
}
func LoadTokenByToken(token string) (tk TokenInfo, err error) {
	var (
		tokenKey = constants.RedisKey(constants.TokenKey, token)
	)
	_, err = cache.Get(tokenKey, &tk)
	//如果是fixedToken那么进行自动处理
	if tk.UserId == "" && strings.HasPrefix(token, constants.LoginTypeFixToken) {
		var (
			ap    = app.App{}
			db    = global.GetLocalGorm()
			uk    userkey.UserKey
			appDb *gorm.DB
		)

		err = caller.NewCaller().
			Call(func(cl *caller.Caller) error {
				uk, err = userkey.LoadByFixToken(token)
				return err
			}).
			Call(func(cl *caller.Caller) error {
				if uk.UserId == "" {
					return errs.NewClientError("can not found the fixToken")
				}
				return nil
			}).
			Call(func(cl *caller.Caller) error {
				ap, err = user2app.LoadDefaultAppByUserId(db, uk.UserId)
				return err
			}).
			Call(func(cl *caller.Caller) error {
				if ap.ID == "" {
					return errs.NewClientError(fmt.Sprintf("can not found the app by user id[%s]", uk.UserId))
				}
				return nil
			}).
			Call(func(cl *caller.Caller) error {
				appDb, err = global.LoadGorm(ap.Database)
				return err
			}).
			Call(func(cl *caller.Caller) error {
				if appDb == nil {
					return errs.NewServerError(fmt.Sprintf("can not create gorm.DB with the app[%s]", ap.ID))
				}
				tk.DeptId, err = loadDeptIdByUserId(appDb, uk.ID)
				return err
			}).Err

		if err != nil {
			return
		}
		tk.Token = token
		tk.UserId = uk.UserId
		tk.AppId = ap.ID
		tk.OnlyOne = true
		return tk, StoreToken(tk)
	}
	return
}
func LoadUserByToken(db *gorm.DB, token string) (us User, err error) {
	tk, err := LoadTokenByToken(token)
	if err != nil {
		return
	}
	if tk.UserId == "" {
		return us, NotLogin
	}
	return LoadUserById(db, tk.UserId)
}

func DelUserToken(token string) error {
	var (
		err    error
		tk     = TokenInfo{}
		cfg, _ = setting.GetLoginConfig(global.GetLocalGorm())
	)
	err = caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			tk, err = LoadTokenByToken(token)
			if tk.UserId == "" {
				cl.Stop()
			}
			return err
		}).
		Call(func(cl *caller.Caller) error {
			_, err = cache.Del(constants.RedisKey(constants.TokenKey, token))
			return err
		}).
		Call(func(cl *caller.Caller) error {
			user2TokenKey := constants.RedisKey(constants.UserTokenKey, tk.UserId)
			tkMap := make(map[string]TokenInfo) //token -> tokenInfo
			_, err = cache.Get(user2TokenKey, &tkMap)
			if err != nil {
				return err
			}
			delete(tkMap, token)
			dur, _ := cache.TTL(user2TokenKey)
			if dur < 0 {
				dur = time.Second * time.Duration(cfg.TokenExpireIn)
			}
			return cache.SetEx(user2TokenKey, tkMap, dur)
		}).Err
	return err
}
func RenewalToken(token string, expiration time.Duration) error {
	var (
		err error
		tk  = TokenInfo{}
	)
	return caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			tk, err = LoadTokenByToken(token)
			return err
		}).
		Call(func(cl *caller.Caller) error {
			tk.ExpiredAt = time.Now().Add(expiration)
			return StoreToken(tk)
		}).Err
}

// LogOutUser
// 退出用指定app登录的用户的会话信息
func LogOutUser(userId string, appId string) error {

	tkList, err := LoadTokenByUserId(userId)

	if err != nil {
		return err
	}
	tkMap := make(map[string]TokenInfo)

	tkList = slice.Filter(tkList, func(index int, item TokenInfo) bool {
		return item.AppId == appId
	})
	for _, tk := range tkList {
		if tk.AppId == appId {
			tkKey := constants.RedisKey(constants.TokenKey, tk.Token)
			_, err = cache.Del(tkKey)
			if err != nil {
				return err
			}
		} else {
			tkMap[tk.Token] = tk
		}
	}

	//1.删除与用户关联的相关的api的token，一个用户可能有多个token
	//2.设置的时候需要注意过期时间
	user2TokenKey := constants.RedisKey(constants.UserTokenKey, userId)
	if len(tkMap) < 1 {
		_, err = cache.Del(user2TokenKey)
		if err != nil {
			return err
		}
	} else {
		cfg, _ := setting.GetLoginConfig(nil)
		dur, _ := cache.TTL(user2TokenKey)
		if dur < 0 {
			dur = time.Second * time.Duration(cfg.TokenExpireIn)
		}
		return cache.SetEx(user2TokenKey, tkMap, dur)
	}
	return nil
}
