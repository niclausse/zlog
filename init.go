package zlog

import (
	"fmt"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

type StringLogLevel string

const (
	DebugStringLogLevel StringLogLevel = "debug"
	InfoStringLogLevel  StringLogLevel = "info"
	WarnStringLogLevel  StringLogLevel = "warn"
	ErrorStringLogLevel StringLogLevel = "error"
	FatalStringLogLevel StringLogLevel = "fatal"
)

func (l StringLogLevel) Level() Level {
	if l == DebugStringLogLevel {
		return DebugLevel
	}
	if l == InfoStringLogLevel {
		return InfoLevel
	}
	if l == WarnStringLogLevel {
		return WarnLevel
	}
	if l == ErrorStringLogLevel {
		return ErrorLevel
	}
	if l == FatalStringLogLevel {
		return FatalLevel
	}

	return DebugLevel
}

type Level = zapcore.Level

const (
	DebugLevel = zapcore.DebugLevel
	InfoLevel  = zapcore.InfoLevel
	WarnLevel  = zapcore.WarnLevel
	ErrorLevel = zapcore.ErrorLevel
	FatalLevel = zapcore.FatalLevel
)

const (
	defaultLogPath      = "./log"
	defaultAppName      = "app"
	ContextKeyURI       = "_uri"
	ContextKeyRequestID = "_zlog_requestId"
	ContextKeyMetas     = "_zlog_metas"
)

type config struct {
	ZapLevel    zapcore.Level
	AppName     string
	EncoderType string // "console", "json"

	Log2Stdout bool // default true
	Log2File   bool // default false
	LogRotate  bool // default false
	LogMaxSize int  // default 100M
	MaxBackups int  // default no limit
	MaxAge     int
	Compress   bool   // default false
	Path       string // default ./log

	// 缓冲区
	BufferSwitch        bool
	BufferSize          int
	BufferFlushInterval time.Duration
}

// 全局配置 仅限InitLogConfig函数进行变更
var logConfig *config

func init() {
	logConfig = &config{
		ZapLevel:    zapcore.InfoLevel,
		AppName:     defaultAppName,
		EncoderType: "",

		Log2Stdout: true,
		Log2File:   false,
		Path:       defaultLogPath,

		// 缓冲区，如果不配置默认使用以下配置
		BufferSwitch:        false,
		BufferSize:          256 * 1024, // 256kb
		BufferFlushInterval: 5 * time.Second,
	}

	InitLogConfig()
}

type Option func()

func WithAppName(appName string) Option {
	return func() {
		logConfig.AppName = appName
	}
}

func WithLevel(level Level) Option {
	return func() {
		logConfig.ZapLevel = level
	}
}

func WithLog2Stdout(stdout bool) Option {
	return func() {
		logConfig.Log2Stdout = stdout
	}
}

func WithLog2File(file bool) Option {
	return func() {
		logConfig.Log2File = file
	}
}

func WithLogRotate(rotate bool) Option {
	return func() {
		logConfig.LogRotate = rotate
	}
}

func WithMaxSize(size int) Option {
	return func() {
		logConfig.LogMaxSize = size
	}
}

func WithMaxBackups(backups int) Option {
	return func() {
		logConfig.MaxBackups = backups
	}
}

func WithMaxAge(age int) Option {
	return func() {
		logConfig.MaxAge = age
	}
}

func WithCompress(compress bool) Option {
	return func() {
		logConfig.Compress = compress
	}
}

func WithPath(path string) Option {
	return func() {
		logConfig.Path = path
	}
}

func WithLogDirPath(path string) Option {
	return func() {
		if len(path) > 0 {
			logConfig.Path = path
		}
	}
}

func WithBuffer(size int, flushInterval time.Duration) Option {
	return func() {
		logConfig.BufferSwitch = true
		logConfig.BufferSize = size
		logConfig.BufferFlushInterval = flushInterval
	}
}

func WithEncoderType(encoderType string) Option {
	return func() {
		logConfig.EncoderType = encoderType
	}
}

func InitLogConfig(opts ...Option) {
	for _, opt := range opts {
		opt()
	}

	// 目录不存在则先创建目录
	if logConfig.Log2File && logConfig.Path != "" {
		if _, err := os.Stat(logConfig.Path); os.IsNotExist(err) {
			err = os.MkdirAll(logConfig.Path, 0777)
			if err != nil {
				panic(fmt.Errorf("log conf err: create log dir '%s' error: %s", logConfig.Path, err))
			}
		}
	}
}
