package database

import (
	"app/src/config"
	"app/src/model"
	"app/src/utils"
	"context"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Global database instance for reconnection
var globalDB *gorm.DB
var dsn string

// Connect establishes a database connection with retry logic
func Connect(dbHost, dbName string) *gorm.DB {
	dsn = fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Dhaka",
		dbHost, config.DBUser, config.DBPassword, dbName, config.DBPort,
	)

	var db *gorm.DB
	var err error

	// Retry connection up to 10 times with exponential backoff
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		utils.Log.Infof("Attempting database connection (attempt %d/%d)...", i+1, maxRetries)

		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger:                 logger.Default.LogMode(logger.Info),
			SkipDefaultTransaction: true,
			PrepareStmt:            false, // Disable prepared statements to avoid stale connection issues
			TranslateError:         true,
		})

		if err == nil {
			// Test the connection
			sqlDB, errDB := db.DB()
			if errDB == nil {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if pingErr := sqlDB.PingContext(ctx); pingErr == nil {
					utils.Log.Infof("✅ Database connection established successfully!")
					break
				} else {
					err = pingErr
				}
			} else {
				err = errDB
			}
		}

		utils.Log.Warnf("Failed to connect to database (attempt %d/%d): %+v", i+1, maxRetries, err)

		// Exponential backoff: 2s, 4s, 8s, 16s, 32s... (max 60s)
		waitTime := time.Duration(min(60, 2<<i)) * time.Second
		utils.Log.Infof("Waiting %v before retry...", waitTime)
		time.Sleep(waitTime)
	}

	if err != nil {
		utils.Log.Errorf("❌ Failed to connect to database after %d attempts: %+v", maxRetries, err)
		panic("Database connection failed - cannot start application")
	}

	sqlDB, errDB := db.DB()
	if errDB != nil {
		utils.Log.Errorf("Failed to get database connection pool: %+v", errDB)
		panic("Database connection pool failed")
	}

	// Configure connection pooling for stability
	sqlDB.SetMaxIdleConns(5)                    // Reduced idle connections
	sqlDB.SetMaxOpenConns(50)                   // Reduced max connections
	sqlDB.SetConnMaxLifetime(30 * time.Minute)  // Shorter lifetime to refresh connections
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)   // Close idle connections after 5 min

	// Run Migrations
	if err := db.AutoMigrate(&model.Article{}, &model.User{}); err != nil {
		utils.Log.Errorf("Failed to auto migrate: %+v", err)
	}

	// Store global reference
	globalDB = db

	// Start connection health monitor
	go monitorConnection(db)

	return db
}

// monitorConnection periodically checks and recovers database connection
func monitorConnection(db *gorm.DB) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		sqlDB, err := db.DB()
		if err != nil {
			utils.Log.Errorf("Connection monitor: Failed to get DB: %+v", err)
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err = sqlDB.PingContext(ctx)
		cancel()

		if err != nil {
			utils.Log.Warnf("Connection monitor: Database ping failed: %+v", err)
			utils.Log.Infof("Connection monitor: Attempting to recover connection...")

			// Close stale connections
			sqlDB.SetMaxIdleConns(0)
			time.Sleep(1 * time.Second)
			sqlDB.SetMaxIdleConns(5)

			// Try ping again
			ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
			if pingErr := sqlDB.PingContext(ctx2); pingErr == nil {
				utils.Log.Infof("Connection monitor: ✅ Connection recovered!")
			} else {
				utils.Log.Errorf("Connection monitor: ❌ Connection recovery failed: %+v", pingErr)
			}
			cancel2()
		}
	}
}

// min returns the smaller of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
