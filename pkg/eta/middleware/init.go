package middleware

import (
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/core/middleware"
	"github.com/davycun/eta/pkg/eta/plugin/plugin_crypt"
	"github.com/davycun/eta/pkg/module/menu/menu_srv"
	"github.com/davycun/eta/pkg/module/optlog"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"slices"
)

func InitMiddleware() error {
	Registry(MidOption{Name: "gin_log", Order: 0, HandlerFunc: gin.LoggerWithConfig(newGinLogConfig())})
	Registry(MidOption{Name: "health", Order: 1, HandlerFunc: middleware.Health})
	Registry(MidOption{Name: "stats", Order: 2, HandlerFunc: middleware.Stats})
	Registry(MidOption{Name: "error_handler", Order: 3, HandlerFunc: middleware.ErrorHandler})
	Registry(MidOption{Name: "error", Order: 4, HandlerFunc: middleware.RequestId})
	Registry(MidOption{Name: "authorize", Order: 6, HandlerFunc: Auth})
	Registry(MidOption{Name: "api_auth", Order: 7, HandlerFunc: menu_srv.ApiCallAuth})
	Registry(MidOption{Name: "table", Order: 10, HandlerFunc: LoadTable})
	Registry(MidOption{Name: "opt_log", Order: 20, HandlerFunc: optlog.Log})
	Registry(MidOption{Name: "crypto", Order: 30, HandlerFunc: plugin_crypt.TransferCrypt})
	Registry(MidOption{Name: "sentry1", Order: 40, HandlerFunc: sentrygin.New(sentrygin.Options{Repanic: true})})
	Registry(MidOption{Name: "sentry2", Order: 50, HandlerFunc: SentryRequestId})

	mds := sortMiddleware()
	for _, v := range mds {
		global.GetGin().Use(v.HandlerFunc)
	}
	return nil
}

func sortMiddleware() []MidOption {
	mds := make([]MidOption, 0, len(middlewareMap))
	for _, v := range middlewareMap {
		mds = append(mds, v)
	}
	slices.SortFunc(mds, func(a, b MidOption) int {
		return a.Order - b.Order
	})
	return mds
}
