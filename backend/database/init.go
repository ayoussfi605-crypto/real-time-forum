package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Init() error {
	var err error

	DB, err = sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	DB.SetMaxOpenConns(1) 
	if _, err = DB.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return fmt.Errorf("enable foreign keys: %w", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	schema, err := os.ReadFile("./database/schema.sql")
	if err != nil {
		return fmt.Errorf("read schema.sql: %w", err)
	}

	if _, err = DB.Exec(string(schema)); err != nil {
		return fmt.Errorf("execute schema.sql: %w", err)
	}

	fmt.Println("Database initialized successfully")
	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}