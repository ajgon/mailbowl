package config

import (
	"errors"
	"fmt"
)

var (
	errLogFormat          = errors.New("invalid log format")
	errLogLevel           = errors.New("invalid log level")
	errLogStacktraceLevel = errors.New("invalid log stacktrace level")
)

func LogFormatError(logFormat string) error {
	return fmt.Errorf("%w: `%s`, allowed formats: `console`, `json` or `logfmt`", errLogFormat, logFormat)
}

func LogLevelError(logLevel string) error {
	return fmt.Errorf("%w: `%s`, allowed levels: `debug`, `info`, `warn` or `error`", errLogLevel, logLevel)
}

func LogStacktraceLevelError(logLevel string) error {
	return fmt.Errorf("%w: `%s`, allowed levels: `debug`, `info`, `warn` or `error`", errLogStacktraceLevel, logLevel)
}
