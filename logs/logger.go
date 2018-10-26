package logs

import (
	"context"
)

type Entry interface {
	WithFields(map[string]interface{}) Entry
	WithField(string, interface{}) Entry

	Print(string)
	Printf(string, ...interface{})

	Debug(string)
	Debugf(string, ...interface{})

	Info(string)
	Infof(string, ...interface{})

	Error(string)
	Errorf(string, ...interface{})

	Warn(string)
	Warnf(string, ...interface{})

	Fatal(string)
	Fatalf(string, ...interface{})
}

type Logger interface {
	Entry() Entry
	WithFields(map[string]interface{}) Entry
	WithField(string, interface{}) Entry
	WithError(error) Entry

	WithContext(context.Context) Logger
	WithPackage(string) Logger

	Print(string)
	Printf(string, ...interface{})

	Debug(string)
	Debugf(string, ...interface{})

	Info(string)
	Infof(string, ...interface{})

	Error(string)
	Errorf(string, ...interface{})

	Warn(string)
	Warnf(string, ...interface{})

	Fatal(string)
	Fatalf(string, ...interface{})
}
