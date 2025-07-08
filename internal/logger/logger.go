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

func Info(args ...any) {
	Logger.Info(args...)
}

func Infof(format string, args ...any) {
	Logger.Infof(format, args...)
}

func Warn(args ...any) {
	Logger.Warn(args...)
}

func Warnf(format string, args ...any) {
	Logger.Warnf(format, args...)
}

func Error(args ...any) {
	Logger.Error(args...)
}

func Errorf(format string, args ...any) {
	Logger.Errorf(format, args...)
}

func Debug(args ...any) {
	Logger.Debug(args...)
}

func Debugf(format string, args ...any) {
	Logger.Debugf(format, args...)
}

func Fatal(args ...any) {
	Logger.Fatal(args...)
}

func Fatalf(format string, args ...any) {
	Logger.Fatalf(format, args...)
}