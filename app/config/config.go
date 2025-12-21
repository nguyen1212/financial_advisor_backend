// Package config defines all the configuration for the app
package config

import (
	"context"
	"sync"
	"time"
)

const (
	// these properties are highly dependent on terminationGracePeriodSeconds set in k8s deployment
	// default is 30s
	ShutdownPeriod      = 15 * time.Second
	ShutdownHardPeriod  = 5 * time.Second
	ReadinessDrainDelay = 5 * time.Second
)

type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
)

type Config struct {
	globalCtx context.Context

	ENV        string `envconfig:"ENV" default:"development"`
	AppName    string `envconfig:"APP_NAME" default:"financial_advisor"`
	AppVersion string `envconfig:"APP_VERSION" default:"v1"`

	Port string   `envconfig:"HTTP_PORT" default:"40000"`
	Cors []string `envconfig:"CORS_ALLOWED_HOSTS" default:"*"`

	LogLevel LogLevel `envconfig:"LOG_LEVEL" default:"info"`

	MySQLUser     string `envconfig:"MYSQL_USER" default:"admin"`
	MySQLPassword string `envconfig:"MYSQL_PASSWORD" default:"root123"`
	MySQLHost     string `envconfig:"MYSQL_HOST" default:"localhost"`
	MySQLPort     string `envconfig:"MYSQL_PORT" default:"40001"`
	MySQLDatabase string `envconfig:"MYSQL_DATABASE" default:"financial_advisor"`
}

var (
	singleton *Config
	once      sync.Once
)

func init() {
	once.Do(func() {
		singleton = &Config{
			globalCtx: context.Background(),
		}
	})
}

func (cfg *Config) SetGlobalCtx(ctx context.Context) {
	if ctx == nil {
		return
	}

	cfg.globalCtx = ctx
}

func Get() *Config {
	return singleton
}
