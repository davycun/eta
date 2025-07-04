package forward

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNeedCache(t *testing.T) {

	vd := Vendor{
		Cache: true,
		CacheUri: []string{
			"GET@^/api/v1/vendor/[a-zA-Z0-9]+/upload$",
		},
	}
	assert.True(t, vd.NeedCache("GET", "/api/v1/vendor/123/upload"))
	assert.False(t, vd.NeedCache("POST", "/api/v1/vendor/123/upload"))
}

func TestVendor_NeedCache_ExactMatch(t *testing.T) {
	vendor := Vendor{
		Cache:    true,
		CacheUri: []string{"GET@/api/a/b"},
		Sorted:   false,
	}

	if !vendor.NeedCache("GET", "/api/a/b") {
		t.Errorf("Expected to cache GET /api/a/b")
	}
}

func TestVendor_NeedCache_WildcardMatch(t *testing.T) {
	vendor := Vendor{
		Cache:    true,
		CacheUri: []string{"*@/api/a/.*"},
		Sorted:   false,
	}

	if !vendor.NeedCache("POST", "/api/a/b") {
		t.Errorf("Expected to cache POST /api/a/b")
	}
}

func TestVendor_NeedCache_MethodMismatch(t *testing.T) {
	vendor := Vendor{
		Cache:  true,
		Sorted: false,
	}

	if !vendor.NeedCache("POST", "/api/a/b") {
		t.Errorf("Expected not to cache POST /api/a/b")
	}
}
func TestVendor_NeedCache_MethodMismatch2(t *testing.T) {
	vendor := Vendor{
		Cache:    true,
		CacheUri: []string{"GET@/api/a/.*"},
		Sorted:   false,
	}

	if vendor.NeedCache("GET", "/api/b") {
		t.Errorf("Expected not to cache POST /api/a/b")
	}
}

func TestVendor_SortCacheUri(t *testing.T) {
	vendor := Vendor{
		Cache:    true,
		CacheUri: []string{"GET@/api/a/c", "GET@/api/a/b/d", "GET@/api/a/b"},
		Sorted:   false,
	}

	vendor.sortCacheUri()

	expectedOrder := []string{"GET@/api/a/b/d", "GET@/api/a/c", "GET@/api/a/b"}
	for i, uri := range vendor.CacheUri {
		if uri != expectedOrder[i] {
			t.Errorf("Expected URI at index %d to be %s, got %s", i, expectedOrder[i], uri)
		}
	}
}
