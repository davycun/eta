package third

import (
	"github.com/davycun/eta/pkg/module/app"
	"github.com/davycun/eta/pkg/module/third/ding"
	"github.com/davycun/eta/pkg/module/third/wecom"
	"github.com/davycun/eta/pkg/module/third/weixin"
	"github.com/davycun/eta/pkg/module/third/zzd"
)

var (
	AppVersion = make(map[string]int64)
)

func CheckAppVersion(ap app.App) {
	id := ap.ID
	currentVersion := ap.UpdatedAt
	v, ok := AppVersion[ap.ID]
	if !ok {
		v = 0
	}
	if currentVersion != v {
		ClearApiCache(ap)
		AppVersion[id] = currentVersion
	}
}

/*
ClearApiCache 删除缓存

只要 app 配置有变更，就全部删除。以后可以优化成只删除当前app的缓存
*/
func ClearApiCache(ap app.App) {
	ding.ApiCache.Flush()
	wecom.ApiCache.Flush()
	weixin.ApiCache.Flush()
	zzd.ApiCache.Flush()
}
