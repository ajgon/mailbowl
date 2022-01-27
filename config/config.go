package config

import (
	"errors"
	"reflect"
)

var (
	ErrUnserializing = errors.New("error unserializing conf")
	config           Config //nolint:gochecknoglobals
)

type Config struct {
	Log   Log
	Relay Relay
	SMTP  SMTP
}

func Get() Config {
	return config
}

func Hook(dataType reflect.Type, targetDataType reflect.Type, rawData interface{}) (interface{}, error) {
	if dataType.Kind() != reflect.Map {
		return rawData, nil
	}

	if targetDataType == reflect.TypeOf(Log{}) {
		return LogHook(dataType, targetDataType, rawData)
	}

	if targetDataType == reflect.TypeOf(Relay{}) {
		return RelayHook(dataType, targetDataType, rawData)
	}

	if targetDataType == reflect.TypeOf(SMTP{}) {
		return SMTPHook(dataType, targetDataType, rawData)
	}

	return rawData, nil
}

func parseBoolString(boolString string) bool {
	return boolString == "true" || boolString == "1"
}
