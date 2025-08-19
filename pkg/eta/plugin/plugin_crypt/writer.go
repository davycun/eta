package plugin_crypt

import (
	"bytes"
	"github.com/gin-gonic/gin"
)

type ResponseCopyBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w ResponseCopyBodyWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
}
