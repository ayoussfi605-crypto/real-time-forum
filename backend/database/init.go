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
	// NOTE: MaxOpenConns(1) here would deadlock any handler that opens a
	// second query while the first is still open (e.g. fetching categories
	// for each post in a loop). Allow a small pool instead, and use WAL mode
	// so concurrent reads/writes on SQLite are safe.
	DB.SetMaxOpenConns(10) // Allow up to 10 simultaneous database connections

	// Enable WAL mode (Write-Ahead Logging) to allow reading and writing at the same time safely
	if _, err = DB.Exec("PRAGMA journal_mode=WAL; PRAGMA foreign_keys=ON"); err != nil {
		return fmt.Errorf("enable WAL/foreign keys: %w", err)
	}
	if _, err = DB.Exec("PRAGMA journal_mode = WAL"); err != nil {
		return fmt.Errorf("enable WAL mode: %w", err)
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