package cache_test

import (
	"github.com/davycun/eta/pkg/common/cache"
	_ "github.com/davycun/eta/pkg/common/http_tes"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCacheExpireKey(t *testing.T) {

	err := cache.SetEx("name", "davy", time.Second*10)
	assert.Nil(t, err)

	var s string
	err, b := cache.Get("name", &s)
	assert.True(t, b)
	assert.Nil(t, err)
	assert.Equal(t, "davy", s)

	//虽然redis过期了但是本地缓存还在
	time.Sleep(time.Second * 15)
	var s1 string
	err1, b1 := cache.Get("name", &s1)
	assert.True(t, b1)
	assert.Nil(t, err1)
	assert.Equal(t, "davy", s1)

	//一分钟后本地缓存也清空了
	time.Sleep(time.Second * 60)
	var s2 string
	err2, b2 := cache.Get("name", &s2)
	assert.False(t, b2)
	assert.Nil(t, err2)
	assert.Empty(t, s2)

}

func TestCacheDelKey(t *testing.T) {

	type Pep struct {
		Name string
		IdNo string
		Age  int
		Boy  bool
	}
	var (
		key = "pep:12"
		p1  = Pep{
			Name: "davy",
			IdNo: "382632839402723",
			Age:  32,
			Boy:  true,
		}
		p2 Pep
		p3 Pep
	)

	err := cache.Set(key, p1)
	assert.Nil(t, err)

	err, b := cache.Get(key, &p2)
	assert.True(t, b)
	assert.Nil(t, err)
	assert.Equal(t, "382632839402723", p2.IdNo)
	assert.Equal(t, 32, p2.Age)
	assert.Equal(t, true, p2.Boy)

	err, b = cache.Del(key)
	assert.Nil(t, err)
	assert.True(t, b)

	//一分钟后本地缓存也清空了
	err2, b2 := cache.Get(key, &p3)
	assert.False(t, b2)
	assert.Nil(t, err2)
	assert.Empty(t, p3.IdNo)
	assert.Empty(t, p3.Name)
}

func TestCacheDelKeyPattern(t *testing.T) {
	type Pep struct {
		Name string
		IdNo string
		Age  int
		Boy  bool
	}
	var (
		key  = "pep:"
		key1 = "pep:1"
		key2 = "pep:2"
		p1   = Pep{
			Name: "nnnnnnn",
			IdNo: "382632839402723",
			Age:  32,
			Boy:  true,
		}
		p2 Pep
		p3 Pep
	)

	err := cache.Set(key1, p1)
	assert.NoError(t, err)
	err = cache.Set(key2, p1)
	assert.NoError(t, err)

	err, b := cache.Get(key1, &p2)
	assert.NoError(t, err)
	assert.True(t, b)

	err, b = cache.DelKeyPattern(key + "*")
	assert.NoError(t, err)
	assert.True(t, b)

	err, b = cache.Get(key1, &p3)
	assert.NoError(t, err)
	assert.False(t, b)
	assert.Empty(t, p3.IdNo)
	assert.Empty(t, p3.Name)
}
