package tag

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJson(t *testing.T) {
	tg := NewTag("json", SplitComma, "")
	tg.Add("pt_code", "").Add("omitempty", "omitempty")

	assert.Equal(t, `json:"pt_code,omitempty"`, tg.String())

	tg = ParseTag("gorm", `column:name;comment:这是什么`, SplitSemicolon, SplitColon)
	assert.Equal(t, `gorm:"column:name;comment:这是什么"`, tg.String())

	tg = NewTag("binding", SplitComma, SplitEq).Add("required", "").Add("default", "123").Add("oneof", "1 2 3 ''")
	assert.Equal(t, `binding:"required,default=123,oneof=1 2 3 ''"`, tg.String())

	tg = NewTag("es", SplitSemicolon, SplitColon).Add("type", "keyword").Add("analyzer", "digit_analyzer")
	assert.Equal(t, `es:"type:keyword;analyzer:digit_analyzer"`, tg.String())

}
