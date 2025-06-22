package logger

import (
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Level           string   `mapstructure:"level"`
	OutPaths        []string `mapstructure:"out_paths"`
	ErrPaths        []string `mapstructure:"err_paths"`
	Colorized       bool     `mapstructure:"colorized"`
	Encoding        string   `mapstructure:"encoding"`
	TrackIterations bool     `mapstructure:"track_iterations"`
}

type Logger interface {
	Info(msg string, args ...interface{})
	Sync() error
	Error(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	WithOptions(opts ...zap.Option) Logger
	StartTimer(msg string, args ...interface{})
	StopTimer(msg string, args ...interface{}) error
}

type logger struct {
	log             *zap.Logger
	trackIterations bool
	stack           *Stack
}

func (l logger) WithOptions(opts ...zap.Option) Logger {
	return &logger{
		log:             l.log.WithOptions(opts...),
		trackIterations: l.trackIterations,
		stack:           l.stack,
	}
}

func (l logger) Sync() error {
	return l.log.Sync()
}

func generateLogFields(keysAndValues []interface{}) []zapcore.Field {
	numKeysAndValues := len(keysAndValues)
	res := make([]zapcore.Field, 0, numKeysAndValues/2)
	for i := 0; i < numKeysAndValues; i += 2 {
		key := keysAndValues[i].(string)
		val := keysAndValues[i+1]
		res = append(res, zap.Any(key, val))
	}
	return res
}

func (l logger) Info(msg string, args ...interface{}) {
	l.log.Info(msg, generateLogFields(args)...)
}

func (l logger) Error(msg string, args ...interface{}) {
	l.log.Error(msg, generateLogFields(args)...)
}

func (l logger) Debug(msg string, args ...interface{}) {
	l.log.Debug(msg, generateLogFields(args)...)
}

func (l logger) Warn(msg string, args ...interface{}) {
	l.log.Warn(msg, generateLogFields(args)...)
}

func (l logger) StartTimer(msg string, args ...interface{}) {
	if !l.trackIterations {
		return
	}
	start := time.Now()
	indent := strings.Repeat("-", l.stack.Len())
	resArgs := make([]zapcore.Field, 0, len(args)+1)
	resArgs = append(resArgs, zap.Time("start", start))
	resArgs = append(resArgs, generateLogFields(args)...)
	l.log.Info(fmt.Sprintf("⏱  %s-> %s", indent, msg), resArgs...)
	l.stack.Push(start)
}

func (l logger) StopTimer(msg string, args ...interface{}) error {
	if !l.trackIterations {
		return nil
	}
	startTime, err := l.stack.Pop()
	if err != nil {
		return err
	}
	indent := strings.Repeat("+", l.stack.Len())
	elapsed := time.Since(startTime)
	resArgs := make([]zapcore.Field, 0, len(args)+1)
	resArgs = append(resArgs, zap.Duration("elapsed", elapsed))
	resArgs = append(resArgs, generateLogFields(args)...)
	l.log.Info(fmt.Sprintf("⏱  %s+> %s", indent, msg), resArgs...)
	return nil
}

func New(conf *Config) (Logger, error) {
	logLevel, err := zap.ParseAtomicLevel(conf.Level)
	if err != nil {
		return nil, err
	}

	cfg := zap.Config{
		Level:            logLevel,
		Encoding:         conf.Encoding,
		Development:      false,
		OutputPaths:      conf.OutPaths,
		ErrorOutputPaths: conf.ErrPaths,
		// "initialFields": {"foo": "bar"},
		EncoderConfig: zap.NewProductionEncoderConfig(),
	}

	if conf.Colorized {
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // zapcore.TimeEncoderOfLayout(time.RFC3339)

	zapLogger, err := cfg.Build(zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))
	if err != nil {
		return nil, err
	}
	return &logger{
		log:             zapLogger,
		trackIterations: conf.TrackIterations,
		stack:           NewStack(),
	}, nil
}
