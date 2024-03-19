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
	LogFile      string
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
	LogFile = logFile
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

// Reapply log filepath to logrus logger to avoid log rotation issues
// Truthfully I cant be bothered making windows & unix compatible file descriptor checks as Log.Out is not updated, so this is will do for now
func logRotateHandler() {
	if LogFile == "" {
		return
	}
	file, _ := os.OpenFile(LogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	Log.SetOutput(file)
}

func logWithLevel(level logrus.Level, args ...interface{}) {
	logRotateHandler()
	switch level {
	case TraceLevel:
		Log.Trace(args...)
	case DebugLevel:
		Log.Debug(args...)
	case InfoLevel:
		Log.Info(args...)
	case WarnLevel:
		Log.Warn(args...)
	case ErrorLevel:
		Log.Error(args...)
	case FatalLevel:
		Log.Fatal(args...)
	case PanicLevel:
		Log.Panic(args...)
	}
}

func logFormattedWithLevel(level logrus.Level, format string, args ...interface{}) {
	logRotateHandler()
	switch level {
	case TraceLevel:
		Log.Tracef(format, args...)
	case DebugLevel:
		Log.Debugf(format, args...)
	case InfoLevel:
		Log.Infof(format, args...)
	case WarnLevel:
		Log.Warnf(format, args...)
	case ErrorLevel:
		Log.Errorf(format, args...)
	case FatalLevel:
		Log.Fatalf(format, args...)
	case PanicLevel:
		Log.Panicf(format, args...)
	}
}

func Trace(args ...interface{}) {
	logWithLevel(TraceLevel, args...)
}

func Debug(args ...interface{}) {
	logWithLevel(DebugLevel, args...)
}

func Info(args ...interface{}) {
	logWithLevel(InfoLevel, args...)
}

func Warn(args ...interface{}) {
	logWithLevel(WarnLevel, args...)
}

func Error(args ...interface{}) {
	logWithLevel(ErrorLevel, args...)
}

func Fatal(args ...interface{}) {
	logWithLevel(FatalLevel, args...)
}

func Panic(args ...interface{}) {
	logWithLevel(PanicLevel, args...)
}

func Tracef(format string, args ...interface{}) {
	logFormattedWithLevel(TraceLevel, format, args...)
}

func Debugf(format string, args ...interface{}) {
	logFormattedWithLevel(DebugLevel, format, args...)
}

func Infof(format string, args ...interface{}) {
	logFormattedWithLevel(InfoLevel, format, args...)
}

func Warnf(format string, args ...interface{}) {
	logFormattedWithLevel(WarnLevel, format, args...)
}

func Errorf(format string, args ...interface{}) {
	logFormattedWithLevel(ErrorLevel, format, args...)
}

func Fatalf(format string, args ...interface{}) {
	logFormattedWithLevel(FatalLevel, format, args...)
}

func Panicf(format string, args ...interface{}) {
	logFormattedWithLevel(PanicLevel, format, args...)
}
