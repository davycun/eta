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
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	HeaderAuthorization = "Authorization"
)

type HttpCase struct {
	Method       string                           // 请求类型
	Path         string                           // 链接
	Headers      map[string]string                // 请求头
	Body         any                              // 参数
	ResponseDest any                              //返回结果的序列化对象，需要放对象指针
	Desc         string                           // 描述
	ShowBody     bool                             // 是否展示返回
	Code         string                           // 希望的Response状态码
	HttpCode     int                              //希望的Http响应码
	Message      string                           // 错误信息
	TransferKey  string                           // 传输加密key
	ValidateFunc []func(t *testing.T, resp *Resp) // 校验方法
}

type Resp struct {
	Code    string      `json:"code"`
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

func NewBufferString(body string) io.Reader {
	return bytes.NewBufferString(body)
}

func PerformRequest(method, url string, headers map[string]string, body any) (c *gin.Context, r *http.Request, w *httptest.ResponseRecorder) {
	router := global.GetGin()
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	switch body.(type) {
	case string:
		bd := body.(string)
		r = httptest.NewRequest(method, url, NewBufferString(bd))
	default:
		bd, err := jsoniter.Marshal(body)
		if err != nil {
			logger.Errorf("")
			return
		}
		r = httptest.NewRequest(method, url, bytes.NewBuffer(bd))
	}

	c.Request = r
	for k, v := range headers {
		c.Request.Header.Set(k, v)
	}
	router.ServeHTTP(w, r)
	return
}

func Call(t *testing.T, testcase ...HttpCase) {
	if LoginToken == "" {
		initLogin()
	}
	for k, v := range testcase {
		if v.Headers == nil {
			v.Headers = make(map[string]string, 4)
		}
		if v.Method == "" {
			v.Method = http.MethodPost
		}
		if v.Headers[HeaderAuthorization] == "" {
			v.Headers[HeaderAuthorization] = LoginToken
		}
		_, _, w := PerformRequest(v.Method, v.Path, v.Headers, v.Body)

		fmt.Printf("第%d个测试用例：%s\n", k+1, v.Desc)
		if v.ShowBody {
			fmt.Printf(" 接口返回：%s\n", w.Body.String())
		}

		var (
			err       error
			resp      = Resp{}
			respBytes = w.Body.Bytes()
		)

		if algo, ok := v.Headers[constants.HeaderCryptSymmetryAlgorithm]; ok {
			var objmap map[string]string
			err = json.Unmarshal(respBytes, &objmap)
			assert.Nil(t, err)
			tk := v.TransferKey
			content, err1 := crypt.DecryptBase64(algo, tk, objmap["content"])
			assert.Nil(t, err1)
			if v.ShowBody {
				fmt.Printf("接口返回解密后明文：%s\n", content)
			}
			err = json.Unmarshal(content, &resp)
			if v.ResponseDest != nil {
				jsoniter.Unmarshal(content, v.ResponseDest)
			}
		} else {
			err = json.Unmarshal(respBytes, &resp)
			if v.ResponseDest != nil {
				jsoniter.Unmarshal(respBytes, v.ResponseDest)
			}
		}

		if v.Code == "" {
			v.Code = "200"
		}
		if v.HttpCode == 0 {
			v.HttpCode = 200
		}

		assert.NoError(t, err)
		assert.Equal(t, v.HttpCode, w.Code, "http状态码不一致")
		assert.Equal(t, v.Code, resp.Code, "错误码不一致")
		if len(v.ValidateFunc) > 0 {
			for _, v := range v.ValidateFunc {
				v(t, &resp)
			}
		}
	}
}
