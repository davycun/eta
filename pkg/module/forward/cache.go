package forward

import (
	"encoding/json"
	"github.com/davycun/eta/pkg/common/crypt"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"net/http"
	"os"
	"path"
	"slices"
	"strings"
	"time"
)

const (
	cacheMd5Key = "123456789abcdefg"
)

type CacheData struct {
	ExpiredAt time.Time   `json:"expired_at,omitempty"`
	Header    http.Header `json:"header,omitempty"`
	Body      []byte      `json:"body,omitempty"`
	Status    int         `json:"status,omitempty"`
}

func (c *CacheData) IsValid() bool {
	return c.Status != 0 && (c.ExpiredAt.IsZero() || c.ExpiredAt.After(time.Now()))
}

type DefaultCache struct {
}

func MakeCacheKey(c *gin.Context, reqBody []byte, vendorName string) (string, error) {

	if vd, ok := GetCacheKeyMaker(vendorName); ok {
		return vd(c, reqBody, vendorName)
	}

	var (
		headerKeys = make([]string, 0, 1)
		queryKeys  = make([]string, 0, 1)
		sb         = strings.Builder{}
	)
	for k := range c.Request.Header {
		headerKeys = append(headerKeys, k)
	}
	for k := range c.Request.URL.Query() {
		queryKeys = append(queryKeys, k)
	}

	sb.WriteString(c.Request.Method)
	sb.WriteString(c.Request.URL.Path)

	//key 排序
	slices.Sort(headerKeys)
	slices.Sort(queryKeys)

	for _, k := range headerKeys {
		sb.WriteString(k + c.GetHeader(k))
	}
	for _, k := range queryKeys {
		sb.WriteString(k + c.Query(k))
	}
	if len(reqBody) > 0 {
		sb.Write(reqBody)
	}
	return crypt.EncryptHex(crypt.AlgoSignHmacMd5, cacheMd5Key, sb.String())
}

func MakeCacheData(vendor Vendor, resp *resty.Response) (CacheData, error) {
	//如果body为空可能是GET或者其他原因，如果之前发生了错误，就用delta的格式返回，如果没有就直接返回
	var (
		err       error
		respBody  = resp.Body()
		cacheData = CacheData{Body: make([]byte, 0), Header: make(http.Header)}
	)
	cacheData.Status = resp.StatusCode()
	for k := range resp.Header() {
		//取消掉跨域相关的头信息和压缩相关的头信息
		//resty和net/http都会处理gzip压缩，如果响应头中为Content-Encoding：gzip，说明服务端已经做了gzip压缩
		//那么 net/http会做解压，可以定制http.Transport.DisableCompression为true来拒绝net/http解压
		//但是如果net/http不解压了，resty还是会解压，但是resty不支持配置不解压
		//所以从resty拿到的body是解压后的，直接往客户端写，但如果给客户端的响应头没有去除Content-Encoding: gzip，那么客户端会拿非压缩的数据进行gzip解压会报错
		//TODO 思考：要不要保留Content-Encoding: gzip，然后对body压缩后再写入
		if strings.HasPrefix(k, "Cross") || strings.HasPrefix(k, "Content-Encoding") {
			continue
		}
		if utils.ContainAnyInsensitive(vendor.ExceptHeader, k) {
			continue
		}

		val := resp.Header().Get(k)
		if val != "" {
			cacheData.Header.Set(k, val)
		}
	}
	if len(respBody) > 0 {
		if fc, ok := GetHandleResponse(vendor.Name); ok {
			realBody, err1 := fc(resp, vendor, respBody)

			if err1 != nil {
				return cacheData, err1
			}
			if len(realBody) > 0 {
				respBody = realBody
			}
		}
		if len(respBody) > 0 {
			cacheData.Body = append(cacheData.Body, respBody...)
		}
	}
	return cacheData, err
}

func LoadCacheData(cacheKey string, vendor Vendor) (CacheData, error) {
	var (
		cacheData CacheData
	)
	if cacheKey == "" {
		return cacheData, nil
	}

	dir := vendor.CacheDir
	if dir == "" {
		dir = os.TempDir()
	}
	fileName := path.Join(dir, cacheKey+".json")
	if !fileutil.IsExist(fileName) {
		return cacheData, nil
	}
	dt, err := os.ReadFile(fileName)
	if err != nil {
		return cacheData, err
	}
	err = json.Unmarshal(dt, &cacheData)
	if err != nil {
		return cacheData, err
	}
	//如果已经过期，那就删除缓存文件
	if !cacheData.IsValid() {
		err = os.Remove(fileName)
		return CacheData{}, err
	}

	return cacheData, err
}

func SaveCacheData(cacheKey string, vendor Vendor, dt CacheData) error {
	if cacheKey == "" || !dt.IsValid() {
		return nil
	}

	//不限制http状态码的时候，默认只缓存200和204
	if len(vendor.CacheStatus) < 1 && !utils.ContainAny([]int{http.StatusOK, 204}, dt.Status) {
		return nil
	}

	//查看是否满足http状态码缓存规则
	if len(vendor.CacheStatus) > 0 && !utils.ContainAny(vendor.CacheStatus, dt.Status) {
		return nil
	}

	//设置过期时间
	if vendor.CacheExpire > 0 {
		dt.ExpiredAt = time.Now().Add(time.Second * time.Duration(vendor.CacheExpire))
	}

	//save data
	dir := vendor.CacheDir
	if dir == "" {
		dir = os.TempDir()
	}
	fileName := path.Join(dir, cacheKey+".json")
	jsDt, err1 := json.Marshal(dt)
	if err1 != nil {
		return err1
	}
	return os.WriteFile(fileName, jsDt, os.FileMode(0666))
}
