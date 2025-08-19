package plugin_crypt

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestPrefix(t *testing.T) {
	data := []struct {
		str    string
		rs     string
		expect bool
	}{

		{"@@347629307620757504@t_address@@$e$asdfasdfasdfasdfasdfadfasdfasdf", "@@347629307620757504@t_address@@", true},
	}

	pattern := "@@[0-9a-zA-z]+@[0-9a-zA-z]+@@"
	compile, err := regexp.Compile(pattern)
	assert.Nil(t, err)

	for _, v := range data {
		assert.Equal(t, v.expect, compile.MatchString(v.str))
		findString := compile.FindString(v.str)
		assert.Equal(t, v.rs, findString)
	}
}
