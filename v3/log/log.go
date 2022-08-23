package log

import (
	gelfTCP "github.com/payports/golib/v3/log/gelf"
	gelf "github.com/snovichkov/zap-gelf"
	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"sync"
)

type Logger interface {
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Debugw(msg string, keysAndValues ...interface{})
	DPanic(args ...interface{})
	DPanicf(template string, args ...interface{})
	DPanicw(msg string, keysAndValues ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Panic(args ...interface{})
	Panicf(template string, args ...interface{})
	Panicw(msg string, keysAndValues ...interface{})
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
}

var logger Logger = defaultSugarLogger()

type LoggerOption func() zapcore.Core

var once sync.Once

func GetLogger() Logger {
	return logger
}

func InitLogger(loggerName string, options ...LoggerOption) {
	once.Do(func() {
		var logOptions []zapcore.Core

		for _, opt := range options {
			logOptions = append(logOptions, opt())
		}

		core := zapcore.NewTee(logOptions...)
		l := zap.New(core, zap.AddCaller())

		// set logger
		logger = l.Sugar().Named(loggerName)
	})
}

func WithIO(w io.Writer, logLevel int, environmentType, logEncoding string) LoggerOption {
	return func() zapcore.Core {
		var (
			encoderConfig zapcore.EncoderConfig
			encoder       zapcore.Encoder
		)

		if environmentType == "production" {
			encoderConfig = zap.NewProductionEncoderConfig()
		} else {
			encoderConfig = zap.NewDevelopmentEncoderConfig()
		}

		if logEncoding == "console" {
			// adds '\t' to eof
			encoder = zapcore.NewConsoleEncoder(encoderConfig)
		} else {
			encoder = zapcore.NewJSONEncoder(encoderConfig)
		}

		stdoutLogger := zapcore.NewCore(encoder, zapcore.AddSync(w), zap.NewAtomicLevelAt(zapcore.Level(logLevel)))

		return stdoutLogger
	}
}

func WithGraylogViaUDP(logLevel int, addr string) LoggerOption {
	return func() zapcore.Core {
		// First go to localhost:9000 on browser
		// login : admin, pass : admin
		// System > Inputs > Select input > GELF UDP > Launch New input
		// note : make docker volume prune
		grayLogger, err := gelf.NewCore(gelf.Addr(addr), gelf.Level(zapcore.Level(logLevel)))
		if err != nil {
			panic(err)
		}

		return grayLogger
	}
}

func WithGraylogViaTCP(logLevel int, addr string) LoggerOption {
	return func() zapcore.Core {
		grayLogger, err := gelfTCP.NewTcpCore(gelfTCP.Addr(addr), gelfTCP.Level(zapcore.Level(logLevel)))
		if err != nil {
			panic(err)
		}

		return grayLogger
	}
}

func WithElasticCompatible(logLevel int) LoggerOption {
	return func() zapcore.Core {
		elastic := ecszap.NewCore(ecszap.NewDefaultEncoderConfig(), zapcore.AddSync(os.Stdout), zap.NewAtomicLevelAt(zapcore.Level(logLevel)))

		return elastic
	}
}

func defaultSugarLogger() Logger {
	return zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentConfig().EncoderConfig),
		zapcore.AddSync(os.Stderr),
		zap.DebugLevel,
	)).Sugar()
}
