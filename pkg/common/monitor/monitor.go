package monitor

import (
	"expvar"
	"github.com/davycun/eta/pkg/common/logger"
	"net/http"
	_ "net/http/pprof"
	"runtime"
)

func Start(addr string) {
	expvar.Publish("metrics", expvar.Func(routineCount))
	go func() {
		if err := http.ListenAndServe(addr, http.DefaultServeMux); err != nil {
			logger.Errorf("start monitor error %s", err.Error())
		}
	}()
}

type Metrics struct {
	NumCPU       int    `json:"num_cpu"`
	NumGoroutine int    `json:"num_goroutine"`
	NumCgoCall   int64  `json:"num_cgo_call"`
	GoVersion    string `json:"go_version"`
	EtaVersion   string `json:"eta_version"`
}

func routineCount() any {
	mts := Metrics{
		NumCPU:       runtime.NumCPU(),
		NumGoroutine: runtime.NumGoroutine(),
		NumCgoCall:   runtime.NumCgoCall(),
		GoVersion:    runtime.Version(),
	}
	return mts
}
