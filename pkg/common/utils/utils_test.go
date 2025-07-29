package utils_test

import (
	"bytes"
	"fmt"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
	"unsafe"
)

func TestUtils(t *testing.T) {

	var mp map[string]interface{}
	var sl []string
	var i int
	var b bool
	var s string
	var it interface{}
	var st struct{}
	var c chan struct{}

	assert.True(t, utils.IsZero(mp))
	assert.True(t, utils.IsZero(sl))
	assert.True(t, utils.IsZero(i))
	assert.True(t, utils.IsZero(b))
	assert.True(t, utils.IsZero(s))
	assert.True(t, utils.IsZero(it))
	assert.True(t, utils.IsZero(st))
	assert.True(t, utils.IsZero(c))

	assert.True(t, utils.IsZero(nil))
	assert.True(t, utils.IsZero(""))
}

func TestColor(t *testing.T) {

	t.Log(fmt.Sprintf("这个季度的总营收是%s\n", utils.FmtTextRed("500万")))
	t.Log(fmt.Sprintf("访问这个网站: %s\n", utils.FmtUrl("https://www.datlas.com")))
	t.Log(fmt.Sprintf("这个是重点: %s\n", utils.FmtColor("重要的", 4, utils.TextGreen, 0)))
}

func TestLocalIp(t *testing.T) {

	t.Log(utils.GetLocalHost())
}

func TestTime(t *testing.T) {

	format := time.Now().Format("2006-01-02 15:04:05Z07:00")
	format2 := time.Now().Format(utils.TimeLayout)
	println(format)
	println(format2)
}

func TestContainAll(t *testing.T) {
	assert.True(t, utils.ContainAll([]string{"a", "b", "c"}, "a", "b", "c"))
	assert.False(t, utils.ContainAll([]string{"a", "b", "c"}, "a", "b", "c", "d"))
	assert.True(t, utils.ContainAll([]string{"a", "b", "c"}, "a", "b"))
	assert.False(t, utils.ContainAll([]string{"a", "b", "c"}, "a", "b", "d"))
	assert.True(t, utils.ContainAll([]int{1, 2, 3}, 1))
	assert.True(t, utils.ContainAll([]int{1, 2, 3}, 2, 3, 1))
}
func TestContainAny(t *testing.T) {
	assert.True(t, utils.ContainAny([]string{"a", "b", "c"}, "a", "b", "c"))
	assert.True(t, utils.ContainAny([]int{4, 6, 9}, 410, 5, 9))
	assert.False(t, utils.ContainAny([]string{"a", "b"}, ""))
}

func TestMap(t *testing.T) {
	var mp map[string][]string
	ids := mp["test"]
	ids = append(ids, "1")
	assert.Equal(t, 1, len(ids))
}

func TestSlice(t *testing.T) {

	var ids []string
	ids = append(ids, "2")
	assert.Equal(t, 1, len(ids))
}

func StringToBytes1(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}
func StringToBytes2(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

var L = 1024 * 1024
var s = bytes.Repeat([]byte{'a'}, L)

func BenchmarkStringToBytesStandard(b *testing.B) {
	var str = strings.Repeat("a", L)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = []byte(str)
		//bt := []byte(str)
		//if len(bt) != L {
		//	b.Fatal()
		//}
	}
}
func BenchmarkStringToBytes(b *testing.B) {
	var str = strings.Repeat("b", L)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = utils.StringToBytes(str)
		//bt := utils.StringToBytes(str)
		//if len(bt) != L {
		//	b.Fatal()
		//}
	}
}
func BenchmarkStringToBytes1(b *testing.B) {
	var str = strings.Repeat("c", L)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = StringToBytes1(str)
		//bt := StringToBytes1(str)
		//if len(bt) != L {
		//	b.Fatal()
		//}
	}
}
func BenchmarkStringToBytes2(b *testing.B) {
	var str = strings.Repeat("d", L)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = StringToBytes2(str)
		//bt := StringToBytes2(str)
		//if len(bt) != L {
		//	b.Fatal()
		//}
	}
}

func TestContainsAnyInsensitive(t *testing.T) {
	dt := []string{
		"Access-Control-Allow-Methods",
		"access-control-allow-origin",
		"Access-Control-Allow-Origin",
		"Access-Control-Allow-Credentials",
		"Cache-control",
	}
	mp := map[string]bool{
		"cache-control":                true,
		"abc":                          false,
		"Access-Control-Allow-Origin":  true,
		"access-control-allow-methods": true,
		"123&^%$":                      false,
	}

	for k, v := range mp {
		assert.Equal(t, v, utils.ContainAnyInsensitive(dt, k))
	}
}
