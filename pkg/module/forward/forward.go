package forward

import (
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/module/setting"
	"github.com/go-resty/resty/v2"
)

var (
	reqHandler  = map[string]HandleRequest{}  //vendor -> HandleRequest
	respHandler = map[string]HandleResponse{} //vendor -> HandleResponse
)

type HandleRequest func(req *resty.Request, credentials setting.BaseCredentials, body []byte) ([]byte, error)
type HandleResponse func(resp *resty.Response, credentials setting.BaseCredentials, body []byte) ([]byte, error)

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

func GetHandleRequest(vendor string) (HandleRequest, bool) {
	fc, ok := reqHandler[vendor]
	return fc, ok
}
func GetHandleResponse(vendor string) (HandleResponse, bool) {
	fc, ok := respHandler[vendor]
	return fc, ok
}
