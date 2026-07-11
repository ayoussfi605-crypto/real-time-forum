package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // Import the SQLite driver
)

var Conn *sql.DB
var err error

func ConnToDb() (*sql.DB, error) {
	Conn, err = sql.Open("sqlite3", "forum.db")
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}
	_, err = Conn.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return nil, fmt.Errorf("error enabling foreign key support: %v", err)
	}
	// ping the database connection to ensure it's valid
	if err := Conn.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %v", err)
	}
	return Conn, nil
}
