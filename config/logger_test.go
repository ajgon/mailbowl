//nolint:paralleltest
package config_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/Masterminds/log-go"
	"github.com/ajgon/mailbowl/config"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
)

const (
	colorDebug = "\x1b[35mDEBUG\x1b[0m"
	colorInfo  = "\x1b[34mINFO\x1b[0m"
	colorWarn  = "\x1b[33mWARN\x1b[0m"
	colorError = "\x1b[31mERROR\x1b[0m"
)

func captureStreams() (string, string, error) {
	var stdout, stderr bytes.Buffer
	stdoutSync := zapcore.AddSync(&stdout)
	stderrSync := zapcore.AddSync(&stderr)

	err := config.ConfigureLogger(stdoutSync, stderrSync)
	if err != nil {
		return "", "", fmt.Errorf("error capturing streams: %w", err)
	}

	log.Debug("testdebug")
	log.Info("testinfo")
	log.Warn("testwarn")
	log.Error("testerror")

	return stdout.String(), stderr.String(), nil
}

func assertIncludes(t *testing.T, streamName, haystack, needle string) {
	t.Helper()

	if strings.Contains(haystack, needle) {
		return
	}

	t.Errorf("expected %s to include `%s`", streamName, needle)
}

func assertExcludes(t *testing.T, streamName, haystack, needle string) {
	t.Helper()

	if !strings.Contains(haystack, needle) {
		return
	}

	t.Errorf("expected %s to not include `%s`", streamName, needle)
}

func TestDefaults(t *testing.T) {
	config.SetDefaults(true)

	stdout, stderr, err := captureStreams()
	if err != nil {
		t.Errorf("an error occurred, while none is expected: %s", err.Error())
	}

	assertExcludes(t, "stdout", stdout, "testdebug")
	assertExcludes(t, "stderr", stderr, "testdebug")
	assertExcludes(t, "stdout", stdout, "testinfo")
	assertExcludes(t, "stderr", stderr, "testinfo")
	assertIncludes(t, "stdout", stdout, colorWarn)
	assertIncludes(t, "stdout", stdout, "testwarn")
	assertExcludes(t, "stdout", stdout, "testing/testing.go")
	assertIncludes(t, "stderr", stderr, colorError)
	assertIncludes(t, "stderr", stderr, "testerror")
	assertIncludes(t, "stderr", stderr, "testing/testing.go")
}

func TestDevelopmentMostVerbose(t *testing.T) {
	config.SetDefaults(true)
	viper.Set("env", "development")
	viper.Set("log.format", "console")
	viper.Set("log.level", "debug")
	viper.Set("log.stacktrace_level", "debug")

	stdout, stderr, err := captureStreams()
	if err != nil {
		t.Errorf("an error occurred, while none is expected: %s", err.Error())
	}

	assertIncludes(t, "stdout", stdout, colorDebug)
	assertIncludes(t, "stdout", stdout, "testdebug")
	assertIncludes(t, "stdout", stdout, "testing/testing.go")
	assertIncludes(t, "stdout", stdout, colorInfo)
	assertIncludes(t, "stdout", stdout, "testinfo")
	assertIncludes(t, "stdout", stdout, "testing/testing.go")
	assertIncludes(t, "stdout", stdout, colorWarn)
	assertIncludes(t, "stdout", stdout, "testwarn")
	assertIncludes(t, "stdout", stdout, "testing/testing.go")
	assertIncludes(t, "stderr", stderr, colorError)
	assertIncludes(t, "stderr", stderr, "testerror")
	assertIncludes(t, "stderr", stderr, "testing/testing.go")
}

func TestDevelopmentJSON(t *testing.T) {
	config.SetDefaults(true)
	viper.Set("log.format", "json")

	stdout, stderr, err := captureStreams()
	if err != nil {
		t.Errorf("an error occurred, while none is expected: %s", err.Error())
	}

	assertIncludes(t, "stdout", stdout, `"L":"WARN"`)
	assertIncludes(t, "stdout", stdout, `"M":"testwarn"`)
	assertExcludes(t, "stdout", stdout, `"S":"github.com/Masterminds/log-go`)
	assertIncludes(t, "stderr", stderr, `"L":"ERROR"`)
	assertIncludes(t, "stderr", stderr, `"M":"testerror"`)
	assertIncludes(t, "stderr", stderr, `"S":"github.com/Masterminds/log-go`)
}

func TestDevelopmentLogfmt(t *testing.T) {
	config.SetDefaults(true)
	viper.Set("log.format", "logfmt")

	stdout, stderr, err := captureStreams()
	if err != nil {
		t.Errorf("an error occurred, while none is expected: %s", err.Error())
	}

	assertIncludes(t, "stdout", stdout, `L=WARN`)
	assertIncludes(t, "stdout", stdout, `M=testwarn`)
	assertExcludes(t, "stdout", stdout, `S="github.com/Masterminds/log-go`)
	assertIncludes(t, "stderr", stderr, `L=ERROR`)
	assertIncludes(t, "stderr", stderr, `M=testerror`)
	assertIncludes(t, "stderr", stderr, `S="github.com/Masterminds/log-go`)
}

func TestProductionDefaults(t *testing.T) {
	config.SetDefaults(true)
	viper.Set("env", "production")

	stdout, stderr, err := captureStreams()
	if err != nil {
		t.Errorf("an error occurred, while none is expected: %s", err.Error())
	}

	assertIncludes(t, "stdout", stdout, `"level":"warn"`)
	assertIncludes(t, "stdout", stdout, `"msg":"testwarn"`)
	assertExcludes(t, "stdout", stdout, "testing/testing.go")
	assertIncludes(t, "stderr", stderr, `"level":"error"`)
	assertIncludes(t, "stderr", stderr, `"msg":"testerror"`)
	assertExcludes(t, "stderr", stderr, "testing/testing.go")
}

func TestProductionMostVerbose(t *testing.T) {
	config.SetDefaults(true)
	viper.Set("env", "production")
	viper.Set("log.format", "json")
	viper.Set("log.level", "debug")
	viper.Set("log.stacktrace_level", "debug")

	stdout, stderr, err := captureStreams()
	if err != nil {
		t.Errorf("an error occurred, while none is expected: %s", err.Error())
	}

	assertIncludes(t, "stdout", stdout, `"level":"debug"`)
	assertIncludes(t, "stdout", stdout, `"msg":"testdebug"`)
	assertIncludes(t, "stdout", stdout, `"stacktrace":"github.com/Masterminds/log-go`)
	assertIncludes(t, "stdout", stdout, `"level":"info"`)
	assertIncludes(t, "stdout", stdout, `"msg":"testinfo"`)
	assertIncludes(t, "stdout", stdout, `"stacktrace":"github.com/Masterminds/log-go`)
	assertIncludes(t, "stdout", stdout, `"level":"warn"`)
	assertIncludes(t, "stdout", stdout, `"msg":"testwarn"`)
	assertIncludes(t, "stdout", stdout, `"stacktrace":"github.com/Masterminds/log-go`)
	assertIncludes(t, "stderr", stderr, `"level":"error"`)
	assertIncludes(t, "stderr", stderr, `"msg":"testerror"`)
	assertIncludes(t, "stderr", stderr, `"stacktrace":"github.com/Masterminds/log-go`)
}

func TestProductionConsole(t *testing.T) {
	config.SetDefaults(true)
	viper.Set("env", "production")
	viper.Set("log.format", "console")

	stdout, stderr, err := captureStreams()
	if err != nil {
		t.Errorf("an error occurred, while none is expected: %s", err.Error())
	}

	assertIncludes(t, "stdout", stdout, "\twarn\t")
	assertIncludes(t, "stdout", stdout, "testwarn")
	assertExcludes(t, "stdout", stdout, "testing/testing.go")
	assertIncludes(t, "stderr", stderr, "\terror\t")
	assertIncludes(t, "stderr", stderr, "testerror")
	assertExcludes(t, "stderr", stderr, "testing/testing.go")
}

func TestProductionLogfmt(t *testing.T) {
	config.SetDefaults(true)
	viper.Set("env", "production")
	viper.Set("log.format", "logfmt")

	stdout, stderr, err := captureStreams()
	if err != nil {
		t.Errorf("an error occurred, while none is expected: %s", err.Error())
	}

	assertIncludes(t, "stdout", stdout, `level=warn`)
	assertIncludes(t, "stdout", stdout, `msg=testwarn`)
	assertExcludes(t, "stdout", stdout, "testing/testing.go")
	assertIncludes(t, "stderr", stderr, `level=error`)
	assertIncludes(t, "stderr", stderr, `msg=testerror`)
	assertExcludes(t, "stderr", stderr, "testing/testing.go")
}

func TestIncorrectLogFormat(t *testing.T) {
	config.SetDefaults(true)
	viper.Set("log.format", "wrong")

	_, _, err := captureStreams()

	if !strings.Contains(err.Error(), "invalid log format: `wrong`, allowed formats: `console`, `json` or `logfmt`") {
		t.Errorf("invalid error message for log format")
	}
}

func TestIncorrectLogLevel(t *testing.T) {
	config.SetDefaults(true)
	viper.Set("log.level", "wrong")

	_, _, err := captureStreams()

	if !strings.Contains(err.Error(), "invalid log level: `wrong`, allowed levels: `debug`, `info`, `warn` or `error`") {
		t.Errorf("invalid error message for log format")
	}
}

func TestIncorrectStacktraceLogLevel(t *testing.T) {
	config.SetDefaults(true)
	viper.Set("log.stacktrace_level", "wrong")

	_, _, err := captureStreams()

	if !strings.Contains(
		err.Error(), "invalid log stacktrace level: `wrong`, allowed levels: `debug`, `info`, `warn` or `error`",
	) {
		t.Errorf("invalid error message for log format")
	}
}
