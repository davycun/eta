package forward

import (
	"errors"
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/setting"
	"gorm.io/gorm"
	"slices"
	"strings"
)

const (
	settingForwardCategory = "setting_forward_category"
	settingForwardName     = "setting_forward_name"
	PathParam              = "path"
	PathVendor             = "vendor"
)

// Vendor
// CacheUri示例: 以xxx@url格式，其中xxx表示请求方法如果有多个用逗号隔开（匹配所有可以用*来代替），url表示的是url请求的path（支持正则表达式）

type Vendor struct {
	setting.BaseCredentials
	Name          string   `json:"name"`           //供应商名字
	Cache         bool     `json:"cache"`          //是否开启缓存
	CacheDir      string   `json:"cache_dir"`      //缓存的目录
	CacheExpire   int64    `json:"cache_expire"`   //过期，单位秒，如果是0表示永不过期
	CacheUri      []string `json:"cache_uri"`      //哪些请求需要缓存，注意这里的格式是METHOD@URI，比如GET@/api/a/b
	CacheStatus   []int    `json:"cache_status"`   //配置哪些响应状态需要缓存，不配置的话，默认只有200的进行缓存
	Sorted        bool     `json:"sorted"`         //表示cacheUri是否已经排序过，避免不断地排序
	Md5Key        string   `json:"md5_key"`        //通过md5生成缓存的key的密钥
	ExcludeHeader []string `json:"exclude_header"` //哪些头不参与缓存及返回
}

func (v *Vendor) NeedCache(method, uri string) bool {
	if !v.Cache {
		return false
	}
	//如果没有配置CacheUri，但是开启了Cache为True，表示所有的路径都缓存
	if len(v.CacheUri) < 1 {
		return true
	}
	v.sortCacheUri()

	return utils.IsMatchedUri(method, uri, v.CacheUri...)
}

func (v *Vendor) sortCacheUri() {
	if v.Sorted {
		return
	}
	//第一斜杠多排前面先匹配，第二长度长排前面先匹配
	slices.SortStableFunc(v.CacheUri, func(a, b string) int {
		la := len(a)
		lb := len(b)
		lsa := strings.Count(a, "/")
		lsb := strings.Count(b, "/")
		if lsa < lsb {
			return 1
		} else if lsa > lsb {
			return -1
		} else {
			if la < lb {
				return 1
			} else if la > lb {
				return -1
			} else {
				return 0
			}
		}
	})
	v.Sorted = true
}

type Config struct {
	Vendors map[string]Vendor `json:"vendors,omitempty"`
}

func (c Config) GetVendor(vendor string) (Vendor, error) {
	if c.Vendors == nil {
		return Vendor{}, errors.New(fmt.Sprintf("can not found the vendor %s setting", vendor))
	}
	vd, ok := c.Vendors[vendor]
	if !ok {
		return Vendor{}, errors.New(fmt.Sprintf("can not found the vendor %s setting", vendor))
	}
	return vd, nil
}

func GetVendor(db *gorm.DB, vendor string) (Vendor, error) {
	cfg, err := setting.GetConfig[Config](db, settingForwardCategory, settingForwardName)
	if err != nil {
		return Vendor{}, err
	}
	if cfg.Vendors == nil {
		cfg.Vendors = make(map[string]Vendor)
	}
	dfCfg := setting.GetDefault[Config](settingForwardCategory, settingForwardName)
	if dfCfg.Vendors != nil {
		for k, v := range dfCfg.Vendors {
			if _, ok := cfg.Vendors[k]; !ok {
				cfg.Vendors[k] = v
			}
		}
	}
	return cfg.GetVendor(vendor)
}

// AddDefaultVendor
// 提供外部初始化扩展，主要是在程序初始化时调用，把一些默认的配置写入到数据库
func AddDefaultVendor(vendor string, cred Vendor) {
	var (
		cfg = setting.GetDefault[Config](settingForwardCategory, settingForwardName)
	)
	if cfg.Vendors == nil {
		cfg.Vendors = make(map[string]Vendor)
	}

	if _, ok := cfg.Vendors[vendor]; ok {
		logger.Warnf("the target %s for the forwarding has already been set will be overwritten", vendor)
	}
	cfg.Vendors[vendor] = cred

	set := setting.Setting{}
	set.Namespace = constants.NamespaceEta
	set.Category = settingForwardCategory
	set.Name = settingForwardName
	set.Content = ctype.NewJson(&cfg)
	setting.Registry(set)
}
