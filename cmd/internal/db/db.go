// Package db handles database connection
package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/financial_advisor/app/config"
	gormAdapter "github.com/financial_advisor/app/external/db/gorm"
	"github.com/financial_advisor/cmd/internal/shutdown"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type task struct{ sqlDB *sql.DB }

func Init(_ context.Context) error {
	if shutdown.Get().IsShuttingDown() {
		return shutdown.ErrAborted
	}

	var (
		task       = &task{}
		cfg        = config.Get()
		gormConfig = gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		}
		mysqlConn = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?multiStatements=true&parseTime=true&charset=utf8mb4",
			cfg.MySQLUser,
			cfg.MySQLPassword,
			cfg.MySQLHost,
			cfg.MySQLPort,
			cfg.MySQLDatabase,
		)
	)

	if cfg.LogLevel == config.LogLevelDebug {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}

	// open DB connection by external adapter
	db, err := gormAdapter.OpenDBConnection(mysqlConn, gormConfig)
	if err != nil {
		return fmt.Errorf("open DB connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("ping DB connection: %w", err)
	}

	logrus.
		WithField("database", cfg.MySQLDatabase).
		WithField("log_level", cfg.LogLevel).
		Infoln("database connection established")

	task.sqlDB = db
	shutdown.Get().Add(task)

	return nil
}

func (task *task) Name() string {
	return "database"
}

func (task *task) Shutdown() error {
	return task.sqlDB.Close()
}
