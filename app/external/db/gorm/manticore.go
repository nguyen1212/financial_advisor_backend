// Package gorm
package gorm

import (
	"database/sql"
	"fmt"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Manticore struct {
	db *gorm.DB
}

var (
	manticoreOnce     sync.Once
	manticoreInstance *Manticore
)

func OpenManticoreConnection(conn string, config gorm.Config) error {
	var err error

	manticoreOnce.Do(func() {
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

		manticoreInstance = &Manticore{db: db}
	})

	if err != nil {
		return err
	}

	return nil
}

func GetManticoreIns() *Manticore {
	return manticoreInstance
}

func (m *Manticore) DB() *gorm.DB {
	return m.db
}

func (m *Manticore) Close() error {
	sqlDB, err := m.db.DB()
	if err != nil {
		return fmt.Errorf("get sql db of manticore: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("close manticore sql db: %w", err)
	}

	return nil
}
