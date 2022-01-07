package config

import (
	"fmt"

	"github.com/Masterminds/log-go"
	zaplog "github.com/Masterminds/log-go/impl/zap"
	zaplogfmt "github.com/jsternberg/zap-logfmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	logFormatConsole = "console"
	logFormatJSON    = "json"
	logFormatLoggmt  = "logfmt"

	logLevelDebug = "debug"
	logLevelInfo  = "info"
	logLevelWarn  = "warn"
	logLevelError = "error"
)

func ConfigureLogger(stdout, stderr zapcore.WriteSyncer) error {
	var (
		core     zapcore.Core
		logLevel zapcore.Level
	)

	level, stacktraceLevel, err := extractLogLevels()
	if err != nil {
		return fmt.Errorf("error while configuring logger: %w", err)
	}

	err = logLevel.UnmarshalText([]byte(level))
	if err != nil {
		return fmt.Errorf("error while configuring logger: %w", err)
	}

	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel && lvl >= logLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel && lvl >= logLevel
	})

	consoleDebugging := zapcore.Lock(stdout)
	consoleErrors := zapcore.Lock(stderr)

	env := viper.GetString("env")

	encoder, err := buildEncoder(env, viper.GetString("log.format"))
	if err != nil {
		return fmt.Errorf("error while configuring logger: %w", err)
	}

	core = zapcore.NewTee(
		zapcore.NewCore(encoder, consoleErrors, highPriority),
		zapcore.NewCore(encoder, consoleDebugging, lowPriority),
	)

	opts, err := buildOpts(env, stacktraceLevel)
	if err != nil {
		return fmt.Errorf("error while configuring logger: %w", err)
	}

	logger := zap.New(core, opts...)
	defer logger.Sync() //nolint:errcheck

	log.Current = zaplog.New(logger)

	return nil
}

func buildOpts(env, stacktraceLevel string) ([]zap.Option, error) {
	var (
		stackLevel zapcore.Level
		opts       []zap.Option
	)

	if stacktraceLevel == "" && env == "development" {
		stacktraceLevel = "error"
	}

	err := stackLevel.UnmarshalText([]byte(stacktraceLevel))
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	if env == "production" {
		if stacktraceLevel != "" {
			opts = append(opts, zap.AddStacktrace(stackLevel))
		}
	} else {
		opts = append(opts, zap.Development())
		opts = append(opts, zap.AddStacktrace(stackLevel))
	}

	opts = append(opts, zap.AddCaller())

	return opts, nil
}

func buildEncoder(env, logFormat string) (zapcore.Encoder, error) { //nolint:ireturn
	var encoderConfig zapcore.EncoderConfig

	if env == "production" {
		encoderConfig = zap.NewProductionEncoderConfig()

		if logFormat == "" {
			logFormat = logFormatJSON
		}
	} else {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

		if logFormat == "" {
			logFormat = logFormatConsole
		}

		if logFormat == logFormatConsole {
			encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		}
	}

	switch logFormat {
	case logFormatConsole:
		return zapcore.NewConsoleEncoder(encoderConfig), nil
	case logFormatJSON:
		return zapcore.NewJSONEncoder(encoderConfig), nil
	case logFormatLoggmt:
		return zaplogfmt.NewEncoder(encoderConfig), nil
	}

	return nil, LogFormatError(logFormat)
}

func extractLogLevels() (string, string, error) {
	level := viper.GetString("log.level")
	if level == "" {
		level = GetDefaults().LogLevel
	}

	if !validLevel(level) {
		return "", "", LogLevelError(level)
	}

	stackLevel := viper.GetString("log.stacktrace_level")
	if stackLevel != "" && !validLevel(stackLevel) {
		return "", "", LogStacktraceLevelError(stackLevel)
	}

	return level, stackLevel, nil
}

func validLevel(level string) bool {
	return level == logLevelDebug || level == logLevelInfo || level == logLevelWarn || level == logLevelError
}
