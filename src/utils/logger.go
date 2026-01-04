package utils

import (
	"log"
)

var Logger *log.Logger

func init() {
	Logger = log.Default()
}

func LogError(msg string, err error) {
	if err != nil {
		Logger.Printf("[ERROR] %s: %v", msg, err)
	}
}

func LogInfo(msg string) {
	Logger.Printf("[INFO] %s", msg)
}
