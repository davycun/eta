package setting

import (
	"encoding/json"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/eta/constants"
)

var SafetyConfigJson = `
{
	"autoLoginOut": {
		"currentTime": 30,
		"active": false,
		"maxTimestamp": 60,
		"minTimestamp": 2
	}
}`

func defaultSafetyConfig() Setting {
	st := Setting{}
	st.Namespace = constants.NamespaceEta
	st.Category = "SafetyConfig"
	st.Name = "SafetyConfig"
	st.Remark = "配置系统安全等，前端多久不操作就自动退出等"
	var SafetyConfig any
	_ = json.Unmarshal([]byte(SafetyConfigJson), &SafetyConfig)
	st.Content = ctype.NewJson(SafetyConfig)
	return st
}
