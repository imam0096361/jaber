package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Connect() {
	var err error
	DB, err = sql.Open("sqlite3", "./news.db")
	if err != nil {
		log.Fatal("Database connection error:", err)
	}

	// Create tables
	CreateTables()
}

func CreateTables() {
	articlesQuery := `
	CREATE TABLE IF NOT EXISTS articles (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		category TEXT NOT NULL,
		author TEXT NOT NULL,
		image TEXT,
		created DATETIME DEFAULT CURRENT_TIMESTAMP,
		featured BOOLEAN DEFAULT 0
	);
	`
	_, err := DB.Exec(articlesQuery)
	if err != nil {
		log.Fatal("Table creation error:", err)
	}

	usersQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err = DB.Exec(usersQuery)
	if err != nil {
		log.Fatal("Users table creation error:", err)
	}
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}
