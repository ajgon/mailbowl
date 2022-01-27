package config_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/ajgon/mailbowl/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

const (
	colorDebug = "\x1b[35mDEBUG\x1b[0m"
	colorInfo  = "\x1b[34mINFO\x1b[0m"
	colorWarn  = "\x1b[33mWARN\x1b[0m"
	colorError = "\x1b[31mERROR\x1b[0m"
)

func captureStreams(conf config.Log) (string, string, error) {
	var stdout, stderr bytes.Buffer
	stdoutSync := zapcore.AddSync(&stdout)
	stderrSync := zapcore.AddSync(&stderr)

	logger, err := config.BuildLogger(conf, stdoutSync, stderrSync)
	if err != nil {
		return "", "", fmt.Errorf("error capturing streams: %w", err)
	}

	logger.Debug("testdebug")
	logger.Info("testinfo")
	logger.Warn("testwarn")
	logger.Error("testerror")

	return stdout.String(), stderr.String(), nil
}

func TestValidLogMarshalFromObjects(t *testing.T) {
	t.Parallel()

	viperConfig := viper.New()
	viperConfig.Set("log.color", "true")
	viperConfig.Set("log.format", "console")
	viperConfig.Set("log.level", "debug")
	viperConfig.Set("log.stacktrace_level", "info")
	conf, err := InitConfig(viperConfig)
	assert.NoError(t, err)

	assert.Equal(t, true, conf.Log.Color)
	assert.Equal(t, config.Console, conf.Log.Format)
	assert.Equal(t, zapcore.DebugLevel, conf.Log.Level)
	assert.Equal(t, zapcore.InfoLevel, conf.Log.StacktraceLevel)
}

func TestValidLogMarshalFromYAML(t *testing.T) {
	t.Parallel()

	yamlExample := `---
log:
  color: true
  format: console
  level: debug
  stacktrace_level: info
`

	viperConfig := viper.New()

	conf, err := InitConfig(viperConfig, yamlExample)
	assert.NoError(t, err)

	assert.Equal(t, true, conf.Log.Color)
	assert.Equal(t, config.Console, conf.Log.Format)
	assert.Equal(t, zapcore.DebugLevel, conf.Log.Level)
	assert.Equal(t, zapcore.InfoLevel, conf.Log.StacktraceLevel)
}

func TestValidLogMarshalFromENV(t *testing.T) {
	t.Setenv("LOG_COLOR", "true")
	t.Setenv("LOG_FORMAT", "console")
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("LOG_STACKTRACE_LEVEL", "info")

	viperConfig := viper.New()
	conf, err := InitConfig(viperConfig)
	assert.NoError(t, err)

	assert.Equal(t, true, conf.Log.Color)
	assert.Equal(t, config.Console, conf.Log.Format)
	assert.Equal(t, zapcore.DebugLevel, conf.Log.Level)
	assert.Equal(t, zapcore.InfoLevel, conf.Log.StacktraceLevel)
}

func TestInvalidFormat(t *testing.T) {
	t.Parallel()

	viperConfig := viper.New()
	viperConfig.Set("log.format", "wrong")
	viperConfig.Set("log.level", "debug")
	viperConfig.Set("log.stacktrace_level", "info")

	_, err := InitConfig(viperConfig)
	assert.EqualError(
		t, err, "error unmarshaling config: 1 error(s) decoding:\n\n* error decoding 'Log': invalid log format",
	)
}

func TestInvalidLevelDisablesLog(t *testing.T) {
	t.Parallel()

	viperConfig := viper.New()
	viperConfig.Set("log.format", "console")
	viperConfig.Set("log.level", "wrong")
	viperConfig.Set("log.stacktrace_level", "info")

	conf, err := InitConfig(viperConfig)
	assert.NoError(t, err)
	assert.Equal(t, zapcore.Level(100), conf.Log.Level)
}

func TestInvalidStacktraceLevel(t *testing.T) {
	t.Parallel()

	viperConfig := viper.New()
	viperConfig.Set("log.format", "console")
	viperConfig.Set("log.level", "debug")
	viperConfig.Set("log.stacktrace_level", "wrong")

	conf, err := InitConfig(viperConfig)
	assert.NoError(t, err)
	assert.Equal(t, zapcore.Level(100), conf.Log.StacktraceLevel)
}

func TestDefaultsLogOutput(t *testing.T) {
	t.Parallel()

	viperConfig := viper.New()
	conf, err := InitConfig(viperConfig)
	assert.NoError(t, err)

	stdout, stderr, err := captureStreams(conf.Log)
	if err != nil {
		t.Errorf("an error occurred, while none is expected: %s", err.Error())
	}

	assert.NotContains(t, stdout, "testdebug")
	assert.NotContains(t, stderr, "testdebug")
	assert.NotContains(t, stdout, "testinfo")
	assert.NotContains(t, stderr, "testinfo")
	assert.Contains(t, stdout, "\tWARN\t")
	assert.Contains(t, stdout, "testwarn")
	assert.NotContains(t, stdout, "testing/testing.go")
	assert.Contains(t, stderr, "\tERROR\t")
	assert.Contains(t, stderr, "testerror")
	assert.Contains(t, stderr, "testing/testing.go")
}

func TestMostVerboseWithColor(t *testing.T) {
	t.Parallel()

	viperConfig := viper.New()
	viperConfig.Set("log.color", true)
	viperConfig.Set("log.format", "console")
	viperConfig.Set("log.level", "debug")
	viperConfig.Set("log.stacktrace_level", "debug")
	conf, err := InitConfig(viperConfig)
	assert.NoError(t, err)

	stdout, stderr, err := captureStreams(conf.Log)
	assert.NoError(t, err)

	assert.Contains(t, stdout, colorDebug)
	assert.Contains(t, stdout, "testdebug")
	assert.Contains(t, stdout, "testing/testing.go")
	assert.Contains(t, stdout, colorInfo)
	assert.Contains(t, stdout, "testinfo")
	assert.Contains(t, stdout, "testing/testing.go")
	assert.Contains(t, stdout, colorWarn)
	assert.Contains(t, stdout, "testwarn")
	assert.Contains(t, stdout, "testing/testing.go")
	assert.Contains(t, stderr, colorError)
	assert.Contains(t, stderr, "testerror")
	assert.Contains(t, stderr, "testing/testing.go")
}

func TestLogFormatJSON(t *testing.T) {
	t.Parallel()

	viperConfig := viper.New()
	viperConfig.Set("log.format", "json")
	viperConfig.Set("log.level", "debug")
	viperConfig.Set("log.stacktrace_level", "debug")
	conf, err := InitConfig(viperConfig)
	assert.NoError(t, err)

	stdout, stderr, err := captureStreams(conf.Log)
	assert.NoError(t, err)

	assert.Contains(t, stdout, `"level":"DEBUG"`)
	assert.Contains(t, stdout, `"msg":"testdebug"`)
	assert.Contains(t, stdout, `"stacktrace":"github.com/Masterminds/log-go`)
	assert.Contains(t, stdout, `"level":"INFO"`)
	assert.Contains(t, stdout, `"msg":"testinfo"`)
	assert.Contains(t, stdout, `"stacktrace":"github.com/Masterminds/log-go`)
	assert.Contains(t, stdout, `"level":"WARN"`)
	assert.Contains(t, stdout, `"msg":"testwarn"`)
	assert.Contains(t, stdout, `"stacktrace":"github.com/Masterminds/log-go`)
	assert.Contains(t, stderr, `"level":"ERROR"`)
	assert.Contains(t, stderr, `"msg":"testerror"`)
	assert.Contains(t, stderr, `"stacktrace":"github.com/Masterminds/log-go`)
}

func TestLogFormatConsole(t *testing.T) {
	t.Parallel()

	viperConfig := viper.New()
	viperConfig.Set("log.format", "console")
	conf, err := InitConfig(viperConfig)
	assert.NoError(t, err)

	stdout, stderr, err := captureStreams(conf.Log)
	assert.NoError(t, err)

	assert.Contains(t, stdout, "\tWARN\t")
	assert.Contains(t, stdout, "testwarn")
	assert.NotContains(t, stdout, "testing/testing.go")
	assert.Contains(t, stderr, "\tERROR\t")
	assert.Contains(t, stderr, "testerror")
	assert.Contains(t, stderr, "testing/testing.go")
}

func TestLogFormatLogfmt(t *testing.T) {
	t.Parallel()

	viperConfig := viper.New()
	viperConfig.Set("log.format", "logfmt")
	conf, err := InitConfig(viperConfig)
	assert.NoError(t, err)

	stdout, stderr, err := captureStreams(conf.Log)
	assert.NoError(t, err)

	assert.Contains(t, stdout, `level=WARN`)
	assert.Contains(t, stdout, `msg=testwarn`)
	assert.NotContains(t, stdout, "testing/testing.go")
	assert.Contains(t, stderr, `level=ERROR`)
	assert.Contains(t, stderr, `msg=testerror`)
	assert.Contains(t, stderr, "testing/testing.go")
}
