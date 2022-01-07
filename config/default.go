package config

import (
	"reflect"

	"github.com/spf13/viper"
)

type Defaults struct {
	Env string `viper:"env"`

	LogFormat          string `viper:"log.format"`
	LogLevel           string `viper:"log.level"`
	LogStacktraceLevel string `viper:"log.stacktrace_level"`

	SMTPHostname            string   `viper:"smtp.hostname"`
	SMTPLimitConnections    int      `viper:"smtp.limit.connections"`
	SMTPLimitMessageSize    int      `viper:"smtp.limit.message_size"`
	SMTPLimitRecipients     int      `viper:"smtp.limit.recipients"`
	SMTPListen              []string `viper:"smtp.listen"`
	SMTPTimeoutRead         string   `viper:"smtp.timeout.read"`
	SMTPTimeoutWrite        string   `viper:"smtp.timeout.write"`
	SMTPTimeoutData         string   `viper:"smtp.timeout.data"`
	SMTPTLSKey              string   `viper:"smtp.tls.key"`
	SMTPTLSCertificate      string   `viper:"smtp.tls.certificate"`
	SMTPTLSKeyFile          string   `viper:"smtp.tls.key_file"`
	SMTPTLSCertificateFile  string   `viper:"smtp.tls.certificate_file"`
	SMTPTLSForceForStartTLS bool     `viper:"smtp.tls.force_for_starttls"`
	SMTPWhitelistCIDRs      []string `viper:"smtp.whitelist.cidrs"`
}

func GetDefaults() *Defaults {
	return &Defaults{
		Env: "development",

		LogFormat:          "",
		LogLevel:           "warn",
		LogStacktraceLevel: "",

		SMTPHostname: "localhost.localdomain",

		SMTPLimitConnections: 100,      //nolint:gomnd
		SMTPLimitMessageSize: 26214400, //nolint:gomnd
		SMTPLimitRecipients:  100,      //nolint:gomnd
		// SMTPListen:           []string{"plain://0.0.0.0:10025", "tls://0.0.0.0:10465", "starttls://0.0.0.0:10587"},
		SMTPListen:              []string{"plain://0.0.0.0:10025"},
		SMTPTimeoutRead:         "60s",
		SMTPTimeoutWrite:        "60s",
		SMTPTimeoutData:         "5m",
		SMTPTLSKey:              "",
		SMTPTLSCertificate:      "",
		SMTPTLSForceForStartTLS: true,
		SMTPTLSKeyFile:          "",
		SMTPTLSCertificateFile:  "",
		SMTPWhitelistCIDRs:      []string{},
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
