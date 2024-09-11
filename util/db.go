package util

import (
	"fmt"
	"github.com/xdire/temporal-async/messaging"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	DB_USER     = "1"
	DB_PASSWORD = ""
	DB_NAME     = "temporal_example"
)

type DB struct {
	gorm *gorm.DB
}

func (db *DB) Conn() *gorm.DB {
	return db.gorm
}

func (db *DB) Connect() error {
	dbc, err := gorm.Open(sqlite.Open("store.db"), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect database")
	}

	err = dbc.AutoMigrate(&messaging.Process{})
	if err != nil {
		return fmt.Errorf("failed to migrate database")
	}

	db.gorm = dbc

	return nil
}
