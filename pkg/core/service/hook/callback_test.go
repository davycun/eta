package hook_test

import (
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/stretchr/testify/assert"
	"reflect"
	"slices"
	"testing"
)

var (
	mp = map[string][]hook.Callback{}
)

func TestFuncEqual(t *testing.T) {
	tbName := "t_people"
	addFunc(tbName)
	delFunc(tbName, firstFunc)
	cb, ok := mp[tbName]
	assert.True(t, ok)
	assert.Equal(t, 1, len(cb))
}

func addFunc(tbName string) {
	mp[tbName] = []hook.Callback{
		firstFunc,
		secondFunc,
	}
}

func delFunc(tbName string, fc hook.Callback) {
	cb := mp[tbName]
	cb = slices.DeleteFunc(cb, func(callback hook.Callback) bool {
		tp := reflect.TypeOf(callback)
		tp2 := reflect.TypeOf(fc)
		logger.Infoln(tp.PkgPath())
		logger.Infoln(tp2.PkgPath())
		val := reflect.ValueOf(fc)
		val2 := reflect.ValueOf(callback)
		logger.Infoln(val.Pointer())
		logger.Infoln(val2.Pointer())
		if val.Pointer() == val2.Pointer() {
			return true
		}
		return false
	})
	mp[tbName] = cb
}

func firstFunc(cfg *hook.SrvConfig, pos hook.CallbackPosition) error {

	return nil
}
func secondFunc(cfg *hook.SrvConfig, pos hook.CallbackPosition) error {

	return nil
}
