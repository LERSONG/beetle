package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

var DefaultSugaredLogger *zap.SugaredLogger

func InitLogger(opts ...Option) {
	options := &Options{}
	for _, opt := range opts {
		opt(options)
	}

	if options.config == nil {
		lvl, err := getLevel(os.Getenv("LOG_LEVEL"))
		if err != nil {
			lvl = zapcore.InfoLevel
		}

		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		options.config = &zap.Config{
			Level:       zap.NewAtomicLevelAt(lvl),
			Development: false,
			Sampling: &zap.SamplingConfig{
				Initial:    100,
				Thereafter: 100,
			},
			Encoding:         "console",
			EncoderConfig:    encoderConfig,
			OutputPaths:      []string{"stderr"},
			ErrorOutputPaths: []string{"stderr"},
		}
	}

	logger, _ := options.config.Build(zap.AddCaller())
	DefaultSugaredLogger = logger.Sugar()
}

func GetLogger() *zap.SugaredLogger {
	return DefaultSugaredLogger
}

func Sync() error {
	return DefaultSugaredLogger.Sync()
}

func Debugf(template string, args ...interface{}) {
	DefaultSugaredLogger.Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	DefaultSugaredLogger.Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	DefaultSugaredLogger.Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	DefaultSugaredLogger.Errorf(template, args...)
}

func DPanicf(template string, args ...interface{}) {
	DefaultSugaredLogger.DPanicf(template, args...)
}

func Panicf(template string, args ...interface{}) {
	DefaultSugaredLogger.Panicf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	DefaultSugaredLogger.Fatalf(template, args...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	DefaultSugaredLogger.Debugw(msg, keysAndValues...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	DefaultSugaredLogger.Infow(msg, keysAndValues...)
}

func DPanicw(msg string, keysAndValues ...interface{}) {
	DefaultSugaredLogger.DPanicw(msg, keysAndValues...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	DefaultSugaredLogger.Warnw(msg, keysAndValues...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	DefaultSugaredLogger.Errorw(msg, keysAndValues...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	DefaultSugaredLogger.Panicw(msg, keysAndValues...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	DefaultSugaredLogger.Fatalw(msg, keysAndValues...)
}

// getLevel converts a level string into a logger Level value.
// returns an error if the input string does not match known values.
func getLevel(levelStr string) (zapcore.Level, error) {
	level := strings.ToLower(levelStr)
	switch level {
	case zapcore.DebugLevel.String():
		return zapcore.DebugLevel, nil
	case zapcore.InfoLevel.String():
		return zapcore.InfoLevel, nil
	case zapcore.WarnLevel.String():
		return zapcore.WarnLevel, nil
	case zapcore.ErrorLevel.String():
		return zapcore.ErrorLevel, nil
	case zapcore.DPanicLevel.String():
		return zapcore.DPanicLevel, nil
	case zapcore.PanicLevel.String():
		return zapcore.PanicLevel, nil
	case zapcore.FatalLevel.String():
		return zapcore.FatalLevel, nil
	}
	return zapcore.InfoLevel, fmt.Errorf("Unknown Level String: '%s', defaulting to InfoLevel", levelStr)
}
