package db

import (
	"database/sql"
	"fmt"
)

func CreateTables(Conn *sql.DB) error {

	// Create the users table
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		Nickname TEXT NOT NULL UNIQUE,
		First_Name TEXT NOT NULL,
		Last_Name TEXT NOT NULL,
		Age INTEGER NOT NULL,
		Gender TEXT NOT NULL,
		Email TEXT NOT NULL UNIQUE, 
		Password TEXT NOT NULL,
		Created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`,

		`CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE
	)`,

		`CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	)`,
		`CREATE TABLE IF NOT EXISTS post_categories (
		post_id INTEGER NOT NULL,
		category_id INTEGER NOT NULL,
		PRIMARY KEY (post_id, category_id),
		FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
		FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
	)`,
		`CREATE TABLE IF NOT EXISTS sessions (
		id INTEGER PRIMARY KEY AUTOINCREMENT ,
 		user_id int NOT NULL,
		token text UNIQUE NOT NULL,
		expiration_date datetime NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id)
	)`,
		`CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	)`,
		`CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		sender_id INTEGER NOT NULL,
		receiver_id INTEGER NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (receiver_id) REFERENCES users(id) ON DELETE CASCADE
	)`,
}
	for _, query := range queries {
		_, err := Conn.Exec(query)
		if err != nil {
			return fmt.Errorf("error creating table: %v", err)
		}
	}
	fmt.Println("Tables created successfully")

	_, err = Conn.Exec(`INSERT OR IGNORE INTO categories (name) VALUES
(?), (?), (?), (?), (?), (?), (?) ON CONFLICT(name) DO NOTHING`,
		"sport", "games", "filmes", "kitchen", "news", "market", "others")
	if err != nil {
		return fmt.Errorf("error inserting into categories: %v", err)
	}
	return nil
}
