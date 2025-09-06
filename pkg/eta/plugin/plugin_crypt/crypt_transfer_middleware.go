package plugin_crypt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/davycun/eta/pkg/common/crypt"
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/controller"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/security"
	"github.com/davycun/eta/pkg/module/user"
	"github.com/gin-gonic/gin"
	"io"
)

func TransferCrypt(c *gin.Context) {
	if !needDec(c) {
		return
	}

	responseBodyWriter := &ResponseCopyBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
	c.Writer = responseBodyWriter

	symmetryKey, err := getSymmetryKey(c)
	if err != nil {
		controller.ProcessResult(c, nil, err)
		return
	}

	next := transferDecrypt(c, symmetryKey)
	if !next {
		return
	}
	c.Next()
	transferEncrypt(c, symmetryKey)
}

func transferEncrypt(c *gin.Context, symmetryKey string) {
	// 响应体加密
	var (
		algo = c.GetHeader(constants.HeaderCryptSymmetryAlgorithm)
	)
	if wt, ok := c.Writer.(*ResponseCopyBodyWriter); ok && algo != "" {
		c.Header(constants.HeaderCryptSymmetryAlgorithm, algo)
		respBody := wt.body.String()
		newRespBody := respBody
		if respBody != "" {
			encryptRespBody, err := crypt.EncryptBase64(algo, symmetryKey, respBody)
			if err != nil {
				logger.Infof("encrypt response body err %s", err)
				return
			}
			newRespBody = fmt.Sprintf(`{"content":"%s"}`, encryptRespBody)
		}
		_, _ = wt.ResponseWriter.WriteString(newRespBody)
		wt.body.Reset()
		return
	}

}
func transferDecrypt(c *gin.Context, symmetryKey string) bool {
	// 请求体解密
	var (
		originalBody, err = io.ReadAll(c.Request.Body)
		algo              = c.GetHeader(constants.HeaderCryptSymmetryAlgorithm)
		newBody           = originalBody
	)
	if err != nil {
		controller.Fail(c, 400, "读取请求体内容失败", nil)
		return false
	}
	if len(originalBody) > 0 {
		jsonMap := make(map[string]interface{})
		err = json.Unmarshal(originalBody, &jsonMap)
		if err != nil {
			controller.Fail(c, 400, "读取请求体内容转为json失败", nil)
			return false
		}
		encryptStr := ctype.ToString(jsonMap["content"])
		if encryptStr != "" {
			newBody, err = crypt.DecryptBase64(algo, symmetryKey, encryptStr)
			if len(newBody) == 0 {
				controller.Fail(c, 400, "读取请求体内容解密失败", nil)
				return false
			}
		}
	}
	c.Request.Body = io.NopCloser(bytes.NewReader(newBody))
	return true
}

// 是否需要解密
func needDec(c *gin.Context) bool {
	// 判断gin请求是否是 json 请求
	// 通过请求头判断是否需要解密
	if c.GetHeader(constants.HeaderCryptSymmetryAlgorithm) != "" {
		return true
	}
	return false
}

// 获取传输加密的key，如果在header中获取一定是经过非对称加密的
// 从token绑定中的key获取密钥，优先级第一
func getSymmetryKey(c *gin.Context) (string, error) {

	var (
		ct    = ctx.GetContext(c)
		tkStr = ct.GetContextToken()
	)
	//先获取与token绑定的key
	if tkStr != "" {
		tk, err := user.LoadTokenByToken(tkStr)
		if err != nil {
			return "", err
		}
		if tk.Key != "" {
			return tk.Key, nil
		}
	}

	//如果token中没有绑定过key，那么就从header中获取
	key, err := security.GetTransferCryptKey(c)
	return key, err
}
