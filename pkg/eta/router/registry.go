package router

type (
	RouteFunc func()
)

var (
	routerFuncList = make([]RouteFunc, 0)
)

// Registry 注册路由函数，一个简单的空函数，实际添加路由由调用者自行决定
func Registry(rf RouteFunc) {
	routerFuncList = append(routerFuncList, rf)
}
