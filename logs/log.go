package logs

import (
	"context"

	"github.com/sirupsen/logrus"
)

var stdLogger = &logrusLogger{logger, nil, ""}

func WithPackage(pkg string) Logger {
	return stdLogger.WithPackage(pkg)
}

func WithContext(ctx context.Context) Logger {
	return stdLogger.WithContext(ctx)
}

func EnableDebug() {
	logger.SetLevel(logrus.DebugLevel)
}

func DisableDebug() {
	logger.SetLevel(logrus.InfoLevel)
}
