package ctx

import (
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/gin-gonic/gin"
	"strings"
)

func GetToken(c *Context) string {
	token := c.GetContextToken()
	if token != "" {
		return token
	}

	token = c.GetGinContext().GetHeader(constants.HeaderAuthorization)
	if token == "" {
		token = c.GetGinContext().Query(strings.ToLower(constants.HeaderAuthorization))
		if token == "" {
			token = c.GetGinContext().Query(constants.HeaderAuthorization)
		}
	}
	if token == "" {
		token, _ = c.GetGinContext().Cookie(constants.HeaderAuthorization)
		if token == "" {
			token, _ = c.GetGinContext().Cookie(strings.ToLower(constants.HeaderAuthorization))
		}
	}
	if token != "" {
		c.SetContextToken(token)
	}

	return token
}

//var (
//	threadLocal = routine.NewInheritableThreadLocal[*Context]()
//)

// Deprecated: not support
func GetCurrentContext() (*Context, bool) {
	//c := threadLocal.Get()
	return nil, false
}

// Deprecated: not support
func SetCurrentContext(c *Context) {
	//threadLocal.Set(c)
}
func CleanCurrentContext(c *gin.Context) {
	cleanContext(c)
	//当routine结束的时候，会被自动回收，这里无需手动删除
	//threadLocal.Remove()
}

func cleanContext(c *gin.Context) {

	//只需要解除gin.Context 和ctx.Context 的引用关系即可，其余的不需要，会被自动回收的
	ct := GetContext(c)
	delete(c.Keys, DeltaGinContextKey)
	ct.keys.Delete(DeltaGinContextKey)

	//ct.keys.Delete(GormContextKey)
	//ct.keys.Delete(GormContextKey)
	//ct.keys.Delete(GormAppKey)
	//ct.keys.Delete(UserIdContextKey)
	//ct.keys.Delete(AppIdContextKey)
	//ct.keys.Delete(CurrentDeptIdContextKey)
	//ct.keys.Delete(TokenContextKey)
	//ct.keys.Delete(StorageContextKey)
}
