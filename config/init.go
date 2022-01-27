package config

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/Masterminds/log-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//nolint:gochecknoglobals
var redactedFields = []string{
	"Key",
	"Password",
	"PasswordHash",
}

func Init(cfgFile string) {
	SetDefaults(false)
	viper.SetEnvPrefix("MAILBOWL")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if cfgFile != "" && err != nil {
		log.Fatalf("problem using config file `%s`: %s", viper.ConfigFileUsed(), err.Error())
	}

	config = Config{}

	if err := viper.Unmarshal(&config, viper.DecodeHook(Hook)); err != nil {
		log.Fatalf("error initializing config: %s", err.Error())
	}

	Reload()

	if err == nil {
		log.Infof("Using config file: %s", viper.ConfigFileUsed())
	}

	configJSON, err := json.MarshalIndent(Get(), "", "  ")
	cobra.CheckErr(err)
	log.Debug("Loaded config: ", redactConfig(string(configJSON)))
}

func LoadConfigFile(cfgFile string) {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)

		return
	}

	// Search config in home directory with name ".mailbowl" (without extension).
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.SetConfigName("mailbowl")
}

func CobraInitialize(cfgFile string) {
	LoadConfigFile(cfgFile)
	Init(cfgFile)
}

func redactConfig(configJSON string) string {
	for _, field := range redactedFields {
		re := regexp.MustCompile(fmt.Sprintf(`"%s":(\s*)"[^"]*"`, field))
		configJSON = re.ReplaceAllString(configJSON, fmt.Sprintf(`"%s":$1"[REDACTED]"`, field))
	}

	return configJSON
}
