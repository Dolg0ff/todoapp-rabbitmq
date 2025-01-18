package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	log *logrus.Logger
}

func NewLogger(level string) *Logger {
	l := logrus.New()

	_, err := os.Getwd()
	if err != nil {
		fmt.Printf("Failed to get current directory: %v\n", err)
	}

	logsDir := "logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		fmt.Printf("Failed to create logs directory: %v\n", err)
	}

	logPath := filepath.Join(logsDir, "app.log")
	fmt.Printf("Log file path: %s\n", logPath)

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
		l.SetOutput(os.Stdout)
	} else {
		fmt.Printf("Successfully opened log file\n")
		l.SetOutput(io.MultiWriter(os.Stdout, file))
	}

	l.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	switch level {
	case "debug":
		l.SetLevel(logrus.DebugLevel)
	case "info":
		l.SetLevel(logrus.InfoLevel)
	case "warn":
		l.SetLevel(logrus.WarnLevel)
	case "error":
		l.SetLevel(logrus.ErrorLevel)
	default:
		l.SetLevel(logrus.InfoLevel)
	}

	return &Logger{
		log: l,
	}
}

func (l *Logger) Info(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.log.WithFields(convertArgsToFields(args...)).Info(msg)
	} else {
		l.log.Info(msg)
	}
}

func (l *Logger) Error(msg string, err error, args ...interface{}) {
	fields := logrus.Fields{}
	if err != nil {
		fields["error"] = err.Error()
	}

	if len(args) > 0 {
		for k, v := range convertArgsToFields(args...) {
			fields[k] = v
		}
	}

	l.log.WithFields(fields).Error(msg)
}

func (l *Logger) Fatal(msg string, err error, args ...interface{}) {
	fields := logrus.Fields{}
	if err != nil {
		fields["error"] = err.Error()
	}

	if len(args) > 0 {
		for k, v := range convertArgsToFields(args...) {
			fields[k] = v
		}
	}

	l.log.WithFields(fields).Fatal(msg)
}

// Вспомогательная функция для конвертации аргументов в поля
func convertArgsToFields(args ...interface{}) logrus.Fields {
	fields := logrus.Fields{}

	// Конвертируем аргументы в пары ключ-значение
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			fields[args[i].(string)] = args[i+1]
		}
	}

	return fields
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.log.WithFields(convertArgsToFields(args...)).Debug(msg)
	} else {
		l.log.Debug(msg)
	}
}

func (l *Logger) Warn(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.log.WithFields(convertArgsToFields(args...)).Warn(msg)
	} else {
		l.log.Warn(msg)
	}
}
