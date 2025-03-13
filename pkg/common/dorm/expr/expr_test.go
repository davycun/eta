package expr

import (
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExplainValue(t *testing.T) {
	var (
		ds = []struct {
			Value            any
			ContainsExpected string
			InExpected       string
		}{
			{
				Value:            34,
				ContainsExpected: "34",
				InExpected:       "34",
			},
			{
				Value:            34.8765,
				ContainsExpected: "34.8765",
				InExpected:       "34.8765",
			},
			{
				Value:            "sijdn28好",
				ContainsExpected: "'sijdn28好'",
				InExpected:       "'sijdn28好'",
			},
			{
				Value:            []string{"a", "b", "c"},
				ContainsExpected: `'a' or 'b' or 'c'`,
				InExpected:       `('a','b','c')`,
			},
			{
				Value:            []int{13, 67, 9876},
				ContainsExpected: `13 or 67 or 9876`,
				InExpected:       `(13,67,9876)`,
			},
		}
	)

	for _, v := range ds {
		assert.Equal(t, v.InExpected, ExplainExprValue(dorm.DaMeng, v.Value))
	}

}
