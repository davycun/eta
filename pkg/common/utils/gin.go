package utils

import "github.com/gin-gonic/gin"

func GetUrlPath(c *gin.Context) string {
	if c == nil || c.Request == nil {
		return ""
	}
	return c.Request.URL.Path
}
func GetHttpMethod(c *gin.Context) string {
	if c == nil || c.Request == nil {
		return ""
	}
	return c.Request.Method
}
