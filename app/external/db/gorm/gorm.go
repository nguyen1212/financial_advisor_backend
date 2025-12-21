// Package gorm
package gorm

import (
	"database/sql"
	"fmt"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	once      sync.Once
	singleton *gorm.DB
)

func OpenDBConnection(conn string, config gorm.Config) (*sql.DB, error) {
	var err error

	once.Do(func() {
		sqlDB, innerErr := sql.Open("mysql", conn)
		if innerErr != nil {
			err = fmt.Errorf("sql open: %w", innerErr)

			return
		}

		db, innerErr := gorm.Open(
			mysql.New(mysql.Config{Conn: sqlDB}),
			&config,
		)
		if innerErr != nil {
			err = fmt.Errorf("gorm open: %w", innerErr)

			return
		}

		singleton = db
	})

	if err != nil {
		return nil, err
	}

	return singleton.DB()
}

func Get() *gorm.DB {
	return singleton
}
