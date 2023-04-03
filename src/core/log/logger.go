package log

import (
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

type Logger interface {
	With(ctx context.Context, args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Error(args ...interface{})
}

type logger struct {
	*zap.SugaredLogger
}

func (l *logger) With(ctx context.Context, args ...interface{}) {

}

func New() Logger {
	l, _ := zap.NewProduction()
	return NewWithZap(l)
}

func NewWithZap(l *zap.Logger) Logger {
	return &logger{l.Sugar()}
}
