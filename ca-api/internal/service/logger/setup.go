// logger/setup.go
package logger

import (
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/sirupsen/logrus"
)

var (
	ErrorLogger *logrus.Logger
	WarnLogger  *logrus.Logger
	InfoLogger  *logrus.Logger
	DebugLogger *logrus.Logger
)

func Init() {
	// Create log files
	errorFile, _ := os.OpenFile("logs/error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	warnFile, _ := os.OpenFile("logs/warn.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	infoFile, _ := os.OpenFile("logs/info.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	debugFile, _ := os.OpenFile("logs/debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	// Setup ERROR logger
	ErrorLogger = logrus.New()
	ErrorLogger.SetLevel(logrus.ErrorLevel)
	ErrorLogger.SetFormatter(&logrus.JSONFormatter{TimestampFormat: "2006-01-02T15:04:05.000Z"})
	ErrorLogger.SetOutput(io.MultiWriter(os.Stdout, errorFile)) // Console + error.log

	// Setup WARN logger
	WarnLogger = logrus.New()
	WarnLogger.SetLevel(logrus.WarnLevel)
	WarnLogger.SetFormatter(&logrus.JSONFormatter{TimestampFormat: "2006-01-02T15:04:05.000Z"})
	WarnLogger.SetOutput(io.MultiWriter(os.Stdout, warnFile)) // Console + warn.log

	// Setup INFO logger
	InfoLogger = logrus.New()
	InfoLogger.SetLevel(logrus.InfoLevel)
	InfoLogger.SetFormatter(&logrus.JSONFormatter{TimestampFormat: "2006-01-02T15:04:05.000Z"})
	InfoLogger.SetOutput(io.MultiWriter(os.Stdout, infoFile)) // Console + info.log

	// Setup DEBUG logger
	DebugLogger = logrus.New()
	DebugLogger.SetLevel(logrus.DebugLevel)
	DebugLogger.SetFormatter(&logrus.JSONFormatter{TimestampFormat: "2006-01-02T15:04:05.000Z"})
	DebugLogger.SetOutput(io.MultiWriter(os.Stdout, debugFile)) // Console + debug.log
}

// Helper functions
func getCallerInfo() (string, string, int) {
	pc, file, line, _ := runtime.Caller(2)
	funcName := runtime.FuncForPC(pc).Name()
	fileName := filepath.Base(file)
	return funcName, fileName, line
}

func Error(operation, message string, err error) {
	funcName, fileName, line := getCallerInfo()
	ErrorLogger.WithFields(logrus.Fields{
		"operation": operation,
		"function":  funcName,
		"file":      fileName,
		"line":      line,
		"error":     err.Error(),
	}).Error(message)
}

func Warn(operation, message string) {
	funcName, fileName, line := getCallerInfo()
	WarnLogger.WithFields(logrus.Fields{
		"operation": operation,
		"function":  funcName,
		"file":      fileName,
		"line":      line,
	}).Warn(message)
}

func Info(operation, message string) {
	funcName, fileName, line := getCallerInfo()
	InfoLogger.WithFields(logrus.Fields{
		"operation": operation,
		"function":  funcName,
		"file":      fileName,
		"line":      line,
	}).Info(message)
}

func Debug(operation, message string) {
	funcName, fileName, line := getCallerInfo()
	DebugLogger.WithFields(logrus.Fields{
		"operation": operation,
		"function":  funcName,
		"file":      fileName,
		"line":      line,
	}).Debug(message)
}
