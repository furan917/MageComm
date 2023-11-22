package logger

// This package is used to Log messages easily It is a wrapper around the logrus package.

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

var (
	Log          *logrus.Logger
	debugFlagSet bool
)

func init() {
	Log = logrus.StandardLogger()
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	Log.SetLevel(logrus.WarnLevel)
}

func ConfigureLogPath(logFile string) {
	folderPath := filepath.Dir(logFile)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		Log.Errorf("Log target folder does not exist, please contact your system administrator")
	}

	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		if err := createLogFile(logFile); err != nil {
			Log.Errorf("Unable to create Log file, please contact your system administrator")
		}
	}

	file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		Log.SetOutput(os.Stdout)
		Log.Errorf("Unable to open Log file, printing to stdout, please contact your system administrator")
		return
	}

	//Log.SetFormatter(&logrus.JSONFormatter{})
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	Log.SetOutput(file)
	Log.Infof("Logging to file: %s", file.Name())
}

func createLogFile(logFile string) error {
	file, err := os.Create(logFile)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			Log.Fatal("Unable to close Log file, please contact your system administrator")
		}
	}(file)
	return nil
}

func EnableDebugMode() {
	SetLogLevel(logrus.DebugLevel.String())
	debugFlagSet = true
}

func SetLogLevel(level string) {
	if debugFlagSet {
		Log.Info("Debug mode enabled by flag, ignoring log level configuration")
		return
	}

	logrusLevel, err := logrus.ParseLevel(strings.ToLower(level))
	if err != nil {
		Log.Warnf("Invalid Log level: %s, defaulting to %s", level, logrus.WarnLevel.String())
		logrusLevel = logrus.WarnLevel
	}

	Log.SetLevel(logrusLevel)
	Log.Infof("Log level set to %s", logrusLevel)
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
	Log.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	Log.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	Log.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	Log.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	Log.Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	Log.Panicf(format, args...)
}
