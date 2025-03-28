package forward_test

import (
	"context"
	"fmt"
	"github.com/davycun/eta/pkg/common/global"
	_ "github.com/davycun/eta/pkg/common/http_tes"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/module/forward"
	"github.com/davycun/eta/pkg/module/setting"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

var (
	vendor = "forward_test"
)

func prepareData(t *testing.T) {

	//注册代理
	forward.AddDefaultVendor(vendor, setting.BaseCredentials{BaseUrl: "http://127.0.0.1:8111"})

	dir := os.TempDir()
	err := os.WriteFile(dir+"/test.txt", []byte("this is a text file"), 0750)
	assert.Nil(t, err)
	http.DefaultServeMux.HandleFunc("GET /file/{name}", processFile)
}

func processFile(w http.ResponseWriter, r *http.Request) {
	dir := os.TempDir()
	value := r.PathValue("name")
	dt, err := os.ReadFile(filepath.Join(dir, string(os.PathSeparator), value))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(dt)
}

func startServer(app *global.Application) {
	server := &http.Server{
		Addr:     ":" + strconv.Itoa(app.GetConfig().Server.Port),
		Handler:  global.GetGin(),
		ErrorLog: logger.Logger.Logger,
	}

	server.ListenAndServe()
}

func TestJson(t *testing.T) {

	prepareData(t)
	server := &http.Server{Addr: "127.0.0.1:8111", Handler: http.DefaultServeMux}
	defer server.Shutdown(context.Background())
	go startServer(global.GetApplication())

	clt := resty.New()
	clt.SetBaseURL("http://127.0.0.1:8080")
	r := clt.R()

	r.SetHeader("Content-Type", "application/json")
	resp, err := r.Get(fmt.Sprintf("/forward/%s/file/text.txt", vendor))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
	assert.Equal(t, "this is a text file", string(resp.Body()))

}
