package utils

import (
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/gin-gonic/gin"
	"net"
	"net/url"
	"strings"
)

const (
	HTTP  = "http"
	HTTPS = "https"
)

func GetLocalHost() string {
	interfaces, err := net.Interfaces()
	if err == nil {
		for _, v := range interfaces {
			if v.Flags&net.FlagUp != 0 {
				ads, _ := v.Addrs()
				for _, addr := range ads {
					if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
						to4 := ipNet.IP.To4()
						if to4 != nil {
							return to4.String()
						}
					}
				}
			}
		}
	}
	return "localhost"
}

func Scheme(c *gin.Context) string {
	// Can't use `r.Request.URL.Scheme`
	// See: https://groups.google.com/forum/#!topic/golang-nuts/pMUkBlQBDF0
	if c.Request.TLS != nil {
		return HTTPS
	}
	if scheme := c.Request.Header.Get("X-Forwarded-Proto"); scheme != "" {
		return scheme
	}
	if scheme := c.Request.Header.Get("X-Forwarded-Protocol"); scheme != "" {
		return scheme
	}
	if ssl := c.Request.Header.Get("X-Forwarded-Ssl"); ssl == "on" {
		return HTTPS
	}
	if scheme := c.Request.Header.Get("X-Url-Scheme"); scheme != "" {
		return scheme
	}
	return HTTP
}

func HttpScheme(secure bool) string {
	if secure {
		return HTTPS
	}
	return HTTP
}

// RequestHost 返回请求的host，会去掉默认端口
func RequestHost(c *gin.Context) string {
	host := c.Request.Host
	s := Scheme(c)
	// 如果s==http,且host以:80结尾，则去掉:80
	if s == HTTP && strings.HasSuffix(host, ":80") {
		return host[:len(host)-3]
	}
	// 如果s==https,且host以:443结尾，则去掉:443
	if s == HTTPS && strings.HasSuffix(host, ":443") {
		return host[:len(host)-4]
	}
	return host
}

func GetIP(input string) string {
	ip := net.ParseIP(input)
	if ip != nil {
		return ip.String()
	}
	addrs, err := net.LookupHost(input)
	if err != nil {
		logger.Errorf("无法解析IP地址: %v", err)
		return ""
	}
	ipAddr := addrs[0]
	logger.Debugf("域名[%s]解析为IP[%s]", input, ipAddr)
	return ipAddr
}

func ParseProxy(proxyUrlStr string) (httpsProxy, httpProxy, socks5Proxy string, err error) {
	if proxyUrlStr == "" {
		return
	}
	proxyUrl, err1 := url.Parse(proxyUrlStr)
	if err1 != nil {
		err = err1
		return
	}
	if proxyUrl == nil {
		return
	}
	if proxyUrl.Scheme == "http" || proxyUrl.Scheme == "https" {
		httpsProxy = proxyUrlStr
		httpProxy = proxyUrlStr
	} else if proxyUrl.Scheme == "socks5" {
		socks5Proxy = proxyUrlStr
	}
	return
}
