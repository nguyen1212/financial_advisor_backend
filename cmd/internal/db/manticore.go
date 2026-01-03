// Package db handles database connection
package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/financial_advisor/app/config"
	gormAdapter "github.com/financial_advisor/app/external/db/gorm"
	"github.com/financial_advisor/cmd/internal/shutdown"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type manticoreTask struct {
	db *gormAdapter.Manticore
}

func InitManticore(_ context.Context) error {
	if shutdown.Get().IsShuttingDown() {
		return shutdown.ErrAborted
	}

	var (
		task          = &manticoreTask{}
		cfg           = config.Get()
		manticoreConn = fmt.Sprintf(
			"admin:admin@tcp(%s:%s)/%s?multiStatements=true&parseTime=true&charset=utf8mb4",
			cfg.ManticoreHost,
			cfg.ManticorePort,
			cfg.ManticoreDatabase,
		)
		gormConfig = gorm.Config{
			Logger: gormAdapter.NewLogger(
				log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
					SlowThreshold:             50 * time.Millisecond,
					LogLevel:                  logger.Warn,
					IgnoreRecordNotFoundError: false,
					Colorful:                  true,
					ParameterizedQueries:      cfg.ENV == config.ENVProduction,
				},
				cfg.ManticoreDatabase,
			),
		}
	)

	if cfg.LogLevel == config.LogLevelDebug {
		gormConfig.Logger = gormConfig.Logger.LogMode(logger.Info)
	}

	// open DB connection by external adapter
	err := gormAdapter.OpenManticoreConnection(manticoreConn, gormConfig)
	if err != nil {
		return fmt.Errorf("open DB connection: %w", err)
	}

	task.db = gormAdapter.GetManticoreIns()

	logrus.
		WithField("database", cfg.ManticoreDatabase).
		WithField("log_level", cfg.LogLevel).
		Infoln("manticore database connection established")

	shutdown.Get().Add(task)

	return nil
}

func (task *manticoreTask) Name() string {
	return "searchengine_manticore"
}

func (task *manticoreTask) Shutdown() error {
	return task.db.Close()
}
