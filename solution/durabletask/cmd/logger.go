package main

import (
	"fmt"
	"log/slog"
)

type logger struct {
	slog.Logger
}

func newLogger() logger {
	return logger{
		Logger: *slog.Default(),
	}
}

func (l logger) Debug(v ...any) {
	l.Logger.Debug(fmt.Sprint(v...))
}

func (l logger) Debugf(format string, v ...any) {
	l.Logger.Debug(fmt.Sprintf(format, v...))
}

func (l logger) Info(v ...any) {
	l.Logger.Info(fmt.Sprint(v...))
}

func (l logger) Infof(format string, v ...any) {
	l.Logger.Info(fmt.Sprintf(format, v...))
}

func (l logger) Warn(v ...any) {
	l.Logger.Warn(fmt.Sprint(v...))
}

func (l logger) Warnf(format string, v ...any) {
	l.Logger.Warn(fmt.Sprintf(format, v...))
}

func (l logger) Error(v ...any) {
	l.Logger.Error(fmt.Sprint(v...))
}

func (l logger) Errorf(format string, v ...any) {
	l.Logger.Error(fmt.Sprintf(format, v...))
}
