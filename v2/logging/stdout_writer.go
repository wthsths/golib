package gl_logging

import (
	"fmt"
	"runtime"
	"time"

	"github.com/google/uuid"
	"github.com/teris-io/shortid"
)

var stdoutWriterShortIDGenerator *shortid.Shortid

type stdoutWriter struct {
	sessionID    string
	title        string
	disablePrint bool
}

func generateStdoutWriterSessionID() string {
	var err error
	if stdoutWriterShortIDGenerator == nil {
		stdoutWriterShortIDGenerator, err = shortid.New(1, shortid.DefaultABC, 2342)
		if err != nil {
			return uuid.NewString()
		}
	}
	generatedID, err := stdoutWriterShortIDGenerator.Generate()
	if err != nil {
		return uuid.NewString()
	}
	return generatedID
}

// NewStdoutWriter creates a stdoutWriter instance which implements Logger interface.
//
// Logging functions can be used to write data to stdout.
//
// NOTE DisablePrint() can be called to prevent outputs. (Can be desirable when writing tests etc...).
func NewStdoutWriter() *stdoutWriter {
	return &stdoutWriter{
		sessionID:    generateStdoutWriterSessionID(),
		disablePrint: false,
	}
}

// NewStdoutWriterWithCustomID creates a stdoutWriter instance with given input ID value.
//
// Input value will override the internally generated session ID value.
func NewStdoutWriterWithCustomID(id string) *stdoutWriter {
	return &stdoutWriter{
		sessionID: id,
	}
}

func (l *stdoutWriter) SetTitle(input string) {
	l.title = input
}

// DisablePrint prevents writing to stdout.
//
// This can be useful when using stdoutWriter as a dummy logger instance in tests.
func (l *stdoutWriter) DisablePrint() {
	l.disablePrint = true
}

// EnablePrint enables writing to stdout.
func (l *stdoutWriter) EnablePrint() {
	l.disablePrint = false
}

func (l *stdoutWriter) Infof(format string, args ...interface{}) {
	formatted := fmt.Sprintf(format, args...)
	log := l.getStructuredLog(logTypeInfo, formatted)
	fmt.Println(log)
}

func (l *stdoutWriter) Warnf(format string, args ...interface{}) {
	formatted := fmt.Sprintf(format, args...)
	log := l.getStructuredLog(logTypeWarn, formatted)
	fmt.Println(log)
}

func (l *stdoutWriter) Errorf(format string, args ...interface{}) {
	formatted := fmt.Sprintf(format, args...)
	log := l.getStructuredLog(logTypeError, formatted)
	fmt.Println(log)
}

func (l *stdoutWriter) Fatalf(format string, args ...interface{}) {
	formatted := fmt.Sprintf(format, args...)
	log := l.getStructuredLog(logTypeFatal, formatted)
	fmt.Println(log)
}

func (l *stdoutWriter) getStructuredLog(logType, content string) string {
	pc, filename, line, _ := runtime.Caller(2)

	var logSuffix string
	if logType != logTypeInfo {
		logSuffix = fmt.Sprintf("(%s%s:%d)", runtime.FuncForPC(pc).Name(), filename, line)
	}

	logTime := time.Now().UTC().Format("2006-01-02 15:04:05.000")
	return fmt.Sprintf("[%s][%s][%s][%s]: %s %s", logTime, l.sessionID, l.title, logType, content, logSuffix)
}
