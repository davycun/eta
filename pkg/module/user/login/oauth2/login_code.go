package oauth2

import (
	"errors"
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/module/setting"
	"github.com/davycun/eta/pkg/module/third/ding"
	"github.com/davycun/eta/pkg/module/third/wecom"
	"github.com/davycun/eta/pkg/module/third/weixin"
	"github.com/davycun/eta/pkg/module/third/zzd"
	"github.com/davycun/eta/pkg/module/user"
	"github.com/davycun/eta/pkg/module/user/login/captcha"
	"strconv"
)

// LoginByDingCode
// 钉钉免密登录、钉钉扫码登录
func LoginByDingCode(c *ctx.Context, args any) (us user.User, err error) {
	var (
		cfg, exists = setting.GetLoginConfig(c.GetAppGorm())
		userFilter  = user.User{}
		param       = args.(*LoginByCodeParam)
	)
	if !exists || !cfg.ThirdDing.Valid() {
		return us, errs.NewServerError("钉钉配置信息有误")
	}

	dingApi := ding.NewDing(cfg.ThirdDing.BaseUrl, cfg.ThirdDing.AppKey, cfg.ThirdDing.AppSecret)
	if dingApi == nil {
		return us, errors.New("钉钉配置信息有误")
	}
	userAccessToken := dingApi.Oauth2().UserAccessToken(param.Code)
	ui := dingApi.User().GetUser("me", userAccessToken.AccessToken)
	if dingApi.Err != nil {
		return us, dingApi.Err
	}
	if ui == nil || ui.UnionId == "" {
		return us, errs.NewServerError("通过登录授权码获取钉钉用户信息失败")
	}
	userFilter.DingUnionId = ctype.NewStringPrt(ui.UnionId)

	return loadUserByFilter(userFilter)
}

// LoginByZzDingCode
// 浙政钉免密登录
func LoginByZzDingCode(c *ctx.Context, args any) (us user.User, err error) {
	var (
		cfg, exists = setting.GetLoginConfig(c.GetAppGorm())
		userFilter  = user.User{}
		param       = args.(*LoginByCodeParam)
	)

	if !exists || !cfg.ThirdZzd.Valid() {
		return us, errs.NewServerError("浙政钉配置信息有误")
	}

	zzdApi := zzd.NewZzd(cfg.ThirdZzd.BaseUrl, cfg.ThirdZzd.AppKey, cfg.ThirdZzd.AppSecret, cfg.ThirdZzd.TenantId)
	if zzdApi == nil {
		return us, errs.NewServerError("浙政钉配置信息有误")
	}
	rr := zzdApi.Oauth2().GetUserInfo(param.Code)
	if zzdApi.Err != nil {
		return us, zzdApi.Err
	}
	if rr == nil || !rr.Success || !rr.Content.Success {
		return us, errs.NewServerError("通过免密登录授权码获取浙政钉用户信息失败")
	}
	userFilter.ZzdAccountId = ctype.NewStringPrt(strconv.Itoa(rr.Content.Data.AccountId))

	return loadUserByFilter(userFilter)
}

// LoginByZzDingQrCode
// 浙政钉扫码登录
func LoginByZzDingQrCode(c *ctx.Context, args any) (us user.User, err error) {
	var (
		cfg, exists = setting.GetLoginConfig(c.GetAppGorm())
		userFilter  = user.User{}
		param       = args.(*LoginByCodeParam)
	)
	if !exists || !cfg.ThirdZzd.Valid() {
		return us, errs.NewServerError("浙政钉配置信息有误")
	}

	zzdApi := zzd.NewZzd(cfg.ThirdZzd.BaseUrl, cfg.ThirdZzd.AppKey, cfg.ThirdZzd.AppSecret, cfg.ThirdZzd.TenantId)
	if zzdApi == nil {
		return us, errs.NewServerError("浙政钉配置信息有误")
	}
	rr := zzdApi.Oauth2().GetUserInfoByCode(param.Code)
	if zzdApi.Err != nil {
		return us, zzdApi.Err
	}
	if rr == nil || !rr.Success || !rr.Content.Success {
		return us, errs.NewServerError("通过扫码登录授权码获取浙政钉用户信息失败")
	}
	userFilter.ZzdAccountId = ctype.NewStringPrt(strconv.Itoa(rr.Content.Data.AccountId))

	return loadUserByFilter(userFilter)
}

// LoginByWechatCode
// 微信免密登录、微信扫码登录
func LoginByWechatCode(c *ctx.Context, args any) (us user.User, err error) {
	var (
		cfg, exists = setting.GetLoginConfig(c.GetAppGorm())
		userFilter  = user.User{}
		param       = args.(*LoginByCodeParam)
	)

	if !exists || !cfg.ThirdWeChat.Valid() {
		return us, errs.NewServerError("微信配置信息有误")
	}

	wx := weixin.NewWeiXin(&weixin.WxKey{
		Key:       cfg.ThirdWeChat.AppKey,
		Type:      "DefaultService",
		AppId:     cfg.ThirdWeChat.AppKey,
		AppSecret: cfg.ThirdWeChat.AppSecret,
		OriginId:  cfg.ThirdWeChat.OriginId,
		Name:      cfg.ThirdWeChat.Name,
	})

	if wx == nil {
		return us, errs.NewServerError("微信配置信息有误")
	}
	atr, err := wx.Oauth2().AccessToken(param.Code)
	if err != nil {
		return us, err
	}
	if atr == nil || atr.Unionid == "" {
		return us, errs.NewServerError("通过登录授权码获取微信用户信息失败")
	}
	userFilter.WechatUnionId = ctype.NewStringPrt(atr.Unionid)
	return loadUserByFilter(userFilter)
}

// LoginByWeComCode
// 企业微信免密登录、企业微信扫码登录
func LoginByWeComCode(c *ctx.Context, args any) (us user.User, err error) {
	var (
		cfg, exists = setting.GetLoginConfig(c.GetAppGorm())
		userFilter  = user.User{}
		param       = args.(*LoginByCodeParam)
	)

	if !exists || !cfg.ThirdWeCom.Valid() {
		return us, errs.NewServerError("企业微信配置信息有误")
	}

	wecomApi := wecom.NewWecom(cfg.ThirdWeCom.BaseUrl, cfg.ThirdWeCom.AppKey, cfg.ThirdWeCom.AppSecret)
	if wecomApi == nil {
		return us, errs.NewServerError("企业微信配置信息有误")
	}
	ui := wecomApi.Oauth2().GetUserInfo(param.Code)
	if wecomApi.Err != nil {
		return us, wecomApi.Err
	}
	if ui == nil || ui.Userid == "" {
		return us, errs.NewServerError("通过登录授权码获取企业微信用户信息失败")
	}
	userFilter.WecomUserId = ctype.NewStringPrt(ui.Userid)

	return loadUserByFilter(userFilter)
}

// LoginBySmsCode
// aliYun短信登录
func LoginBySmsCode(c *ctx.Context, args any) (us user.User, err error) {
	var (
		userFilter = user.User{}
		param      = args.(*LoginByCodeParam)
	)

	if !captcha.Verify(captcha.Captcha{Code: param.Code, Phone: param.Phone}) {
		return us, errs.NewClientError("验证码错误")
	}
	userFilter.Phone = ctype.NewStringPrt(param.Phone)

	return loadUserByFilter(userFilter)
}

func loadUserByFilter(userFilter user.User) (us user.User, err error) {

	var (
		userList []user.User
	)
	err = global.GetLocalGorm().Model(&userList).Where(&userFilter).Find(&userList).Error
	if err != nil {
		return
	}
	if len(userList) < 1 {
		return us, errs.NewClientError("没有绑定，请先绑定")
	}
	if len(userList) > 1 {
		return us, errs.NewClientError("账号被绑定了多个用户，请联系管理员")
	}
	return userList[0], nil
}
