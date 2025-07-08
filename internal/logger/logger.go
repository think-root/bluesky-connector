package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func Init(level string) {
	Logger = logrus.New()
	Logger.SetOutput(os.Stdout)
	Logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logLevel = logrus.InfoLevel
		Logger.Warnf("Invalid log level '%s', defaulting to info", level)
	}
	Logger.SetLevel(logLevel)

	Logger.Info("Logger initialized")
}

func Info(args ...interface{}) {
	Logger.Info(args...)
}

func Infof(format string, args ...interface{}) {
	Logger.Infof(format, args...)
}

func Warn(args ...interface{}) {
	Logger.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	Logger.Warnf(format, args...)
}

func Error(args ...interface{}) {
	Logger.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	Logger.Errorf(format, args...)
}

func Debug(args ...interface{}) {
	Logger.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	Logger.Debugf(format, args...)
}

func Fatal(args ...interface{}) {
	Logger.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	Logger.Fatalf(format, args...)
}