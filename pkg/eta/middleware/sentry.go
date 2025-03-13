package middleware

import (
	"github.com/davycun/eta/pkg/eta/constants"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

func SentryRequestId(c *gin.Context) {
	if hub := sentrygin.GetHubFromContext(c); hub != nil {
		hub.Scope().SetTag(constants.HeaderRequestId, c.GetString(constants.HeaderRequestId))
	}
}
