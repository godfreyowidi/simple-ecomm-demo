package db

import (
	"database/sql"
	"fmt"
	"os"
)

var DB *sql.DB

func InitDB() error {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		return fmt.Errorf("DATABASE_URL not set")
	}
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	return DB.Ping()
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
