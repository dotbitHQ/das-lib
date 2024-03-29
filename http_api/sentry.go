package http_api

import (
	"github.com/getsentry/sentry-go"
	"time"
)

func SentryInit(dsn string) (err error) {
	//if dsn == "" {
	//	return fmt.Errorf("sentry dsn is empty")
	//}
	err = sentry.Init(sentry.ClientOptions{
		Dsn:           dsn,
		EnableTracing: true,
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 0,
	})
	return

}

func RecoverPanic() {
	if err := recover(); err != nil {
		sentry.CurrentHub().Recover(err)
		sentry.Flush(time.Second * 2)
		panic(err)
	}
}

func CaptureMessage(msg string) {
	sentry.CurrentHub().CaptureMessage(msg)
	sentry.Flush(time.Second * 2)
}
