package logs

import (
	"context"

	"github.com/sirupsen/logrus"
)

var stdLogger = &logrusLogger{logrus.New(), nil, ""}

func WithPackage(pkg string) Logger {
	return stdLogger.WithPackage(pkg)
}

func WithContext(ctx context.Context) Logger {
	return stdLogger.WithContext(ctx)
}
