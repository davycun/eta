package hook_test

import (
	"github.com/davycun/eta/pkg/core/service/hook"
	"reflect"
	"testing"
)

func TestType(t *testing.T) {
	et := reflect.TypeOf(hook.SrvConfig{})

	x := reflect.New(et).Interface()
	y := reflect.New(reflect.SliceOf(et)).Interface()

	if tmp, ok := y.(*[]hook.SrvConfig); ok {
		println(len(*tmp))
	}
	if tmp, ok := x.(*hook.SrvConfig); ok {
		println(tmp.GetTableName())
	}
}
