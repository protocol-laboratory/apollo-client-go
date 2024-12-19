package apollo

import (
	"fmt"
	"golang.org/x/exp/slog"
)

type Logger interface {
	Info(format string, args ...interface{})

	Error(format string, args ...interface{})

	Warn(format string, args ...interface{})

	Infof(format string, args ...interface{})

	Errorf(format string, args ...interface{})

	Warnf(format string, args ...interface{})
}

type defaultLogger struct {
	Logger *slog.Logger
}

var log Logger = &defaultLogger{
	Logger: slog.Default(),
}

func SetLogger(logger Logger) {
	log = logger
}

func (d *defaultLogger) Info(format string, args ...interface{}) {
	d.Logger.Info(format, args...)
}

func (d *defaultLogger) Error(format string, args ...interface{}) {
	d.Logger.Error(format, args...)
}

func (d *defaultLogger) Warn(format string, args ...interface{}) {
	d.Logger.Warn(format, args...)
}

func (d *defaultLogger) Infof(format string, args ...interface{}) {
	d.Logger.Info(fmt.Sprintf(format, args...))
}

func (d *defaultLogger) Errorf(format string, args ...interface{}) {
	d.Logger.Error(fmt.Sprintf(format, args...))
}

func (d *defaultLogger) Warnf(format string, args ...interface{}) {
	d.Logger.Warn(fmt.Sprintf(format, args...))
}
