package logger

// This package is used to log messages easily It is a wrapper around the logrus package.

import (
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

const (
	TraceLevel = logrus.TraceLevel
	DebugLevel = logrus.DebugLevel
	InfoLevel  = logrus.InfoLevel
	WarnLevel  = logrus.WarnLevel
	ErrorLevel = logrus.ErrorLevel
	FatalLevel = logrus.FatalLevel
	PanicLevel = logrus.PanicLevel
)

var Log *logrus.Logger

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.WarnLevel)
	Log = logrus.StandardLogger()
}

func ConfigureLogPath(logFile string) {
	folderPath := filepath.Dir(logFile)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		Log.Fatal("Log target folder does not exist, please contact your system administrator")
	}

	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		file, err := os.Create(logFile)
		if err != nil {
			Log.Fatal("Unable to create log file, please contact your system administrator")
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				Log.Fatal("Unable to close log file, please contact your system administrator")
			}
		}(file)
	}

	f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		Log.Fatal("Unable to open log file, please contact your system administrator")
	}
	logrus.SetOutput(f)
}

func EnableDebugMode() {
	SetLogLevel(logrus.DebugLevel.String())
}

func SetLogLevel(level string) {
	logrusLevel, err := logrus.ParseLevel(strings.ToLower(level))
	if err != nil {
		logrus.Warnf("Invalid log level: %s, defaulting to %s", level, logrus.WarnLevel.String())
	}

	logrus.SetLevel(logrusLevel)
}

func Trace(args ...interface{}) {
	Log.Trace(args...)
}

func Debug(args ...interface{}) {
	Log.Debug(args...)
}

func Info(args ...interface{}) {
	Log.Info(args...)
}

func Warn(args ...interface{}) {
	Log.Warn(args...)
}

func Error(args ...interface{}) {
	Log.Error(args...)
}

func Fatal(args ...interface{}) {
	Log.Fatal(args...)
}

func Panic(args ...interface{}) {
	Log.Panic(args...)
}

func Tracef(format string, args ...interface{}) {
	Log.Tracef(format, args...)
}

func Debugf(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	logrus.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	logrus.Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	logrus.Panicf(format, args...)
}
