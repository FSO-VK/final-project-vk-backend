package logcon

import (
	"context"

	"github.com/sirupsen/logrus"
)

// key is a context key
type key string

const (
	loggerKey key = "logger"
)

// WithContext adds a logger to the context
func WithContext(ctx context.Context, logger *logrus.Entry) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// FromContext returns a logger from the context
func FromContext(ctx context.Context) (*logrus.Entry, bool) {
	l, ok := ctx.Value(loggerKey).(*logrus.Entry)
	return l, ok
}
