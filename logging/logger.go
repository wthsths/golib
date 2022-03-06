package logging

import (
	"fmt"
	"runtime"
	"time"

	"github.com/google/uuid"
)

const (
	logTypeInfo  = "INFO"
	logTypeWarn  = "WARN"
	logTypeError = "ERROR"
	logTypeFatal = "FATAL"

	LogIgnored = "[LOG-IGNORED]"
)

var isInitialized = false

// "glog" implementation is built upon: "https://github.com/birlesikodeme/glog"

// Logger is an abstract representation of sessionLogger.
//
// Actual implementation is built upon a modified version of glog.
type Logger interface {
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

// Init globally prepares logger for runtime.
// Must pass without errors to be able to write to stderr and file system.
func Init(name, dir string) error {
	err := glogInit(name, dir)
	if err != nil {
		return err
	}

	isInitialized = true
	return nil
}

type sessionLogger struct {
	sessionID string
}

// NewSessionLogger creates new logger instance with a unique ID.
// Created instance can be passed accross application layers for structed session logging.
//
// Example log output when Fatalf("fatal error!") is called.
//
// [2022-02-27 17:58:03.565][7cc2d15b-6069-495c-9f84-89c4ba4c5566][FATAL]: fatal error!
func NewSessionLogger() *sessionLogger {
	if !isInitialized {
		panic("Logger is not initialized yet. logging.Init() must be executed first.")
	}

	return &sessionLogger{
		sessionID: uuid.NewString(),
	}
}

// NewSessionLoggerWithCustomID creates a logger instance with given input ID value.
// Input value will override the interbally generated session ID value.
func NewSessionLoggerWithCustomID(id string) *sessionLogger {
	if !isInitialized {
		panic("Logger is not initialized yet. logging.Init() must be executed first.")
	}

	return &sessionLogger{
		sessionID: id,
	}
}

func (l *sessionLogger) Infof(format string, args ...interface{}) {
	formatted := fmt.Sprintf(format, args...)
	log := l.getStructuredLog(logTypeInfo, formatted)
	infoln(log)
}

func (l *sessionLogger) Warnf(format string, args ...interface{}) {
	formatted := fmt.Sprintf(format, args...)
	log := l.getStructuredLog(logTypeWarn, formatted)
	warningln(log)
}

func (l *sessionLogger) Errorf(format string, args ...interface{}) {
	formatted := fmt.Sprintf(format, args...)
	log := l.getStructuredLog(logTypeError, formatted)
	errorln(log)
}

func (l *sessionLogger) Fatalf(format string, args ...interface{}) {
	formatted := fmt.Sprintf(format, args...)
	log := l.getStructuredLog(logTypeFatal, formatted)
	fatalln(log)
}

func (l *sessionLogger) getStructuredLog(logType, content string) string {
	pc, filename, line, _ := runtime.Caller(2)

	var logSuffix string
	if logType != logTypeInfo {
		logSuffix = fmt.Sprintf("(%s%s:%d)", runtime.FuncForPC(pc).Name(), filename, line)
	}

	logTime := time.Now().UTC().Format("2006-01-02 15:04:05.000")
	return fmt.Sprintf("[%s][%s][%s]: %s %s", logTime, l.sessionID, logType, content, logSuffix)
}
