package logging

import "example.com/api/config"

type ILogger interface {
	Init()

	Debug(cat Category, sub SubCategory, msg string, extra map[ExtraKey]any)
	Debugf(template string, args ...any)

	Info(cat Category, sub SubCategory, msg string, extra map[ExtraKey]any)
	Infof(template string, args ...any)

	Warn(cat Category, sub SubCategory, msg string, extra map[ExtraKey]any)
	Warnf(template string, args ...any)

	Error(cat Category, sub SubCategory, msg string, extra map[ExtraKey]any)
	Errorf(template string, args ...any)

	Fatal(cat Category, sub SubCategory, msg string, extra map[ExtraKey]any)
	Fatalf(template string, args ...any)
}

func NewLogger(cfg *config.Config) ILogger {
	return newZeroLogger(cfg)
}