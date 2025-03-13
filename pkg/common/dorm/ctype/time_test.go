package ctype_test

import (
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTime(t *testing.T) {

	tl := []string{`"2024-07-09T16:25:30.892437+08:00"`, `"2024-07-09T16:25:30+08:00"`}

	for _, v := range tl {
		tm := &ctype.LocalTime{}
		err := tm.UnmarshalJSON([]byte(v))
		assert.Nil(t, err)
		assert.Equal(t, "2024-07-09T16:25:30", tm.Data.Format("2006-01-02T15:04:05"))
	}
}
