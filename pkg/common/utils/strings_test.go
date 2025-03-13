package utils_test

import (
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHumpToUnderline(t *testing.T) {
	tt := []struct {
		src    string
		target string
	}{
		{"Name", "name"}, {"AgE", "ag_e"}, {"SEX", "sex"}, {"ABCaBD", "abca_bd"},
		{"PeopleType", "people_type"}, {"12Name", "12_name"}, {"中国", "中国"},
	}

	for _, v := range tt {
		assert.Equal(t, v.target, utils.HumpToUnderline(v.src))
	}
}
func TestUnderlineToHump(t *testing.T) {
	tt := []struct {
		src    string
		target string
	}{
		{"name", "Name"}, {"ag_e", "AgE"}, {"sex", "Sex"}, {"abca_bd", "AbcaBd"},
		{"people_type", "PeopleType"}, {"12_name", "12Name"}, {"中国", "中国"},
	}

	for _, v := range tt {
		assert.Equal(t, v.target, utils.UnderlineToHump(v.src))
	}
}

func TestSearch(t *testing.T) {

	str := "小写，未婚,本市户籍，13917283763,未婚"

	split := utils.Split(str, ",", "，")
	assert.Equal(t, []string{"小写", "未婚", "本市户籍", "13917283763", "未婚"}, split)

	split = utils.Split(str, ",")
	assert.Equal(t, []string{"小写，未婚", "本市户籍，13917283763", "未婚"}, split)

	split = utils.Split(str, "@")
	assert.Equal(t, []string{"小写，未婚,本市户籍，13917283763,未婚"}, split)

}
