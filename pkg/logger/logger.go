package logger

type Logger struct {
	level string
}

func NewLogger(level string) *Logger {
	return &Logger{
		level: level,
	}
}

func (l *Logger) Info(msg string, args ...interface{}) {
	// Реализация логирования
}

func (l *Logger) Error(msg string, err error, args ...interface{}) {
	// Реализация логирования ошибок
}

func (l *Logger) Fatal(msg string, err error, args ...interface{}) {
	// Логирование и завершение программы
}
