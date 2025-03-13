package httpclient_test

import (
	"encoding/json"
	"fmt"
	"github.com/davycun/eta/pkg/common/httpclient"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	server = httptest.NewServer(http.DefaultServeMux)
)

type Person struct {
	Name    string
	Age     int
	Deleted bool
}

func TestJSON(t *testing.T) {
	http.HandleFunc("/test", func(writer http.ResponseWriter, request *http.Request) {
		p := Person{
			Name:    "davy",
			Age:     23,
			Deleted: true,
		}
		marshal, err := json.Marshal(p)
		if err == nil {
			fmt.Fprint(writer, string(marshal))
		}
	})
	var p Person

	err := httpclient.DefaultHttpClient.Url(server.URL + "/test").Method("GET").Do(&p).Error
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, p.Name, "davy")
	assert.Equal(t, p.Age, 23)
	assert.Equal(t, p.Deleted, true)
}
