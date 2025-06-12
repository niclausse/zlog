package zlog

import (
	"context"
	"go.uber.org/zap"
	"strings"
)

func GetSugaredLogger() *zap.SugaredLogger {
	if sugaredLogger == nil {
		sugaredLogger = GetZapLogger().Sugar()
	}

	return sugaredLogger
}

func Debug(args ...interface{}) {
	GetSugaredLogger().Debug(args)
}

func Debugf(format string, args ...interface{}) {
	GetSugaredLogger().Debugf(format, args...)
}

func Info(args ...interface{}) {
	GetSugaredLogger().Info(args)
}

func Infof(format string, args ...interface{}) {
	GetSugaredLogger().Infof(format, args...)
}

func Warn(args ...interface{}) {
	GetSugaredLogger().Warn(args)
}

func Warnf(format string, args ...interface{}) {
	GetSugaredLogger().Warnf(format, args...)
}

func Error(args ...interface{}) {
	GetSugaredLogger().Error(args)
}

func Errorf(format string, args ...interface{}) {
	GetSugaredLogger().Errorf(format, args...)
}

func Panic(args ...interface{}) {
	GetSugaredLogger().Panic(args)
}

func Panicf(format string, args ...interface{}) {
	GetSugaredLogger().Panicf(format, args...)
}

func Fatal(args ...interface{}) {
	GetSugaredLogger().Fatal(args)
}

func Fatalf(format string, args ...interface{}) {
	GetSugaredLogger().Fatalf(format, args...)
}

func WithContext(ctx context.Context) *zap.SugaredLogger {
	args := func() (a []interface{}) {
		uri, _ := ctx.Value(ContextKeyURI).(string)
		requestId, _ := ctx.Value(ContextKeyRequestID).(string)
		a = append(a, uri, requestId)

		metas, _ := ctx.Value(ContextKeyMetas).(string)
		for _, meta := range strings.Split(metas, ",") {
			v, _ := ctx.Value(meta).(string)
			a = append(a, String(meta, v))
		}
		return
	}()

	return GetSugaredLogger().With(args...).WithOptions(AddCallerSkip(-1))
}
