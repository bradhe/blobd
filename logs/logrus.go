package logs

import (
	"context"

	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

type logrusLogger struct {
	base *logrus.Logger
	ctx  context.Context
	pkg  string
}

func (l *logrusLogger) clone() *logrusLogger {
	return &logrusLogger{l.base, l.ctx, l.pkg}
}

func (l *logrusLogger) Entry() Entry {
	baseEntry := logrus.NewEntry(l.base)

	if l.pkg != "" {
		baseEntry = baseEntry.WithField("package", l.pkg)
	}

	return &logrusEntry{baseEntry}
}

func (l *logrusLogger) WithFields(fields map[string]interface{}) Entry {
	return l.Entry().WithFields(fields)
}

func (l *logrusLogger) WithField(name string, value interface{}) Entry {
	return l.WithFields(map[string]interface{}{name: value})
}

func (l *logrusLogger) WithError(err error) Entry {
	return l.Entry().WithFields(logrus.Fields{"error": err})
}

func (l *logrusLogger) WithPackage(pkg string) Logger {
	cp := l.clone()
	cp.pkg = pkg
	return cp
}

func (l *logrusLogger) WithContext(ctx context.Context) Logger {
	cp := l.clone()
	cp.ctx = ctx
	return cp
}

func (l *logrusLogger) Print(str string) {
	l.Entry().Print(str)
}

func (l *logrusLogger) Printf(format string, args ...interface{}) {
	l.Entry().Printf(format, args...)
}

func (l *logrusLogger) Debug(str string) {
	l.Entry().Debug(str)
}

func (l *logrusLogger) Debugf(format string, args ...interface{}) {
	l.Entry().Debugf(format, args...)
}

func (l *logrusLogger) Info(str string) {
	l.Entry().Info(str)
}

func (l *logrusLogger) Infof(format string, args ...interface{}) {
	l.Entry().Infof(format, args...)
}

func (l *logrusLogger) Error(str string) {
	l.Entry().Error(str)
}

func (l *logrusLogger) Errorf(format string, args ...interface{}) {
	l.Entry().Errorf(format, args...)
}

func (l *logrusLogger) Warn(str string) {
	l.Entry().Warn(str)
}

func (l *logrusLogger) Warnf(format string, args ...interface{}) {
	l.Entry().Warnf(format, args...)
}

func (l *logrusLogger) Fatal(str string) {
	l.Entry().Warn(str)
}

func (l *logrusLogger) Fatalf(format string, args ...interface{}) {
	l.Entry().Fatalf(format, args...)
}

type logrusEntry struct {
	base *logrus.Entry
}

func (e *logrusEntry) WithFields(fields map[string]interface{}) Entry {
	newEntry := e.base.WithFields(logrus.Fields(fields))
	return &logrusEntry{newEntry}
}

func (e *logrusEntry) WithField(name string, value interface{}) Entry {
	return e.WithFields(map[string]interface{}{name: value})
}

func (e *logrusEntry) Print(str string) {
	e.base.Print(str)
}

func (e *logrusEntry) Printf(format string, args ...interface{}) {
	e.base.Printf(format, args...)
}

func (e *logrusEntry) Debug(str string) {
	e.base.Debug(str)
}

func (e *logrusEntry) Debugf(format string, args ...interface{}) {
	e.base.Debugf(format, args...)
}

func (e *logrusEntry) Info(str string) {
	e.base.Info(str)
}

func (e *logrusEntry) Infof(format string, args ...interface{}) {
	e.base.Infof(format, args...)
}

func (e *logrusEntry) Error(str string) {
	e.base.Error(str)
}

func (e *logrusEntry) Errorf(format string, args ...interface{}) {
	e.base.Errorf(format, args...)
}

func (e *logrusEntry) Warn(str string) {
	e.base.Warn(str)
}

func (e *logrusEntry) Warnf(format string, args ...interface{}) {
	e.base.Warnf(format, args...)
}

func (e *logrusEntry) Fatal(str string) {
	e.base.Warn(str)
}

func (e *logrusEntry) Fatalf(format string, args ...interface{}) {
	e.base.Fatalf(format, args...)
}
