package authorize

import (
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/gin-gonic/gin"
)

func Authorization(c *gin.Context) {
	c.Header(constants.HeaderUserId, ctx.GetContext(c).GetContextUserId())
	c.Header(constants.HeaderAppId, ctx.GetContext(c).GetContextAppId())
}
