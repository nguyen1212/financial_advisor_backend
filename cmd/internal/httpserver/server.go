// Package httpserver serves http request
package httpserver

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/financial_advisor/app/config"
	framework "github.com/financial_advisor/app/external/framework/gin"
	"github.com/financial_advisor/cmd/internal/shutdown"
	"github.com/sirupsen/logrus"
)

type serverTask struct {
	srv *http.Server
}

var (
	ongoingCtx       context.Context
	ongoingCtxCancel context.CancelFunc
)

func Init(
	ctx context.Context,
) error {
	if shutdown.Get().IsShuttingDown() {
		return shutdown.ErrAborted
	}

	var (
		task    = &serverTask{}
		handler = framework.Handler()
	)

	// server self-manages its context
	ongoingCtx, ongoingCtxCancel = context.WithCancel(context.Background())

	task.srv = &http.Server{
		Addr:    ":" + config.Get().Port,
		Handler: handler,
		// all requests share the same ongoingCtx
		BaseContext: func(_ net.Listener) context.Context { return ongoingCtx },
	}

	go func() {
		logrus.WithField("Port", config.Get().Port).Infoln("start http server")
		if err := task.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.WithError(err).Errorln("start http server")

			shutdown.Get().Shutdown()
		}
	}()

	shutdown.Get().Add(task)

	return nil
}

func (task *serverTask) Name() string {
	return "http_server"
}

func (task *serverTask) Shutdown() error {
	serverShutdownCtx, serverShutdownCancel := context.WithTimeout(context.Background(), config.ShutdownPeriod)
	defer serverShutdownCancel()

	// stop server from receiving new requests
	err := task.srv.Shutdown(serverShutdownCtx)

	// cancel ongoing requests context
	ongoingCtxCancel()

	if err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}

	return nil
}
