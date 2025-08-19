package dgorm

import "gorm.io/gorm"

var (
	gormPlugins = make([]gorm.Plugin, 0)
)

func Registry(plugin gorm.Plugin) {
	if plugin == nil {
		return
	}
	gormPlugins = append(gormPlugins, plugin)
}
