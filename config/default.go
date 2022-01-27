package config

import "github.com/spf13/viper"

const (
	defaultConnectionsLimit   = 100
	defaultMessageSizeInBytes = 26214400
	defaultRecipientsLimit    = 100
)

//nolint:gochecknoglobals
var defaults = map[string]interface{}{
	"log.color":                             false,
	"log.format":                            "console",
	"log.level":                             "warn",
	"log.stacktrace_level":                  "error",
	"relay.outgoing_server.address":         "",
	"relay.outgoing_server.auth_method":     "plain",
	"relay.outgoing_server.connection_type": "tls",
	"relay.outgoing_server.from_email":      "",
	"relay.outgoing_server.host":            "",
	"relay.outgoing_server.password":        "",
	"relay.outgoing_server.port":            0,
	"relay.outgoing_server.username":        "",
	"relay.outgoing_server.verify_tls":      true,
	"smtp.auth.enabled":                     false,
	"smtp.auth.users":                       []interface{}{},
	"smtp.hostname":                         "",
	"smtp.limit.connections":                defaultConnectionsLimit,
	"smtp.limit.message_size":               defaultMessageSizeInBytes,
	"smtp.limit.recipients":                 defaultRecipientsLimit,
	"smtp.listen":                           []string{},
	"smtp.timeout.read":                     "60s",
	"smtp.timeout.write":                    "60s",
	"smtp.timeout.data":                     "5m",
	"smtp.tls.key":                          "",
	"smtp.tls.certificate":                  "",
	"smtp.tls.key_file":                     "",
	"smtp.tls.certificate_file":             "",
	"smtp.tls.force_for_starttls":           true,
	"smtp.whitelist":                        []string{},
}

func SetDefaults(force bool) {
	for key, value := range defaults {
		if force {
			viper.Set(key, value)
		} else {
			viper.SetDefault(key, value)
		}
	}
}

func SetDefaultsForViper(force bool, v *viper.Viper) {
	for key, value := range defaults {
		if force {
			v.Set(key, value)
		} else {
			v.SetDefault(key, value)
		}
	}
}

func GetDefaultInt(name string) int {
	var (
		value int
		ok    bool
	)

	if value, ok = defaults[name].(int); !ok {
		return 0
	}

	return value
}

func GetDefaultString(name string) string {
	var (
		value string
		ok    bool
	)

	if value, ok = defaults[name].(string); !ok {
		return ""
	}

	return value
}
