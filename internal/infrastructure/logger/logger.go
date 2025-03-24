package logger

import (
	"log"
	"os"

	"github.com/vadimkiryanov/tg-bot-lebowski/internal/domain/ports"
)

type Logger struct {
	logger *log.Logger
	fields map[string]interface{}
}

func NewLogger() ports.Logger {
	return &Logger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
		fields: make(map[string]interface{}),
	}
}

func (l *Logger) Info(args ...interface{}) {
	l.logger.Println(append([]interface{}{"[INFO]"}, args...)...)
}

func (l *Logger) Error(args ...interface{}) {
	l.logger.Println(append([]interface{}{"[ERROR]"}, args...)...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.logger.Fatal(append([]interface{}{"[FATAL]"}, args...)...)
}

func (l *Logger) WithFields(fields map[string]interface{}) ports.Logger {
	return &Logger{
		logger: l.logger,
		fields: fields,
	}
}
