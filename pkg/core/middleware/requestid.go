package middleware

import (
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
)

func RequestId(c *gin.Context) {
	rid := c.GetHeader(constants.HeaderRequestId)
	if rid == "" {
		rid = ulid.Make().String() // ulid 能解析出时间，方便按照时间追溯
		c.Header(constants.HeaderRequestId, rid)
	}
	c.Set(constants.HeaderRequestId, rid)
	c.Next()
}
