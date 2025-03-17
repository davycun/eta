package oauth2

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/cache"
	"github.com/davycun/eta/pkg/common/crypt"
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/user"
	"github.com/davycun/eta/pkg/module/user/userkey"
	"github.com/duke-git/lancet/v2/convertor"
	"time"
)

var (
	RequestValidIn = int64(60) // 秒
)

func LoginByAccessToken(c *ctx.Context, args any) (us user.User, err error) {

	var (
		param      = args.(*AccessKeyParam)
		userFilter = user.User{}
	)

	uk, err := userkey.LoadByKey(param.AccessKey)
	if err != nil {
		return us, err
	}
	if uk.UserId == "" {
		return us, errs.NewClientError("access key not exists")
	}

	userFilter.ID = uk.UserId
	userFilter.Valid = ctype.Boolean{Valid: true, Data: true}

	us, err = loadUserByFilter(userFilter)
	if err != nil {
		return
	}
	err = validateClientCredential(&us, uk, param)
	return us, nil
}

func validateClientCredential(u *user.User, uk userkey.UserKey, param *AccessKeyParam) (err error) {
	nowTs := time.Now().UTC().Unix()
	dur := nowTs - param.Ts
	if dur < -RequestValidIn || dur > RequestValidIn { // 前后兼容5秒
		logger.Errorf("u: %v, param:%v, 请求过期", u, param)
		return errs.NewClientError("请求过期")
	}
	//通过AccessKey和AccessSecure方式只能是OpenApi类型的用户
	if u.Category != constants.UserTypeOpenApi {
		return errs.NewClientError("用户类型错误")
	}
	exists, _ := cache.Exists(constants.RedisKey(constants.AccessTokenNonceKey, param.AccessKey, param.Nonce))
	if exists {
		logger.Errorf("u: %v, param:%v, 请求重复", u, param)
		return errs.NewClientError("请求重复")
	}
	calcSign, err := crypt.NewEncrypt(param.Algo, ctype.ToString(uk.AccessSecure)).
		FromRawString(fmt.Sprintf("%s%s", convertor.ToString(param.Ts), param.Nonce)).
		ToHexString()

	if err != nil || calcSign == "" {
		logger.Errorf("u: %v, param:%v, 签名计算失败: %v", u, param, err)
		return errs.NewClientError("签名计算失败")
	}
	if param.Sign == "" || calcSign != param.Sign {
		logger.Errorf("u: %v, param:%v, 签名验证失败", u, param)
		return errs.NewClientError("签名验证失败")
	}

	err = cache.SetEx(constants.RedisKey(constants.AccessTokenNonceKey, param.AccessKey, param.Nonce), 1, time.Second*time.Duration(RequestValidIn))
	return err
}
