package menu

var (
	defaultMenu = make([]Menu, 0, 10)
)

// Registry
// 默认初始化的菜单
func Registry(menuList ...Menu) {
	defaultMenu = append(defaultMenu, menuList...)
}
