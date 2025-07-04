package forward

import (
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

var (
	reqHandler    = map[string]HandleRequest{}  //vendorName -> HandleRequest
	respHandler   = map[string]HandleResponse{} //vendorName -> HandleResponse
	cacheKeyMaker = map[string]CacheKeyMaker{}  //vendorName -> CacheKeyMaker
)

type HandleRequest func(req *resty.Request, credentials Vendor, body []byte) ([]byte, error)
type HandleResponse func(resp *resty.Response, credentials Vendor, body []byte) ([]byte, error)

type CacheKeyMaker func(c *gin.Context, reqBody []byte, vendorName string) (string, error) //

func RegistryHandleRequest(vendor string, fc HandleRequest) {
	if _, ok := reqHandler[vendor]; ok {
		logger.Warnf("the HandleRequest of %s has already been set will be overwritten", vendor)
	}
	reqHandler[vendor] = fc
}
func RegistryHandleResponse(vendor string, fc HandleResponse) {
	if _, ok := respHandler[vendor]; ok {
		logger.Warnf("the HandleResponse of %s has already been set will be overwritten", vendor)
	}
	respHandler[vendor] = fc
}
func RegistryCacheKeyMaker(vendor string, fc CacheKeyMaker) {
	if _, ok := cacheKeyMaker[vendor]; ok {
		logger.Warnf("the CacheKeyMaker of %s has already been set will be overwritten", vendor)
	}
	cacheKeyMaker[vendor] = fc
}

func GetHandleRequest(vendor string) (HandleRequest, bool) {
	fc, ok := reqHandler[vendor]
	return fc, ok
}
func GetHandleResponse(vendor string) (HandleResponse, bool) {
	fc, ok := respHandler[vendor]
	return fc, ok
}
func GetCacheKeyMaker(vendor string) (CacheKeyMaker, bool) {
	fc, ok := cacheKeyMaker[vendor]
	return fc, ok
}
