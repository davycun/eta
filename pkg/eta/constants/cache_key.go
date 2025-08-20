package constants

import (
	"fmt"
	"strings"
)

const (
	// user
	AccessTokenNonceKey    = "eta:user:access_token_nonce:%s:%s" //第一个是accessKey，第二个是随机数nonce
	UserKey                = "eta:user:id:info:%s"
	User2DeptCacheKey      = `eta:user:id:user2dept:%s`
	UserTokenKey           = "eta:user:id:token:%s"
	UserOpenApiIdKey       = "eta:user:openapi:user_id:%s"
	TokenKey               = "eta:user:token:user_id_dept_id:%s"
	DatlasUserIdKey        = "eta:user:datlas_user_id:user_id_dept_id:%s"
	DatlasDeltaTokenKey    = "eta:user:datlas_eta_token:%s"
	TransferCryptoTokenKey = "eta:user:token:crypto:%s:%s"
	// login fail lock
	LoginFailLock    = "eta:login:fail:lock:%s"
	LoginFailCounter = "eta:login:fail:counter:%s:%s"
	// app
	AppKey = "eta:app:id:info:%s"
	// dept
	DeptCacheKey = `eta:dept:id:info:%s`
	// captcha
	CaptchaCodeKey  = "eta:captcha:code:%s"
	CaptchaPhoneKey = "eta:captcha:phone:%s:%s"
	// auth
	Auth2RoleKey       = `eta:auth:role_id:auth2role:%s:%s`
	PermissionCacheKey = `eta:auth:permission_id:permission:%s`
	UserRoleIdsKey     = "eta:auth:user_id:role_ids:%s"

	// tourist_forecast
	TfCouponClaimUserLock = "eta:tourist_forecast:coupon:cliam_user_lock:%s"
	//AllDataCacheKey
	CacheAllDataMenu           = "eta:allData:eta:menu:%s"
	CacheAllDataSetting        = "eta:allData:eta:setting:%s"        //存放在appDB下
	CacheAllDataConfig         = "eta:allData:eta:config:%s"         ///从localDB下取congif
	CacheAllDataConfigSetting  = "eta:allData:eta:config_setting:%s" //从app下的setting取Config
	CacheAllDataUser2App       = "eta:allData:eta:user2app:%s"
	CacheAllDataDictionary     = "eta:allData:eta:dictionary:%s"
	CacheAllDataTemplate       = "eta:allData:eta:template:%s"
	CacheAllDataSubscriber     = "eta:allData:eta:subscriber:%s"
	CacheAllDataLabel          = "eta:allData:citizen:label:%s"
	CacheAllDataAddress        = "eta:allData:citizen:address:%s"
	CacheAllDataBuilding       = "eta:allData:citizen:building:%s"
	CacheAllDataBd2Addr        = "eta:allData:citizen:bd2addr:%s"
	CacheAllDataFloor          = "eta:allData:citizen:floor:%s"
	CacheAllDataRoom           = "eta:allData:citizen:room:%s"
	CacheAllDataIdNameUser     = "eta:allData:eta:idNameUser:%s"
	CacheAllDataIdNameDept     = "eta:allData:eta:idNameDept:%s"
	CacheAllDataSearchBuilding = "eta:allData:eta:searchBuilding:%s"
	CacheAllDataSearchRoom     = "eta:allData:eta:searchRoom:%s"

	LockPep2RoomLive = "eta:citizen:lock:pep2room_live"

	// api
	APIDatlasTokenKey = "eta:api:datlas_token"
)

func RedisKey(key string, params ...any) string {
	return fmt.Sprintf(key, params...)
}
func GetAllDataRedisKeyId(key string) string {
	split := strings.Split(key, ":")
	return split[len(split)-1]
}
