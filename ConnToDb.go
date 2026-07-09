package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // Import the PostgreSQL driver
)

var Conn *sql.DB

func ConnToDb() (*sql.DB, error) {
Conn, err := sql.Open("sqlite3", "forum.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	return Conn, nil
}
