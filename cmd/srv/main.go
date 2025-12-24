package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/financial_advisor/app/config"
	"github.com/financial_advisor/cmd/internal/db"
	"github.com/financial_advisor/cmd/internal/httpserver"
	"github.com/financial_advisor/cmd/internal/shutdown"
	"github.com/financial_advisor/cmd/internal/worker"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

func main() {
	rootCtx, rootCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer rootCancel()

	shutdown.Init(rootCtx, rootCancel)

	// shutdown all tasks gracefully on exit
	defer shutdown.Get().ShutdownTasks()

	cfg := config.Get()

	// global context is used to manage goroutines orphans
	cfg.SetGlobalCtx(rootCtx)

	if err := envconfig.Process("", cfg); err != nil {
		logrus.WithError(err).Fatal("Failed to process envconfig")
	}

	// Initialize and start the HTTP serverdocker
	if err := httpserver.Init(rootCtx); err != nil {
		logrus.WithError(err).Errorln("start HTTP server")

		return
	}

	// Initialize database connections
	if err := db.Init(rootCtx); err != nil {
		logrus.WithError(err).Errorln("initialize database connections")

		return
	}

	// Initialize workers
	if err := worker.Init(rootCtx); err != nil {
		logrus.WithError(err).Errorln("initialize workers")

		return
	}

	// wait for termination signal
	shutdown.Get().WaitForSignals()
}
