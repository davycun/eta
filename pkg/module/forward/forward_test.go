package forward_test

import (
	"context"
	"fmt"
	"github.com/davycun/eta/pkg/common/http_tes"
	_ "github.com/davycun/eta/pkg/common/http_tes"
	"github.com/davycun/eta/pkg/module/forward"
	"github.com/davycun/eta/pkg/module/setting"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var (
	vendor = "forward_test"
)

func TestMain(t *testing.M) {

	//启动backendServer
	backServer := &http.Server{Addr: "127.0.0.1:8111", Handler: http.DefaultServeMux}
	go backServer.ListenAndServe()
	prepareData()
	t.Run()
	backServer.Shutdown(context.Background())
}

func prepareData() {
	//注册代理
	forward.AddDefaultVendor(vendor, forward.Vendor{
		BaseCredentials: setting.BaseCredentials{BaseUrl: "http://127.0.0.1:8111"},
		Cache:           true,
		CacheUri: []string{
			"GET@/api/data/.*",
		},
	})
	forward.AddDefaultVendor("baidu", forward.Vendor{BaseCredentials: setting.BaseCredentials{BaseUrl: "https://www.baidu.com"}})
	dir := os.TempDir()
	err := os.WriteFile(dir+"/text.txt", []byte("this is a text file"), 0750)
	if err != nil {
		panic(err)
	}
	http.DefaultServeMux.HandleFunc("GET /file/{name}", readFile)
	http.DefaultServeMux.HandleFunc("POST /file/upload", uploadFile)
	http.DefaultServeMux.HandleFunc("GET /api/data/cache", testCache)
}

func readFile(w http.ResponseWriter, r *http.Request) {
	dir := os.TempDir()
	value := r.PathValue("name")
	dt, err := os.ReadFile(filepath.Join(dir, string(os.PathSeparator), value))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(dt)
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	dt, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	cd := r.Header.Get("Content-Disposition")
	fileName := strings.Split(strings.Split(cd, ";")[1], "=")[1]

	path := filepath.Join(os.TempDir(), strings.TrimSpace(fileName))
	os.WriteFile(path, dt, 0750)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(path))
}

func testCache(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("this is /api/data/cache test result"))
}
func TestForwardCache(t *testing.T) {
	http_tes.Call(t, http_tes.HttpCase{
		Method: "GET",
		Path:   fmt.Sprintf("/forward/%s/api/data/cache", vendor),
		ValidateFunc: []http_tes.ValidateFunc{
			func(t *testing.T, resp *http_tes.Response) {
				assert.Equal(t, "this is /api/data/cache test result", string(resp.RawBody))
			},
		},
	})
	http_tes.Call(t, http_tes.HttpCase{
		Method: "GET",
		Path:   fmt.Sprintf("/forward/%s/api/data/cache", vendor),
		ValidateFunc: []http_tes.ValidateFunc{
			func(t *testing.T, resp *http_tes.Response) {
				assert.Equal(t, "this is /api/data/cache test result", string(resp.RawBody))
				assert.NotEmpty(t, resp.Header.Get("X-Cache-Key"))
			},
		},
	})
}

func TestReadFile(t *testing.T) {
	http_tes.Call(t, http_tes.HttpCase{
		Method: "GET",
		Path:   fmt.Sprintf("/forward/%s/file/text.txt", vendor),
		ValidateFunc: []http_tes.ValidateFunc{
			func(t *testing.T, resp *http_tes.Response) {
				assert.Equal(t, "this is a text file", string(resp.RawBody))
			},
		},
	})
}

func TestUploadFile(t *testing.T) {
	content := "upload new file"
	http_tes.Call(t, http_tes.HttpCase{
		Method: "POST",
		Path:   fmt.Sprintf("/forward/%s/file/upload", vendor),
		Headers: map[string]string{
			"Content-Type":        "application/octet-stream",
			"Content-Disposition": `attachment; filename=eta_forward_upload_test.txt`,
		},
		Body: []byte(content),
		ValidateFunc: []http_tes.ValidateFunc{
			func(t *testing.T, resp *http_tes.Response) {
				file, err := os.ReadFile(string(resp.RawBody))
				assert.Nil(t, err)
				assert.Equal(t, content, string(file))
			},
		},
	})
}
func TestBaidu(t *testing.T) {
	http_tes.Call(t, http_tes.HttpCase{
		Method: "GET",
		Path:   fmt.Sprintf("/forward/%s/", "baidu"),
		ValidateFunc: []http_tes.ValidateFunc{
			func(t *testing.T, resp *http_tes.Response) {
				assert.Greater(t, len(resp.RawBody), 0)
			},
		},
	})
}
