package zlog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type (
	Field  = zap.Field
	Logger = zap.Logger
)

var (
	AddCallerSkip = zap.AddCallerSkip
	String        = zap.String
	Int           = zap.Int
	Int64         = zap.Int64
	Float64       = zap.Float64
)

var (
	zapLogger     *zap.Logger
	sugaredLogger *zap.SugaredLogger
)

// log文件后缀类型
const (
	txtLogStdout          = 0
	txtLogNormal          = 1 // 正常的日志：info、debug
	txtLogWarnFatal       = 2 // 异常的日志： warn、error、fatal
	txtLogRotateNormal    = 3
	txtLogRotateWarnFatal = 4
)

func newLogger() *zap.Logger {
	var stdLevel = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= logConfig.ZapLevel && lvl >= DebugLevel
	})

	var errLevel = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= logConfig.ZapLevel && lvl >= WarnLevel
	})

	var infoLevel = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= logConfig.ZapLevel && lvl <= InfoLevel
	})

	var zapCores []zapcore.Core
	if logConfig.Log2Stdout {
		var encoder zapcore.Encoder
		if logConfig.EncoderType == "console" {
			encoder = getConsoleEncoder()
		} else {
			encoder = getJsonEncoder()
		}
		zapCores = append(zapCores, zapcore.NewCore(encoder, getLogWriter(txtLogStdout), stdLevel))
	}

	if logConfig.Log2File {
		if logConfig.LogRotate {
			zapCores = append(zapCores, zapcore.NewCore(getJsonEncoder(), getLogWriter(txtLogRotateNormal), infoLevel))
			zapCores = append(zapCores, zapcore.NewCore(getJsonEncoder(), getLogWriter(txtLogRotateWarnFatal), errLevel))
		} else {
			zapCores = append(zapCores, zapcore.NewCore(getJsonEncoder(), getLogWriter(txtLogNormal), infoLevel))
			zapCores = append(zapCores, zapcore.NewCore(getJsonEncoder(), getLogWriter(txtLogWarnFatal), errLevel))
		}
	}

	core := zapcore.NewTee(zapCores...)

	return zap.New(core, zap.AddCaller(), zap.Fields(), zap.Development())
}

func getConsoleEncoder() zapcore.Encoder {
	// time字段编码器
	timeEncoder := zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.999999")

	encoderCfg := zapcore.EncoderConfig{
		LevelKey:       "level",
		TimeKey:        "time",
		CallerKey:      "file",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 短路径编码器
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     timeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}
	return zapcore.NewConsoleEncoder(encoderCfg)
}

func getJsonEncoder() zapcore.Encoder {
	// time字段编码器
	timeEncoder := zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.999999")

	encoderCfg := zapcore.EncoderConfig{
		LevelKey:       "level",
		TimeKey:        "time",
		CallerKey:      "file",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 短路径编码器
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     timeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}

	return zapcore.NewJSONEncoder(encoderCfg) // todo: custom webkit json_encoder
}

func getLogWriter(loggerType int8) (ws zapcore.WriteSyncer) {
	var w io.Writer
	if loggerType == txtLogStdout {
		w = os.Stdout
	} else if loggerType < 3 {
		// 打印到 name.log[.wf] 中
		var err error
		filename := filepath.Join(strings.TrimSuffix(logConfig.Path, "/"), appendLogFileTail(logConfig.AppName, loggerType))
		w, err = os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			panic("open log file error: " + err.Error())
		}
	} else {
		// Configure lumberjack for log rotation
		w = &lumberjack.Logger{
			Filename:   filepath.Join(strings.TrimSuffix(logConfig.Path, "/"), appendLogFileTail(logConfig.AppName, loggerType)),
			MaxSize:    logConfig.LogMaxSize,
			MaxBackups: logConfig.MaxBackups,
			MaxAge:     logConfig.MaxAge,
			Compress:   logConfig.Compress,
		}
	}

	if !logConfig.BufferSwitch {
		return zapcore.AddSync(w)
	}

	// 开启缓冲区
	ws = &zapcore.BufferedWriteSyncer{
		WS:            zapcore.AddSync(w),
		Size:          logConfig.BufferSize,
		FlushInterval: logConfig.BufferFlushInterval,
		Clock:         nil,
	}
	return ws
}

// genFilename 拼装完整文件名
func appendLogFileTail(appName string, loggerType int8) string {
	var tailFixed string
	switch loggerType {
	case txtLogNormal:
		tailFixed = ".log"
	case txtLogRotateNormal:
		tailFixed = ".log"
	case txtLogWarnFatal:
		tailFixed = ".log.wf"
	case txtLogRotateWarnFatal:
		tailFixed = ".log.wf"
	default:
		tailFixed = ".log"
	}
	return appName + tailFixed
}
