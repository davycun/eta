package weixin

import (
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/patrickmn/go-cache"
	"time"
)

var (
	WK         = WxKeys{}
	ApiCache   = cache.New(5*time.Minute, 10*time.Minute)
	HuaMu2025  = "huamu_2025"
	HuaMuYaoSu = "huamu_yaosu"
	MaiCeYaoSu = "maice_yaosu"
)

func init() {
	WK[HuaMu2025] = &WxKey{
		Key:       HuaMu2025,
		Type:      "DefaultService",
		AppId:     "wx906da58f13354e6c",
		AppSecret: "9a88cb67a759681b949ae55e5c0855d5",
		OriginId:  "gh_f94a72c11367",
		Name:      "花木街道数字孪生城市2025",
	}
	WK[HuaMuYaoSu] = &WxKey{
		Key:       HuaMuYaoSu,
		Type:      "Mini",
		AppId:     "wx6382728eda473357",
		AppSecret: "ec021e1bcb52925369670281d1cde98c",
		OriginId:  "gh_1a7ded290bdc",
		Name:      "花木街道市民数字要素",
	}
	WK[MaiCeYaoSu] = &WxKey{
		Key:       MaiCeYaoSu,
		Type:      "Mini",
		AppId:     "wx343feaa298672c6c",
		AppSecret: "2ed955bc0202a8e06828a9e6cbdb423b",
		OriginId:  "gh_90fe0c004fc2",
		Name:      "脉策城市全要素数字底座",
	}

	ApiCache.Set(HuaMu2025, NewWeiXin(WK[HuaMu2025]), cache.NoExpiration)
	ApiCache.Set(HuaMuYaoSu, NewWeiXin(WK[HuaMuYaoSu]), cache.NoExpiration)
	ApiCache.Set(MaiCeYaoSu, NewWeiXin(WK[MaiCeYaoSu]), cache.NoExpiration)

}

type WxKey struct {
	//微信类型：小程序、服务号、订阅号
	Key       string `json:"key"`
	Type      string `json:"type"`
	AppId     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
	OriginId  string `json:"origin_id"`
	Name      string `json:"name"`
}

type WxKeys map[string]*WxKey

func (w WxKeys) Get(key string) *WxKey {
	m, ok := w[key]
	if !ok {
		logger.Errorf("not found the key %s from WxKeys", key)
		return nil
	} else {
		return m
	}
}

func (w WxKeys) Set(key string, wxKey *WxKey) {
	w[key] = wxKey
}
