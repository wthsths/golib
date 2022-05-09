package gl_logging

import (
	"fmt"
	"runtime"
	"time"
)

const (
	logTypeInfo  = "INFO"
	logTypeWarn  = "WARN"
	logTypeError = "ERROR"
	logTypeFatal = "FATAL"

	LogIgnored = "[LOG-IGNORED]"
)

var isInitialized = false
var writeToFileSystem = true

/* "glog" implementation is built upon: "https://github.com/birlesikodeme/glog" */

// Logger is an abstract representation of sessionLogger.
//
// Actual implementation is built upon a modified version of glog.
type Logger interface {
	// SetTitle adds title text which will be displayed in title bracket in logs.
	//
	// If SetTitle("title") was called prior to logging, output will be as follows:
	//
	// [2022-02-27 17:58:03.565][your-id][title][log-level]: fatal error!
	SetTitle(input string)

	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

// Init globally prepares logger for runtime.
// Must finish without errors to be able to write to stderr and file system.
//
// Creating a logger instance before the execution of Init will produce a panic.
func Init(name, dir string) error {
	err := logInit(name, dir)
	if err != nil {
		return err
	}

	isInitialized = true
	return nil
}

// SetFileSystemWrite can be used to enable or disable file system writes.
//
// Outputs to stdout will always occur.
func SetFileSystemWrite(active bool) {
	writeToFileSystem = active
}

type logger struct {
	sessionID string
	title     string
}

// NewLogger creates new logger instance with a internally generated unique ID.
//
// Created instance can be passed accross application layers for structed session logging.
//
// Example log output when Fatalf("fatal error!") is called:
//
// [2022-02-27 17:58:03.565][your-id][your-title][FATAL]: fatal error!
func NewLogger(title, sessionID string) *logger {
	return &logger{
		title:     title,
		sessionID: sessionID,
	}
}

// NewLoggerWithID creates a logger instance with given input ID value.
func NewLoggerWithID(id string) *logger {
	return &logger{
		sessionID: id,
	}
}

func (l *logger) SetTitle(input string) {
	l.title = input
}

func (l *logger) Infof(format string, args ...interface{}) {
	if writeToFileSystem && !isInitialized {
		panic("Logger is not initialized yet. logging.Init() must be executed first to write logs to file system.")
	}

	formatted := fmt.Sprintf(format, args...)
	log := l.getStructuredLog(logTypeInfo, formatted)

	if writeToFileSystem {
		writeLn(log)
	}
	fmt.Println(log)

}

func (l *logger) Warnf(format string, args ...interface{}) {
	if writeToFileSystem && !isInitialized {
		panic("Logger is not initialized yet. logging.Init() must be executed first to write logs to file system.")
	}

	formatted := fmt.Sprintf(format, args...)
	log := l.getStructuredLog(logTypeWarn, formatted)

	if writeToFileSystem {
		writeLn(log)
	}
	fmt.Println(log)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	if writeToFileSystem && !isInitialized {
		panic("Logger is not initialized yet. logging.Init() must be executed first to write logs to file system.")
	}

	formatted := fmt.Sprintf(format, args...)
	log := l.getStructuredLog(logTypeError, formatted)

	if writeToFileSystem {
		writeLn(log)
	}
	fmt.Println(log)
}

func (l *logger) Fatalf(format string, args ...interface{}) {
	if writeToFileSystem && !isInitialized {
		panic("Logger is not initialized yet. logging.Init() must be executed first to write logs to file system.")
	}

	formatted := fmt.Sprintf(format, args...)
	log := l.getStructuredLog(logTypeFatal, formatted)

	if writeToFileSystem {
		writeLn(log)
	}
	fmt.Println(log)
}

func (l *logger) getStructuredLog(logType, content string) string {
	pc, filename, line, _ := runtime.Caller(2)

	var logSuffix string
	if logType != logTypeInfo {
		logSuffix = fmt.Sprintf("(%s%s:%d)", runtime.FuncForPC(pc).Name(), filename, line)
	}

	logTime := time.Now().UTC().Format("2006-01-02 15:04:05.000")
	return fmt.Sprintf("[%s][%s][%s][%s]: %s %s", logTime, l.sessionID, l.title, logType, content, logSuffix)
}
