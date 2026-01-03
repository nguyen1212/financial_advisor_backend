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

type mysqlTask struct{ db *gormAdapter.MySQL }

func InitMySQL(_ context.Context) error {
	if shutdown.Get().IsShuttingDown() {
		return shutdown.ErrAborted
	}

	var (
		task       = &mysqlTask{}
		cfg        = config.Get()
		gormConfig = gorm.Config{
			Logger: gormAdapter.NewLogger(
				log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
					SlowThreshold:             50 * time.Millisecond,
					LogLevel:                  logger.Warn,
					IgnoreRecordNotFoundError: false,
					Colorful:                  true,
					ParameterizedQueries:      cfg.ENV == config.ENVProduction,
				},
				cfg.MySQLDatabase,
			),
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
		gormConfig.Logger = gormConfig.Logger.LogMode(logger.Info)
	}

	// open DB connection by external adapter
	if err := gormAdapter.OpenMySQLConnection(mysqlConn, gormConfig); err != nil {
		return fmt.Errorf("open DB connection: %w", err)
	}

	logrus.
		WithField("database", cfg.MySQLDatabase).
		WithField("log_level", cfg.LogLevel).
		Infoln("database connection established")

	task.db = gormAdapter.GetMySQLIns()

	shutdown.Get().Add(task)

	return nil
}

func (task *mysqlTask) Name() string {
	return "database_mysql"
}

func (task *mysqlTask) Shutdown() error {
	return task.db.Close()
}
