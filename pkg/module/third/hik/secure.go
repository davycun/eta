package hik

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/davycun/eta/pkg/common/id/nanoid"
	"golang.org/x/net/proxy"
	"io"
	"net"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

// HKClient 海康OpenAPI配置参数
type HKClient struct {
	Ip      string //平台ip
	Port    int    //平台端口
	AppKey  string //平台APPKey
	Secret  string //平台APPSecret
	IsHttps bool   //是否使用HTTPS协议
	dial    *net.Dialer
	client  *http.Client
}

type Result struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type Data struct {
	Total    int                      `json:"total"`
	PageSize int                      `json:"pageSize"`
	PageNo   int                      `json:"pageNo"`
	List     []map[string]interface{} `json:"list"`
}

func NewHikClient(ip string, port int, key, sec string) *HKClient {
	hk := &HKClient{
		Ip:      ip,
		Port:    port,
		AppKey:  key,
		Secret:  sec,
		IsHttps: true,
	}
	hk.initClient()
	return hk
}

func (hk *HKClient) initClient() {

	if hk.client != nil {
		return
	}
	hk.client = &http.Client{}

	hk.dial = &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	tran := &http.Transport{
		ForceAttemptHTTP2:     true,
		DialContext:           hk.dial.DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	if hk.IsHttps {
		tran.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	hk.client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	hk.client.Timeout = time.Duration(30) * time.Second

	hk.client.Transport = tran
}

func (hk *HKClient) SocksProxy(address string, at *proxy.Auth) *HKClient {
	hk.initClient()
	tran := hk.client.Transport.(*http.Transport)
	proxyDialer, _ := proxy.SOCKS5("tcp", address, at, hk.dial)

	if cd, ok := proxyDialer.(proxy.ContextDialer); ok {
		tran.DialContext = cd.DialContext
	} else {
		tran.Dial = proxyDialer.Dial
	}
	return hk
}

func (hk *HKClient) HttpPost(url string, body map[string]any, rs any, header map[string]string) error {
	hk.initClient()
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return err
	}
	err = hk.initRequest(header, url, string(bodyJson), true)
	if err != nil {
		return err
	}
	var sb []string
	if hk.IsHttps {
		sb = append(sb, "https://")
	} else {
		sb = append(sb, "http://")
	}
	sb = append(sb, fmt.Sprintf("%s:%d", hk.Ip, hk.Port))
	sb = append(sb, url)

	req, err := http.NewRequest("POST", strings.Join(sb, ""), bytes.NewReader(bodyJson))
	if err != nil {
		return err
	}

	req.Header.Set("Accept", header["Accept"])
	req.Header.Set("Content-Type", header["Content-Type"])
	for k, v := range header {
		if strings.Contains(k, "x-ca-") {
			req.Header.Set(k, v)
		}
	}

	resp, err := hk.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		var resBody []byte
		resBody, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		err = json.Unmarshal(resBody, rs)
	} else if resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusMovedPermanently {
		reqUrl := resp.Header.Get("Location")
		err = fmt.Errorf("HttpPost Response StatusCode：%d，Location：%s", resp.StatusCode, reqUrl)
	} else {
		err = fmt.Errorf("HttpPost Response StatusCode：%d", resp.StatusCode)
	}
	return err
}

// initRequest 初始化请求头
func (hk *HKClient) initRequest(header map[string]string, url, body string, isPost bool) error {
	header["Accept"] = "application/json"
	header["Content-Type"] = "application/json"
	if isPost {
		var err error
		header["content-md5"], err = computeContentMd5(body)
		if err != nil {
			return err
		}
	}
	header["x-ca-timestamp"] = strconv.FormatInt(time.Now().UnixMilli(), 10)
	header["x-ca-nonce"] = nanoid.New()
	header["x-ca-key"] = hk.AppKey

	var strToSign string
	if isPost {
		strToSign = buildSignString(header, url, "POST")
	} else {
		strToSign = buildSignString(header, url, "GET")
	}
	signedStr, err := computeForHMACSHA256(strToSign, hk.Secret)
	if err != nil {
		return err
	}
	header["x-ca-signature"] = signedStr
	return nil
}

// computeContentMd5 计算content-md5
func computeContentMd5(body string) (string, error) {
	h := md5.New()
	_, err := h.Write([]byte(body))
	if err != nil {
		return "", err
	}
	md5Str := hex.EncodeToString(h.Sum(nil))
	return base64.StdEncoding.EncodeToString([]byte(md5Str)), nil
}

// computeForHMACSHA256 计算HMACSHA265
func computeForHMACSHA256(str, secret string) (string, error) {
	mac := hmac.New(sha256.New, []byte(secret))
	_, err := mac.Write([]byte(str))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(mac.Sum(nil)), nil
}

// buildSignString 计算签名字符串
func buildSignString(header map[string]string, url, method string) string {
	var sb []string
	sb = append(sb, strings.ToUpper(method))
	sb = append(sb, "\n")

	if header != nil {
		if _, ok := header["Accept"]; ok {
			sb = append(sb, header["Accept"])
			sb = append(sb, "\n")
		}
		if _, ok := header["Content-MD5"]; ok {
			sb = append(sb, header["Content-MD5"])
			sb = append(sb, "\n")
		}
		if _, ok := header["Content-Type"]; ok {
			sb = append(sb, header["Content-Type"])
			sb = append(sb, "\n")
		}
		if _, ok := header["Date"]; ok {
			sb = append(sb, header["Date"])
			sb = append(sb, "\n")
		}
	}
	sb = append(sb, buildSignHeader(header))
	sb = append(sb, url)
	return strings.Join(sb, "")
}

// buildSignHeader 计算签名头
func buildSignHeader(header map[string]string) string {
	var sortedDicHeader map[string]string
	sortedDicHeader = header

	var sslice []string
	for key, _ := range sortedDicHeader {
		sslice = append(sslice, key)
	}
	sort.Strings(sslice)

	var sbSignHeader []string
	var sb []string
	//在将key输出
	for _, k := range sslice {
		if strings.Contains(strings.ReplaceAll(k, " ", ""), "x-ca-") {
			sb = append(sb, k+":")
			if sortedDicHeader[k] != "" {
				sb = append(sb, sortedDicHeader[k])
			}
			sb = append(sb, "\n")
			if len(sbSignHeader) > 0 {
				sbSignHeader = append(sbSignHeader, ",")
			}
			sbSignHeader = append(sbSignHeader, k)
		}
	}

	header["x-ca-signature-headers"] = strings.Join(sbSignHeader, "")
	return strings.Join(sb, "")
}
