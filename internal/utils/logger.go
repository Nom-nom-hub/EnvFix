package utils

import (
	"fmt"
	"os"
	"time"
)

type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarningLevel
	ErrorLevel
)

// Logger is unused and can be removed in future refactoring
// type Logger struct {
// 	level LogLevel
// 	file  *os.File
// }

var logFile *os.File

func InitLogger(filepath string) error {
	var err error
	logFile, err = os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	return err
}

func CloseLogger() {
	if logFile != nil {
		logFile.Close()
	}
}

func logMessage(level string, msg string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	formatted := fmt.Sprintf(msg, args...)
	output := fmt.Sprintf("[%s] %s: %s\n", timestamp, level, formatted)

	// Print to stdout
	fmt.Print(output)

	// Write to log file
	if logFile != nil {
		_, _ = logFile.WriteString(output)
	}
}

func Debug(msg string, args ...interface{}) {
	logMessage("DEBUG", msg, args...)
}

func Info(msg string, args ...interface{}) {
	logMessage("INFO", msg, args...)
}

func Warning(msg string, args ...interface{}) {
	logMessage("WARN", msg, args...)
}

func Error(msg string, args ...interface{}) {
	logMessage("ERROR", msg, args...)
}

func Success(msg string, args ...interface{}) {
	logMessage("✓", msg, args...)
}

func Fatal(msg string, args ...interface{}) {
	logMessage("FATAL", msg, args...)
	os.Exit(1)
}
