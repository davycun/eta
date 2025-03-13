package datlas

import (
	"errors"
	"fmt"
	delta_cache "github.com/davycun/eta/pkg/common/cache"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/patrickmn/go-cache"
	"time"
)

var (
	ApiCache            = cache.New(5*time.Minute, 10*time.Minute)
	RequestDatlasHeader = "X-Need-Request-Datlas" // 请求 datlas 接口时，添加这个 header。当因为配置错误而导致请求到delta时，可以通过这个header来判断
)

type Datlas struct {
	BaseUrl  string `json:"base_url"`
	Username string `json:"username"`
	Password string `json:"password"`
	Err      error
	client   *resty.Client
}

func NewDatlas(baseUrl, username, password string, c *gin.Context) *Datlas {
	if baseUrl == "" {
		baseUrl = DefaultBaseUrl(c)
	}
	apiCacheKey := fmt.Sprintf(`%s%s`, baseUrl, username)
	if a, found := ApiCache.Get(apiCacheKey); found {
		return a.(*Datlas)
	}

	w := &Datlas{
		BaseUrl:  baseUrl,
		Username: username,
		Password: password,
		client:   resty.New(), // 创建一个resty客户端
	}
	w.client.SetTimeout(1 * time.Minute)
	w.client.SetBaseURL(baseUrl)
	ApiCache.Set(apiCacheKey, w, cache.NoExpiration)
	return w
}

func (d *Datlas) Auth() *Auth {
	return &Auth{
		Datlas: *d,
	}
}

func (d *Datlas) GetToken() string {
	token := &Token{}

	err, _ := delta_cache.Get(constants.RedisKey(constants.APIDatlasTokenKey), token)
	if err != nil {
		logger.Errorf("获取datlas token 失败, %v", err)
		return ""
	}

	if token.IsExpired() {
		loginResp := d.Auth().LoginByName(&LoginByNameParam{Name: d.Username, Password: d.Password})
		if loginResp == nil || loginResp.Rc != 0 || loginResp.Result.Auth == "" {
			logger.Errorf("登录datlas失败")
			d.Err = errors.New("登录datlas失败")
			return ""
		}
		token.Auth = loginResp.Result.Auth
		token.MdtUser = loginResp.Result.MdtUser
		token.ExpiresIn = TokenExpireIn
		token.ExpiredAt = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)

		err = delta_cache.SetEx(constants.RedisKey(constants.APIDatlasTokenKey), token, time.Second*time.Duration(TokenExpireIn))
		if err != nil {
			logger.Errorf("保存datlas token 失败, %v", err)
		}
	}
	return token.Auth
}

func DefaultBaseUrl(c *gin.Context) string {
	if c == nil {
		return ""
	}
	schema := utils.Scheme(c)
	datlasBaseUrl := fmt.Sprintf("%s://%s", schema, c.Request.Host)
	logger.Warnf("生成默认 datlas base url：%s", datlasBaseUrl)
	return datlasBaseUrl
}
