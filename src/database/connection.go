package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// Database holds the database instance
type Database struct {
	db *sql.DB
	dbType string
}

var db *Database

// GetDB returns the database instance
func GetDB() *Database {
	return db
}

// Init initializes the database connection
func Init() error {
	database := &Database{}
	
	dbType := os.Getenv("DB_TYPE")
	if dbType == "" {
		dbType = "sqlite"
	}

	var err error
	if dbType == "postgres" {
		err = database.initPostgres()
	} else {
		err = database.initSQLite()
	}

	if err != nil {
		return err
	}

	db = database
	return nil
}

func (d *Database) initPostgres() error {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "news_portal"
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	d.db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool parameters
	d.db.SetMaxOpenConns(25)
	d.db.SetMaxIdleConns(5)
	d.db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := d.db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}

	d.dbType = "postgres"
	log.Println("✅ Connected to PostgreSQL database")

	// Create tables and indexes
	if err := d.createTables(); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	// Insert sample data
	d.insertSampleData()

	return nil
}

func (d *Database) initSQLite() error {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./news.db"
	}

	var err error
	d.db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open SQLite database: %w", err)
	}

	// Set connection pool parameters
	d.db.SetMaxOpenConns(10)
	d.db.SetMaxIdleConns(2)
	d.db.SetConnMaxLifetime(10 * time.Minute)

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := d.db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping SQLite: %w", err)
	}

	d.dbType = "sqlite"
	log.Println("✅ Connected to SQLite database")

	// Create tables and indexes
	if err := d.createTables(); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	// Insert sample data
	d.insertSampleData()

	return nil
}

func (d *Database) createTables() error {
	if d.dbType == "postgres" {
		return d.createPostgresTables()
	}
	return d.createSQLiteTables()
}

func (d *Database) createPostgresTables() error {
	schema := `
	-- Articles Table
	CREATE TABLE IF NOT EXISTS articles (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		content TEXT NOT NULL,
		category VARCHAR(100) NOT NULL,
		author VARCHAR(100) NOT NULL,
		image VARCHAR(255),
		created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		featured BOOLEAN DEFAULT FALSE
	);

	CREATE INDEX IF NOT EXISTS idx_articles_category ON articles(category);
	CREATE INDEX IF NOT EXISTS idx_articles_featured ON articles(featured);
	CREATE INDEX IF NOT EXISTS idx_articles_created ON articles(created DESC);

	-- Users Table
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	`

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Execute schema
	_, err := d.db.ExecContext(ctx, schema)
	if err != nil {
		log.Printf("Warning: Error creating schema: %v", err)
		// Continue anyway, tables might already exist
	}

	log.Println("✅ PostgreSQL schema ready")
	return nil
}

func (d *Database) createSQLiteTables() error {
	schema := `
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

	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Execute schema
	_, err := d.db.ExecContext(ctx, schema)
	if err != nil {
		log.Printf("Warning: Error creating schema: %v", err)
		// Continue anyway, tables might already exist
	}

	log.Println("✅ SQLite schema ready")
	return nil
}

func (d *Database) insertSampleData() {
	// Check if articles already exist
	var count int
	err := d.db.QueryRow("SELECT COUNT(*) FROM articles").Scan(&count)
	if err != nil || count > 0 {
		return // Already has data
	}

	// Sample articles
	articles := []struct {
		title    string
		content  string
		category string
		author   string
		featured bool
	}{
		{
			title:    "দেশের নতুন অর্থনীতি নীতি",
			content:  "বাংলাদেশ সরকার নতুন অর্থনীতি নীতি ঘোষণা করেছে যা দেশের উন্নয়নে গুরুত্বপূর্ণ ভূমিকা রাখবে।",
			category: "জাতীয়",
			author:   "সিনিয়র রিপোর্টার",
			featured: true,
		},
		{
			title:    "ক্রিকেট চ্যাম্পিয়নশিপে বাংলাদেশের জয়",
			content:  "বাংলাদেশ ক্রিকেট দল আন্তর্জাতিক চ্যাম্পিয়নশিপে দারুণ পারফরম্যান্স দেখিয়েছে।",
			category: "খেলাধুলা",
			author:   "স্পোর্টস এডিটর",
			featured: false,
		},
		{
			title:    "কৃত্রিম বুদ্ধিমত্তার নতুন উদ্ভাবন",
			content:  "বিশ্বব্যাপী কৃত্রিম বুদ্ধিমত্তার ক্ষেত্রে নতুন উদ্ভাবনী প্রযুক্তি আবিষ্কৃত হয়েছে।",
			category: "প্রযুক্তি",
			author:   "টেক করেসপন্ডেন্ট",
			featured: true,
		},
		{
			title:    "স্বাস্থ্যসেবায় নতুন মাত্রা",
			content:  "দেশে স্বাস্থ্যসেবা খাতে যুগান্তকারী পরিবর্তন আসতে চলেছে যা সাধারণ মানুষের জীবনযাত্রার মান উন্নত করবে।",
			category: "স্বাস্থ্য",
			author:   "হেলথ করেসপন্ডেন্ট",
			featured: false,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for _, article := range articles {
		query := `INSERT INTO articles (title, content, category, author, featured) VALUES ($1, $2, $3, $4, $5)`
		_, err := d.db.ExecContext(ctx, query, article.title, article.content, article.category, article.author, article.featured)
		if err != nil {
			log.Printf("Warning: Error inserting article: %v", err)
		}
	}

	log.Println("✅ Sample data initialized")
}

// Query executes a query
func (d *Database) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return d.db.Query(query, args...)
}

// QueryRow executes a query that returns a single row
func (d *Database) QueryRow(query string, args ...interface{}) *sql.Row {
	return d.db.QueryRow(query, args...)
}

// Exec executes a query without returning rows
func (d *Database) Exec(query string, args ...interface{}) (sql.Result, error) {
	return d.db.Exec(query, args...)
}

// BeginTx starts a new transaction
func (d *Database) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return d.db.BeginTx(ctx, nil)
}

// Close closes the database connection
func (d *Database) Close() error {
	if d != nil && d.db != nil {
		return d.db.Close()
	}
	return nil
}

// GetConnection returns the raw sql.DB for advanced usage
func (d *Database) GetConnection() *sql.DB {
	return d.db
}

// GetDBType returns the database type (postgres or sqlite)
func (d *Database) GetDBType() string {
	return d.dbType
}
