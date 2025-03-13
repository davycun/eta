package cache_test

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/http_tes"
	"testing"
)

var (
	key = "delta:ut:123"
)

func TestCache(t *testing.T) {
	http_tes.Call(t, set(), detail(), scan(), del())
	http_tes.Call(t, set1(), detail1())
	http_tes.Call(t, set2(), detail2())
	http_tes.Call(t, set3(), detail3())
	http_tes.Call(t, detail4())
}

func set() http_tes.HttpCase {
	b := map[string]any{
		"key": key,
		//"value": `{"kkk":"xxxx","log_level":4,"slow_threshold":200}`,
		"value":      map[string]any{"kkk": "xxxx", "log_level": 4, "slow_threshold": 200},
		"expiration": 7200,
	}

	return http_tes.HttpCase{
		Desc:     "缓存set",
		Method:   "POST",
		Path:     "/cache/set",
		Headers:  map[string]string{"Content-Type": "application/json"},
		Body:     b,
		ShowBody: true,
		Code:     "200",
	}
}

func detail() http_tes.HttpCase {
	return http_tes.HttpCase{
		Desc:    "缓存detail",
		Method:  "GET",
		Path:    "/cache/detail/" + key,
		Headers: map[string]string{"Content-Type": "application/json"},
		//Body:     b,
		ShowBody: true,
		Code:     "200",
	}
}

func set1() http_tes.HttpCase {
	b := map[string]any{
		"key":        "delta:ut:set1",
		"value":      2342423,
		"expiration": 567,
	}

	return http_tes.HttpCase{
		Desc:     "缓存set1",
		Method:   "POST",
		Path:     "/cache/set",
		Headers:  map[string]string{"Content-Type": "application/json"},
		Body:     b,
		ShowBody: true,
		Code:     "200",
	}
}

func detail1() http_tes.HttpCase {
	return http_tes.HttpCase{
		Desc:    "缓存detail1",
		Method:  "GET",
		Path:    "/cache/detail/delta:ut:set1",
		Headers: map[string]string{"Content-Type": "application/json"},
		//Body:     b,
		ShowBody: true,
		Code:     "200",
	}
}

func set2() http_tes.HttpCase {
	b := map[string]any{
		"key":   "delta:ut:set2",
		"value": "dfasdfasdf22323rds",
		//"value":      "null",
		"expiration": 567,
	}

	return http_tes.HttpCase{
		Desc:     "缓存set2",
		Method:   "POST",
		Path:     "/cache/set",
		Headers:  map[string]string{"Content-Type": "application/json"},
		Body:     b,
		ShowBody: true,
		Code:     "200",
	}
}

func detail2() http_tes.HttpCase {
	return http_tes.HttpCase{
		Desc:    "缓存detail2",
		Method:  "GET",
		Path:    "/cache/detail/delta:ut:set2",
		Headers: map[string]string{"Content-Type": "application/json"},
		//Body:     b,
		ShowBody: true,
		Code:     "200",
	}
}

func set3() http_tes.HttpCase {
	b := map[string]any{
		"key":        "delta:ut:set3",
		"value":      nil,
		"expiration": 567,
	}

	return http_tes.HttpCase{
		Desc:     "缓存set3",
		Method:   "POST",
		Path:     "/cache/set",
		Headers:  map[string]string{"Content-Type": "application/json"},
		Body:     b,
		ShowBody: true,
		Code:     "200",
	}
}

func detail3() http_tes.HttpCase {
	return http_tes.HttpCase{
		Desc:    "缓存detail3",
		Method:  "GET",
		Path:    "/cache/detail/delta:ut:set3",
		Headers: map[string]string{"Content-Type": "application/json"},
		//Body:     b,
		ShowBody: true,
		Code:     "200",
	}
}

func detail4() http_tes.HttpCase {
	return http_tes.HttpCase{
		Desc:    "缓存detail4",
		Method:  "GET",
		Path:    "/cache/detail/delta:ut:not_exists",
		Headers: map[string]string{"Content-Type": "application/json"},
		//Body:     b,
		ShowBody: true,
		Code:     "200",
	}
}

func scan() http_tes.HttpCase {
	return http_tes.HttpCase{
		Desc:    "缓存scan",
		Method:  "POST",
		Path:    "/cache/scan",
		Headers: map[string]string{"Content-Type": "application/json"},
		Body: `{
					"cursor": 0,
					"match": "*",
					"count": 10
				}`,
		ShowBody: true,
		Code:     "200",
	}
}

func del() http_tes.HttpCase {
	return http_tes.HttpCase{
		Desc:    "缓存del",
		Method:  "POST",
		Path:    "/cache/del",
		Headers: map[string]string{"Content-Type": "application/json"},
		Body: fmt.Sprintf(`{
					"keys": [ "%s", "123", "xxx:*", "*ksd*" ]
				}`, key),
		ShowBody: true,
		Code:     "200",
	}
}
