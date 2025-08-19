package setting

import (
	"github.com/davycun/eta/pkg/common/dorm/ctype"
)

var (
	defaultSettingMap = make(map[string]Setting) //category+name -> Setting
)

// Registry 添加默认的配置，在setting表初次创建的时候会把内容写入数据库
func Registry(settingList ...Setting) {
	for _, v := range settingList {
		defaultSettingMap[v.Category+v.Name] = v
	}
}

// GetDefault
// 可以提供默认值的获取后修改在Registry
func GetDefault[T any](category, name string) T {
	var (
		t     T
		s, ok = defaultSettingMap[category+name]
	)
	if !ok {
		return t
	}
	if ctype.IsValid(s.Content) {
		switch v := s.Content.Data.(type) {
		case T:
			t = v
		case *T:
			t = *v
		}
	}
	return t
}
