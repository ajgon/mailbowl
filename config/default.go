package config

import (
	"reflect"

	"github.com/spf13/viper"
)

type Defaults struct {
	Env                string `viper:"env"`
	LogFormat          string `viper:"log.format"`
	LogLevel           string `viper:"log.level"`
	LogStacktraceLevel string `viper:"log.stacktrace_level"`
}

func GetDefaults() *Defaults {
	return &Defaults{
		Env:                "development",
		LogFormat:          "",
		LogLevel:           "warn",
		LogStacktraceLevel: "",
	}
}

func SetDefaults(force ...bool) {
	forceDefaults := false

	if len(force) == 1 {
		forceDefaults = force[0]
	}

	defaults := GetDefaults()

	defaultsTypes := reflect.TypeOf(*defaults)
	defaultsValues := reflect.ValueOf(*defaults)

	for i := 0; i < defaultsTypes.NumField(); i++ {
		structType := defaultsTypes.Field(i)
		structValue := defaultsValues.Field(i)

		viper.SetDefault(structType.Tag.Get("viper"), structValue.Interface())

		if forceDefaults {
			viper.Set(structType.Tag.Get("viper"), structValue.Interface())
		}
	}
}
