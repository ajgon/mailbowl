package config_test

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ajgon/mailbowl/config"
	"github.com/spf13/viper"
)

func InitConfig(viperConfig *viper.Viper, yaml ...string) (*config.Config, error) {
	config.SetDefaultsForViper(false, viperConfig)
	viperConfig.AutomaticEnv()
	viperConfig.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viperConfig.SetConfigType("yaml")

	if len(yaml) > 0 {
		if err := viperConfig.ReadConfig(bytes.NewReader([]byte(yaml[0]))); err != nil {
			return nil, fmt.Errorf("error parsing yaml: %w", err)
		}
	}

	conf := config.Config{}

	if err := viperConfig.Unmarshal(&conf, viper.DecodeHook(config.Hook)); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &conf, nil
}
