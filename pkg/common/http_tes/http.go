package http_tes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/davycun/eta/pkg/common/crypt"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	HeaderAuthorization = "Authorization"
)

type (
	ValidateFunc func(t *testing.T, resp *Response)
)

type HttpCase struct {
	Method       string            // 请求类型
	Path         string            // 链接
	Headers      map[string]string // 请求头
	Body         any               // 参数
	ResponseDest any               //返回结果的序列化对象，需要放对象指针
	Desc         string            // 描述
	ShowBody     bool              // 是否展示返回
	Code         string            // 希望的Response状态码
	HttpCode     int               //希望的Http响应码
	Message      string            // 错误信息
	TransferKey  string            // 传输加密key
	ValidateFunc []ValidateFunc    // 校验方法
}

type Resp struct {
	Code    string      `json:"code"`
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}
type Response struct {
	RawBody []byte
	Resp    Resp
}

func PerformRequest(method, url string, headers map[string]string, body any) *httptest.ResponseRecorder {

	var (
		r      *http.Request
		router = global.GetGin()
		w      = httptest.NewRecorder()
		c, _   = gin.CreateTestContext(w)
	)
	switch bd := body.(type) {
	case string:
		r = httptest.NewRequest(method, url, bytes.NewBufferString(bd))
	case *string:
		r = httptest.NewRequest(method, url, bytes.NewBufferString(*bd))
	case []byte:
		r = httptest.NewRequest(method, url, bytes.NewBuffer(bd))
	default:
		bds, err := jsoniter.Marshal(body)
		if err != nil {
			logger.Errorf("")
			return w
		}
		r = httptest.NewRequest(method, url, bytes.NewBuffer(bds))
	}

	c.Request = r
	for k, v := range headers {
		c.Request.Header.Set(k, v)
	}
	router.ServeHTTP(w, r)
	return w
}

func Call(t *testing.T, testcase ...HttpCase) {
	for _, v := range testcase {
		if v.Headers == nil {
			v.Headers = make(map[string]string, 4)
		}
		if v.Method == "" {
			v.Method = http.MethodPost
		}
		if v.Headers[HeaderAuthorization] == "" {
			v.Headers[HeaderAuthorization] = LoginToken
		}
		w := PerformRequest(v.Method, v.Path, v.Headers, v.Body)

		var (
			err      error
			response = Response{
				Resp: Resp{},
			}
		)
		processResponse(t, w, &response, v)
		if v.HttpCode == 0 {
			v.HttpCode = 200
		}

		assert.Nil(t, err)
		assert.Equal(t, v.HttpCode, w.Code)
		if len(v.ValidateFunc) > 0 {
			for _, fc := range v.ValidateFunc {
				fc(t, &response)
			}
		}
	}
}

func processResponse(t *testing.T, w *httptest.ResponseRecorder, response *Response, httpCase HttpCase) {
	response.RawBody = w.Body.Bytes()
	if len(response.RawBody) < 1 {
		return
	}
	var (
		err    error
		ct     = w.Header().Get("Content-Type")
		isJson = strings.Contains(ct, "application/json")
	)

	if algo, ok := httpCase.Headers[constants.HeaderCryptSymmetryAlgorithm]; ok {
		objMap := make(map[string]string)
		err = json.Unmarshal(response.RawBody, &objMap)
		assert.Nil(t, err)

		response.RawBody, err = crypt.DecryptBase64(algo, httpCase.TransferKey, objMap["content"])
		assert.Nil(t, err)
	}

	if isJson {
		err = json.Unmarshal(response.RawBody, &response.Resp)
		assert.Nil(t, err)
		if httpCase.ResponseDest != nil {
			err = jsoniter.Unmarshal(response.RawBody, httpCase.ResponseDest)
			assert.Nil(t, err)
		}
	}
	if httpCase.ShowBody {
		fmt.Printf("response body is：%s\n", response.RawBody)
	}
}
