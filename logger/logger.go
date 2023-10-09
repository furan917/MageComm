package logger

// This package is used to log messages easily It is a wrapper around the logrus package.

import (
	"fmt"
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

	DefaultPrintoutLevel = ErrorLevel
)

var Log *logrus.Logger
var printoutLevel = DefaultPrintoutLevel

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
		Log.Errorf("Log target folder does not exist, please contact your system administrator")
	}

	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		if err := createLogFile(logFile); err != nil {
			Log.Errorf("Unable to create log file, please contact your system administrator")
		}
	}

	file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		Log.Errorf("Unable to open log file, printing to stdout, please contact your system administrator")
		logrus.SetOutput(os.Stdout)
		return
	}

	logrus.SetOutput(file)
}

func createLogFile(logFile string) error {
	file, err := os.Create(logFile)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			Log.Fatal("Unable to close log file, please contact your system administrator")
		}
	}(file)
	return nil
}

func EnableDebugMode() {
	SetLogLevel(logrus.DebugLevel.String())
}

func SetLogLevel(level string) {
	logrusLevel, err := logrus.ParseLevel(strings.ToLower(level))
	if err != nil {
		logrus.Warnf("Invalid log level: %s, defaulting to %s", level, logrus.WarnLevel.String())
		logrusLevel = logrus.WarnLevel
	}

	logrus.SetLevel(logrusLevel)
}

// future improvement: add a config option for this
func SetPrintoutLevel(level string) {
	logrusLevel, err := logrus.ParseLevel(strings.ToLower(level))
	if err != nil {
		logrus.Warnf("Invalid log level given for CLI print out level definement: %s, defaulting to %s", level, logrus.ErrorLevel.String())
		logrusLevel = logrus.ErrorLevel
	}
	printoutLevel = logrusLevel
}

func Trace(args ...interface{}) {
	printLog(TraceLevel, args...)
	Log.Trace(args...)
}

func Debug(args ...interface{}) {
	printLog(DebugLevel, args...)
	Log.Debug(args...)
}

func Info(args ...interface{}) {
	printLog(InfoLevel, args...)
	Log.Info(args...)
}

func Warn(args ...interface{}) {
	printLog(WarnLevel, args...)
	Log.Warn(args...)
}

func Error(args ...interface{}) {
	printLog(ErrorLevel, args...)
	Log.Error(args...)
}

func Fatal(args ...interface{}) {
	printLog(FatalLevel, args...)
	Log.Fatal(args...)
}

func Panic(args ...interface{}) {
	printLog(PanicLevel, args...)
	Log.Panic(args...)
}

func Tracef(format string, args ...interface{}) {
	printLogF(TraceLevel, format, args...)
	Log.Tracef(format, args...)
}

func Debugf(format string, args ...interface{}) {
	printLogF(DebugLevel, format, args...)
	logrus.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	printLogF(InfoLevel, format, args...)
	logrus.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	printLogF(WarnLevel, format, args...)
	logrus.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	printLogF(ErrorLevel, format, args...)
	logrus.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	printLogF(FatalLevel, format, args...)
	logrus.Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	printLogF(PanicLevel, format, args...)
	logrus.Panicf(format, args...)
}

func printLogF(level logrus.Level, format string, args ...interface{}) {
	if int(level) <= int(printoutLevel) {
		fmt.Printf(format, args...)
	}
}

func printLog(level logrus.Level, args ...interface{}) {
	if int(level) <= int(printoutLevel) {
		fmt.Print(args...)
	}
}
