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
)

var (
	DBPath          = "./news.db"
	SERVER_PORT     = ":9999"
	APP_ENV         = "development"
	LOG_LEVEL       = "info"
	MAX_UPLOAD_SIZE = int64(10485760)
)

func init() {
	// Load from environment variables
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
