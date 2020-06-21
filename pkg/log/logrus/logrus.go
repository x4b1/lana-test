package logrus

import (
	"context"

	"github.com/xabi93/lana-test/internal/server"
	"github.com/xabi93/lana-test/pkg/log"

	"github.com/sirupsen/logrus"
)

//New creates a new instance of logrus logger
func New() log.Logger {
	return &logger{logrus.New()}
}

type logger struct {
	*logrus.Logger
}

func (l logger) Error(ctx context.Context, err error) {
	l.WithDefaultFields(ctx).Error(err)
}

func (l logger) Fatal(ctx context.Context, err error) {
	l.WithDefaultFields(ctx).Fatal(err)
}

func (l logger) Info(ctx context.Context, msg string, args ...interface{}) {
	l.Infof(msg, args...)
}

func (l logger) WithDefaultFields(ctx context.Context) *logrus.Entry {
	endpoint, _ := server.Endpoint(ctx)

	f := logrus.Fields{"endpoint": endpoint}

	return l.WithFields(f)
}
