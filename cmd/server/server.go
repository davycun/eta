package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/davycun/eta/pkg/common/broker"
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/config"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/monitor"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/eta/middleware"
	"github.com/davycun/eta/pkg/eta/migrator"
	"github.com/davycun/eta/pkg/eta/router"
	"github.com/davycun/eta/pkg/eta/validator"
	"github.com/davycun/eta/pkg/module"
	"github.com/gin-gonic/gin/binding"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

var (
	StartCommand = &cobra.Command{
		Use:   "server",
		Short: "启动一个eta服务",
		Args: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
	confFile  = ""
	argConfig = config.Configuration{}
)

func init() {
	config.BindArgConfig(StartCommand, &confFile, &argConfig)
}

func run() error {
	destCfg := config.LoadConfig(confFile, &argConfig)
	//如果json串中有一个字段是数值，但是反序列化的时候针对这个字段没有指定具体的是float或者int
	//那么默认json会反序列化为float64类型，这也就是为什么我用map去接受反序列化的时候，明明序列化之前是int，但是反序列化后map里面是float64的原因
	//如果设置了EnableDecoderUseNumber，那么这种情况下反序列化的目标就会被指定为json.Number对象（其实是个string，type Number string）
	//binding.EnableDecoderUseNumber = true
	binding.EnableDecoderDisallowUnknownFields = true
	err := caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			return callStageCallback(BeforeInitApplication)
		}).
		Call(func(cl *caller.Caller) error {
			global.InitApplication(destCfg)
			return nil
		}).
		Call(func(cl *caller.Caller) error {
			return callStageCallback(AfterInitApplication)
		}).
		Call(func(cl *caller.Caller) error {
			return callStageCallback(BeforeInitMiddleware)
		}).
		Call(func(cl *caller.Caller) error {
			//eta.InitEta()
			//初始化模块需放第一
			module.InitModules()
			validator.AddValidate()
			middleware.InitMiddleware()
			return nil
		}).
		Call(func(cl *caller.Caller) error {
			return callStageCallback(AfterInitMiddleware)
		}).
		Call(func(cl *caller.Caller) error {
			return callStageCallback(BeforeInitRouter)
		}).
		Call(func(cl *caller.Caller) error {
			router.InitRouter()
			return nil
		}).
		Call(func(cl *caller.Caller) error {
			return callStageCallback(AfterInitRouter)
		}).
		Call(func(cl *caller.Caller) error {
			return callStageCallback(BeforeMigrate)
		}).
		Call(func(cl *caller.Caller) error {
			return migrator.MigrateLocal(global.GetLocalGorm())
		}).
		Call(func(cl *caller.Caller) error {
			return callStageCallback(AfterMigrate)
		}).
		Call(func(cl *caller.Caller) error {
			return callStageCallback(BeforeStartServer)
		}).
		Call(func(cl *caller.Caller) error {
			startServer(global.GetApplication())
			return nil
		}).
		Call(func(cl *caller.Caller) error {
			return callStageCallback(AfterStartServer)
		}).Err
	return err
}

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

	_ = callStageCallback(BeforeShutdown)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	var done = make(chan struct{}, 1)
	defer cancel()
	go func() {
		if err := server.Shutdown(ctx); err != nil {
			logger.Info(fmt.Sprintf("Server Shutdown err: %s", err.Error()))
		} else {
			logger.Info("Server Shutdown ok")
		}
		if err := app.Shutdown(ctx); err != nil {
			logger.Info(fmt.Sprintf("Application Shutdown err: %s", err.Error()))
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

	_ = callStageCallback(AfterShutdown)
}
