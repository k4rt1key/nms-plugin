package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	logFile *os.File
	logger  *log.Logger
	once    sync.Once
)

// Initialize sets up the logger
func Initialize() {
	once.Do(func() {
		// Create logs directory if it doesn't exist
		logsDir := "logs"
		if _, err := os.Stat(logsDir); os.IsNotExist(err) {
			if err := os.Mkdir(logsDir, 0755); err != nil {
				log.Fatalf("Failed to create logs directory: %v", err)
			}
		}

		// Create log file with timestamp
		timestamp := time.Now().Format("2006-01-02_15-04-05")
		logFilePath := filepath.Join(logsDir, fmt.Sprintf("winrm-plugin_%s.log", timestamp))

		var err error
		logFile, err = os.Create(logFilePath)
		if err != nil {
			log.Fatalf("Failed to create log file: %v", err)
		}

		// Log only to file (no stdout)
		logger = log.New(logFile, "", log.LstdFlags|log.Lshortfile)

		Info("Logger initialized")
	})
}

// Close closes the log file
func Close() {
	if logFile != nil {
		Info("Closing logger")
		logFile.Close()
	}
}

// Info logs an informational message
func Info(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	logger.Printf("[INFO] %s", msg)
	// fmt.Printf("[INFO] %s\n", msg) // Mute stdout
}

// Error logs an error message
func Error(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	logger.Printf("[ERROR] %s", msg)
	// fmt.Printf("[ERROR] %s\n", msg) // Mute stdout
}

// Debug logs a debug message
func Debug(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	logger.Printf("[DEBUG] %s", msg)
	// fmt.Printf("[DEBUG] %s\n", msg) // Mute stdout
}

// Fatal logs a fatal message and exits
func Fatal(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	logger.Printf("[FATAL] %s", msg)
	// fmt.Printf("[FATAL] %s\n", msg) // Mute stdout
	os.Exit(1)
}
