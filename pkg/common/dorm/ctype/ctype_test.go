package ctype_test

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestType(t *testing.T) {

	m := &ctype.Map{}
	v := reflect.ValueOf(m)
	assert.Equal(t, v.Kind(), reflect.Pointer)
}

func TestSet(t *testing.T) {
	assert.Equal(t, "true", fmt.Sprintf("%t", true))
	assert.Equal(t, "false", fmt.Sprintf("%t", false))
}
