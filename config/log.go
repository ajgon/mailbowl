package config

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/Masterminds/log-go"
	zaplog "github.com/Masterminds/log-go/impl/zap"
	zaplogfmt "github.com/jsternberg/zap-logfmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogFormat int

const (
	Console LogFormat = iota
	JSON
	Logfmt
)

const disabledLogLevel zapcore.Level = 100

var ErrInvalidLogFormat = errors.New("invalid log format")

type Log struct {
	Color           bool
	Format          LogFormat
	Level           zapcore.Level
	StacktraceLevel zapcore.Level
}

func ConfigureLogger(conf Log, stdout, stderr zapcore.WriteSyncer) (err error) {
	log.Current, err = BuildLogger(conf, stdout, stderr)

	if err != nil {
		return fmt.Errorf("error configuring logger: %w", err)
	}

	return nil
}

func BuildLogger(conf Log, stdout, stderr zapcore.WriteSyncer) (*zaplog.Logger, error) {
	var core zapcore.Core

	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel && lvl >= conf.Level
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel && lvl >= conf.Level
	})

	consoleDebugging := zapcore.Lock(stdout)
	consoleErrors := zapcore.Lock(stderr)

	encoder, err := buildEncoder(conf)
	if err != nil {
		return nil, fmt.Errorf("error while configuring logger: %w", err)
	}

	core = zapcore.NewTee(
		zapcore.NewCore(encoder, consoleErrors, highPriority),
		zapcore.NewCore(encoder, consoleDebugging, lowPriority),
	)

	opts := []zap.Option{
		zap.AddStacktrace(conf.StacktraceLevel),
		zap.AddCaller(),
	}

	logger := zap.New(core, opts...)
	defer logger.Sync() //nolint:errcheck

	return zaplog.New(logger), nil
}

//nolint:cyclop
func LogHook(dataType reflect.Type, targetDataType reflect.Type, rawData interface{}) (interface{}, error) {
	var (
		data                           map[string]interface{}
		err                            error
		ok                             bool
		format, level, stacktraceLevel string
	)

	if dataType.Kind() != reflect.Map {
		return rawData, nil
	}

	if targetDataType != reflect.TypeOf(Log{}) {
		return rawData, nil
	}

	if data, ok = rawData.(map[string]interface{}); !ok {
		return nil, ErrUnserializing
	}

	if format, ok = data["format"].(string); !ok {
		return nil, ErrUnserializing
	}

	if level, ok = data["level"].(string); !ok {
		return nil, ErrUnserializing
	}

	if stacktraceLevel, ok = data["stacktrace_level"].(string); !ok {
		return nil, ErrUnserializing
	}

	logConfig := &Log{}

	switch logColor := data["color"].(type) {
	case bool:
		logConfig.Color = logColor
	case string:
		logConfig.Color = parseBoolString(logColor)
	default:
		return nil, ErrUnserializing
	}

	logConfig.Format, err = buildLogFormat(format)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	err = logConfig.Level.UnmarshalText([]byte(level))
	if err != nil {
		logConfig.Level = disabledLogLevel
	}

	err = logConfig.StacktraceLevel.UnmarshalText([]byte(stacktraceLevel))
	if err != nil {
		logConfig.StacktraceLevel = disabledLogLevel
	}

	return logConfig, nil
}

func buildLogFormat(format string) (LogFormat, error) {
	switch format {
	case "console":
		return Console, nil
	case "json":
		return JSON, nil
	case "logfmt":
		return Logfmt, nil
	}

	return -1, ErrInvalidLogFormat
}

func buildEncoder(conf Log) (zapcore.Encoder, error) { //nolint:ireturn
	encoderConfig := zap.NewProductionEncoderConfig()

	if conf.Color {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	}

	switch conf.Format {
	case Console:
		return zapcore.NewConsoleEncoder(encoderConfig), nil
	case JSON:
		return zapcore.NewJSONEncoder(encoderConfig), nil
	case Logfmt:
		return zaplogfmt.NewEncoder(encoderConfig), nil
	}

	return nil, ErrInvalidLogFormat
}
