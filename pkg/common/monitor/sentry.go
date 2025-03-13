package monitor

import (
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/duke-git/lancet/v2/maputil"
	"github.com/getsentry/sentry-go"
	"log"
	"os"
	"time"
)

var (
	EnvInfo = map[string]string{
		"commit_tag":  os.Getenv("COMMIT_TAG"),
		"commit_hash": os.Getenv("COMMIT_HASH"),
		"target_arch": os.Getenv("TARGET_ARCH"),
	}
)

func SentryEnable() bool {
	return global.GetConfig().Monitor.SentryEnable
}
func StartSentry() {
	if !SentryEnable() {
		logger.Infof("Sentry disabled!")
		return
	}
	conf := global.GetConfig()
	logger.Infof("Sentry config: %v", conf.Monitor)
	err := sentry.Init(sentry.ClientOptions{
		Dsn:                   conf.Monitor.SentryDsn,
		TracesSampleRate:      conf.Monitor.SentryTracesSampleRate,
		EnableTracing:         conf.Monitor.SentryEnableTracing,
		SendDefaultPII:        conf.Monitor.SentrySendDefaultPii,
		Environment:           conf.Server.Env,
		BeforeSend:            beforeSend,
		BeforeSendTransaction: beforeSend,
		//Debug:                 true,
	})
	if err != nil {
		log.Fatalf("Sentry init error! %s", err)
	}

	defer sentry.Flush(2 * time.Second)
	defer sentry.Recover()

	logger.Infof("Sentry init success!")
}

func beforeSend(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
	event.Tags = maputil.Merge(event.Tags, EnvInfo)
	return event
}
