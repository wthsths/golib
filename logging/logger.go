package logging

import (
	"fmt"
	"runtime"
	"time"

	"github.com/teris-io/shortid"
)

const (
	logTypeInfo  = "INFO"
	logTypeWarn  = "WARN"
	logTypeError = "ERROR"
	logTypeFatal = "FATAL"

	LogIgnored = "[LOG-IGNORED]"
)

var isInitialized = false
var shortIDGenerator *shortid.Shortid

/* "glog" implementation is built upon: "https://github.com/birlesikodeme/glog" */

// Logger is an abstract representation of sessionLogger.
//
// Actual implementation is built upon a modified version of glog.
type Logger interface {
	// SetTitle adds title text which will be displayed in title bracket in logs.
	//
	// If SetTitle("title") was called prior to logging, output will be as follows:
	//
	// [2022-02-27 17:58:03.565][KFTGcuiQ9p][title][FATAL]: fatal error!
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
	err := glogInit(name, dir)
	if err != nil {
		return err
	}

	shortIDGenerator, err = shortid.New(1, shortid.DefaultABC, 2342)
	if err != nil {
		return fmt.Errorf("unable to initialize short ID generator: %s", err.Error())
	}

	isInitialized = true
	return nil
}

type sessionLogger struct {
	sessionID string
	title     string
}

// NewSessionLogger creates new logger instance with a internally generated unique ID.
//
// Created instance can be passed accross application layers for structed session logging.
//
// Example log output when Fatalf("fatal error!") is called:
//
// [2022-02-27 17:58:03.565][KFTGcuiQ9p][][FATAL]: fatal error!
//
// If SetTitle("title") was called prior to logging, output will be as follows:
//
// [2022-02-27 17:58:03.565][KFTGcuiQ9p][title][FATAL]: fatal error!
func NewSessionLogger() *sessionLogger {
	if !isInitialized {
		panic("Logger is not initialized yet. logging.Init() must be executed first.")
	}

	generatedID, _ := shortIDGenerator.Generate()
	return &sessionLogger{
		sessionID: generatedID,
	}
}

// NewSessionLoggerWithCustomID creates a logger instance with given input ID value.
//
// Input value will override the internally generated session ID value.
func NewSessionLoggerWithCustomID(id string) *sessionLogger {
	if !isInitialized {
		panic("Logger is not initialized yet. logging.Init() must be executed first.")
	}

	return &sessionLogger{
		sessionID: id,
	}
}

func (l *sessionLogger) SetTitle(input string) {
	l.title = input
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
	return fmt.Sprintf("[%s][%s][%s][%s]: %s %s", logTime, l.sessionID, l.title, logType, content, logSuffix)
}
