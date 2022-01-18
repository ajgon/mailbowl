package config

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/Masterminds/log-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//nolint:gochecknoglobals
var redactedFields = []string{
	"key",
	"password",
	"password_hash",
}

// Init reads in config file and ENV variables if set.
func Init() {
	SetDefaults()

	viper.SetEnvPrefix("MAILBOWL")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	err := ConfigureLogger(os.Stdout, os.Stderr)
	cobra.CheckErr(err)

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		Reload()
		log.Info("Using config file:", viper.ConfigFileUsed())
	} else {
		log.Errorf("problem using config file `%s`: %s", viper.ConfigFileUsed(), err.Error())
	}

	normalizeSliceOfMaps("smtp.auth.users")

	viperSettingsJSON, err := json.MarshalIndent(viper.AllSettings(), "", "  ")
	cobra.CheckErr(err)
	log.Debug("Loaded config: ", redactConfig(string(viperSettingsJSON)))
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

func CobraInitialize(cfgFile string) func() {
	return func() {
		LoadConfigFile(cfgFile)
		Init()
	}
}

func normalizeSliceOfMaps(viperPath string) {
	normalizedSlice := make([]map[string]string, 0)

	if slice, ok := viper.Get(viperPath).([]interface{}); ok {
		normalizedMap := make(map[string]string)

		for _, hash := range slice {
			for key, value := range hash.(map[interface{}]interface{}) {
				if normalizedKey, ok := key.(string); ok {
					if normalizedValue, ok := value.(string); ok {
						normalizedMap[normalizedKey] = normalizedValue
					}
				}
			}
		}

		normalizedSlice = append(normalizedSlice, normalizedMap)
	}

	viper.Set(viperPath, normalizedSlice)
}

func redactConfig(configJSON string) string {
	for _, field := range redactedFields {
		re := regexp.MustCompile(fmt.Sprintf(`"%s":(\s*)"[^"]*"`, field))
		configJSON = re.ReplaceAllString(configJSON, fmt.Sprintf(`"%s":$1"[REDACTED]"`, field))
	}

	return configJSON
}
