// Package gorm
package gorm

import (
	"database/sql"
	"fmt"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQL struct {
	db *gorm.DB
}

var (
	once      sync.Once
	singleton *MySQL
)

func OpenMySQLConnection(conn string, config gorm.Config) error {
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

		sqlDB, _ = db.DB()

		if err = sqlDB.Ping(); err != nil {
			err = fmt.Errorf("ping DB connection: %w", err)

			return
		}

		singleton = &MySQL{db: db}
	})

	if err != nil {
		return err
	}

	return nil
}

func GetMySQLIns() *MySQL {
	return singleton
}

func (m *MySQL) DB() *gorm.DB {
	return m.db
}

func (m *MySQL) Close() error {
	sqlDB, err := m.db.DB()
	if err != nil {
		return fmt.Errorf("get sql db of mysql: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("close mysql sql db: %w", err)
	}

	return nil
}
