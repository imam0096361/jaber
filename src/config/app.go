package config

import (
	"os"
	"strconv"
)

const (
	SERVER_PORT_DEFAULT     = ":9999"
	UPLOADS_DIR             = "frontend/uploads"
	MAX_FILE_SIZE           = 10 << 20 // 10 MB
	MAX_UPLOAD_SIZE_DEFAULT = 10485760
	DB_TYPE_DEFAULT         = "sqlite"
)

var (
	DBPath          = "./news.db"
	DBType          = "sqlite"
	DBHost          = "localhost"
	DBPort          = "5432"
	DBUser          = "postgres"
	DBPassword      = ""
	DBName          = "news_portal"
	SERVER_PORT     = ":9999"
	APP_ENV         = "development"
	LOG_LEVEL       = "info"
	MAX_UPLOAD_SIZE = int64(10485760)
)

func init() {
	// Load from environment variables
	if dbType := os.Getenv("DB_TYPE"); dbType != "" {
		DBType = dbType
	}

	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		DBHost = dbHost
	}

	if dbPort := os.Getenv("DB_PORT"); dbPort != "" {
		DBPort = dbPort
	}

	if dbUser := os.Getenv("DB_USER"); dbUser != "" {
		DBUser = dbUser
	}

	if dbPassword := os.Getenv("DB_PASSWORD"); dbPassword != "" {
		DBPassword = dbPassword
	}

	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		DBName = dbName
	}

	if port := os.Getenv("APP_PORT"); port != "" {
		SERVER_PORT = ":" + port
	}

	if dbPath := os.Getenv("DB_PATH"); dbPath != "" {
		DBPath = dbPath
	}

	if env := os.Getenv("APP_ENV"); env != "" {
		APP_ENV = env
	}

	if level := os.Getenv("LOG_LEVEL"); level != "" {
		LOG_LEVEL = level
	}

	if maxSize := os.Getenv("MAX_UPLOAD_SIZE"); maxSize != "" {
		if size, err := strconv.ParseInt(maxSize, 10, 64); err == nil {
			MAX_UPLOAD_SIZE = size
		}
	}
}
