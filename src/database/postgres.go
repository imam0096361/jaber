package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB
var DBType string

// Connect initializes database based on environment variable
func Connect() {
	dbType := os.Getenv("DB_TYPE")
	if dbType == "" {
		dbType = "sqlite"
	}

	var err error

	if dbType == "postgres" {
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

		err = InitPostgres(host, port, user, password, dbname)
	} else {
		dbPath := os.Getenv("DB_PATH")
		if dbPath == "" {
			dbPath = "./news.db"
		}
		err = InitSQLite(dbPath)
	}

	if err != nil {
		log.Fatal(err)
	}
}

// InitPostgres initializes PostgreSQL database connection
func InitPostgres(host, port, user, password, dbname string) error {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return fmt.Errorf("database connection error: %w", err)
	}

	// Test the connection
	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("database ping error: %w", err)
	}

	DBType = "postgres"
	log.Println("✅ PostgreSQL database connected successfully")

	// Create tables
	CreateTables()

	return nil
}

// InitSQLite initializes SQLite database connection (fallback)
func InitSQLite(dbPath string) error {
	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("database connection error: %w", err)
	}

	DBType = "sqlite"
	log.Println("✅ SQLite database connected successfully")

	// Create tables
	CreateTables()

	return nil
}

// CreateTables creates all necessary tables
func CreateTables() {
	if DBType == "postgres" {
		createPostgresTables()
	} else {
		createSQLiteTables()
	}
}

func createPostgresTables() {
	articlesQuery := `
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
	CREATE INDEX IF NOT EXISTS idx_articles_created ON articles(created);
	`

	_, err := DB.Exec(articlesQuery)
	if err != nil {
		log.Printf("Error creating articles table: %v", err)
	}

	usersQuery := `
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

	_, err = DB.Exec(usersQuery)
	if err != nil {
		log.Printf("Error creating users table: %v", err)
	}

	log.Println("✅ PostgreSQL tables created/verified")
}

func createSQLiteTables() {
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
		log.Printf("Table creation error: %v", err)
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
		log.Printf("Table creation error: %v", err)
	}

	log.Println("✅ SQLite tables created/verified")
}

// Close closes the database connection
func Close() {
	if DB != nil {
		DB.Close()
	}
}

// InsertSampleData inserts sample data into the database
func InsertSampleData() {
	// Check if articles already exist
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM articles").Scan(&count)
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

	for _, article := range articles {
		query := `INSERT INTO articles (title, content, category, author, featured) VALUES ($1, $2, $3, $4, $5)`
		_, err := DB.Exec(query, article.title, article.content, article.category, article.author, article.featured)
		if err != nil {
			log.Printf("Error inserting article: %v", err)
		}
	}

	log.Println("✅ Sample data inserted successfully")
}
