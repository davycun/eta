package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/davycun/eta/pkg/common/broker"
	"github.com/davycun/eta/pkg/common/config"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/monitor"
	"github.com/davycun/eta/pkg/common/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

func startServer(app *global.Application) {
	//start monitor
	var (
		cfg = app.GetConfig()
	)
	monitor.Start(fmt.Sprintf("%s:%d", "", cfg.Monitor.GetPort()))
	monitor.StartSentry()
	//start http server
	server := &http.Server{
		Addr:     ":" + strconv.Itoa(app.GetConfig().Server.Port),
		Handler:  global.GetGin(),
		ErrorLog: logger.Logger.Logger,
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go watchSignal(app, server, &wg)
	tip(cfg)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf(fmt.Sprintf("server start error: %s", err.Error()))
	}
	wg.Wait()
}

func tip(conf *config.Configuration) {
	tips := fmt.Sprintf("欢迎使用eta, 您可以通过%s 来查看帮助!", utils.FmtTextRed(" -h"))
	fmt.Println(tips)
	fmt.Println("服务已经启动，您可以通过如下地址访问: ")
	pt := strconv.Itoa(conf.Server.Port)
	wlan := `http://` + utils.GetLocalHost() + ":" + pt
	fmt.Printf("local: %s\n", utils.FmtUrl("http://localhost:"+pt))
	fmt.Printf("wlan: %s\n", utils.FmtUrl(wlan))
}

// TODO 这里也需要支持hook
func watchSignal(app *global.Application, server *http.Server, group *sync.WaitGroup) {
	defer group.Done()
	quit := make(chan os.Signal, 1)
	//Ctrl-C -> SIGINT   Ctrl-\ -> SIGQUIT
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
	s := <-quit
	logger.Info(fmt.Sprintf("Server shutdown with the signal %s", s.String()))

	// call shutdown before
	err := callLifeCycleHook(Shutdown, BeforeStage)
	if err != nil {
		logger.Errorf("Call lifecycle stage before Shutdown err %s", err)
	}

	//shutdown everything
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	var done = make(chan struct{}, 1)
	defer cancel()
	go func() {
		if err1 := server.Shutdown(ctx); err1 != nil {
			logger.Infof("Server Shutdown err: %s", err1)
		} else {
			logger.Info("Server Shutdown ok")
		}
		if err1 := app.Shutdown(ctx); err1 != nil {
			logger.Infof("Application Shutdown err: %s", err1)
		} else {
			logger.Info("Application Shutdown ok")
		}

		broker.CloseAll(ctx)

		done <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		logger.Info("Server shutdown timeout")
	case <-done:
	}

	_ = callLifeCycleHook(Shutdown, AfterStage)

}
