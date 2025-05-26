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

	logger *log.Logger

	once sync.Once
)

func Initialize() {

	once.Do(func() {

		logsDir := "logs"

		if _, err := os.Stat(logsDir); os.IsNotExist(err) {

			if err := os.Mkdir(logsDir, 0755); err != nil {
				log.Fatalf("Failed to create logs directory: %v", err)
			}

		}

		timestamp := time.Now().Format("2006-01-02_15-04-05")

		logFilePath := filepath.Join(logsDir, fmt.Sprintf("winrm-plugin_%s.log", timestamp))

		var err error

		logFile, err = os.Create(logFilePath)

		if err != nil {
			log.Fatalf("Failed to create log file: %v", err)
		}

		logger = log.New(logFile, "", log.LstdFlags|log.Lshortfile)

		Info("Logger initialized")

	})
}

func Close() {

	if logFile != nil {

		Info("Closing logger")

		logFile.Close()

	}
}

func Info(format string, args ...interface{}) {

	msg := fmt.Sprintf(format, args...)

	logger.Printf("[INFO] %s", msg)

}

func Error(format string, args ...interface{}) {

	msg := fmt.Sprintf(format, args...)

	logger.Printf("[ERROR] %s", msg)

}

func Debug(format string, args ...interface{}) {

	msg := fmt.Sprintf(format, args...)

	logger.Printf("[DEBUG] %s", msg)

}

func Fatal(format string, args ...interface{}) {

	msg := fmt.Sprintf(format, args...)

	logger.Printf("[FATAL] %s", msg)

	os.Exit(1)

}
