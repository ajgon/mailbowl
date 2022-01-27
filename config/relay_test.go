package config_test

import (
	"testing"

	"github.com/ajgon/mailbowl/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestValidRelayMarshalFromObjects(t *testing.T) {
	t.Parallel()

	viperConfig := viper.New()
	// address takes precedence if set
	viperConfig.Set("relay.outgoing_server.address", "tls://192.168.42.1:10025")
	viperConfig.Set("relay.outgoing_server.host", "192.168.10.10")
	viperConfig.Set("relay.outgoing_server.port", 10465)
	viperConfig.Set("relay.outgoing_server.connection_type", "plain")
	viperConfig.Set("relay.outgoing_server.auth_method", "plain")
	viperConfig.Set("relay.outgoing_server.from_email", "server@example.local")
	viperConfig.Set("relay.outgoing_server.password", "secret")
	viperConfig.Set("relay.outgoing_server.username", "user@example.local")
	viperConfig.Set("relay.outgoing_server.verify_tls", true)
	conf, err := InitConfig(viperConfig)
	assert.NoError(t, err)

	assert.Equal(t, "192.168.42.1", conf.Relay.OutgoingServer.Host)
	assert.Equal(t, 10025, conf.Relay.OutgoingServer.Port)
	assert.Equal(t, config.ConnectionTLS, conf.Relay.OutgoingServer.ConnectionType)
	assert.Equal(t, config.AuthPlain, conf.Relay.OutgoingServer.AuthMethod)
	assert.Equal(t, "server@example.local", conf.Relay.OutgoingServer.FromEmail)
	assert.Equal(t, "secret", conf.Relay.OutgoingServer.Password)
	assert.Equal(t, "user@example.local", conf.Relay.OutgoingServer.Username)
	assert.True(t, conf.Relay.OutgoingServer.VerifyTLS)
}

func TestValidRelayMarshalFromYAML(t *testing.T) {
	t.Parallel()

	yamlExample := `---
relay:
  outgoing_server:
    address: plain://192.168.45.1:20025
    host: 192.168.10.10
    port: 10465
    connection_type: tls
    auth_method: plain
    from_email: yaml@example.local
    password: secretyaml
    username: useryaml@example.local
    verify_tls: true
`

	viperConfig := viper.New()

	conf, err := InitConfig(viperConfig, yamlExample)
	assert.NoError(t, err)

	assert.Equal(t, "192.168.45.1", conf.Relay.OutgoingServer.Host)
	assert.Equal(t, 20025, conf.Relay.OutgoingServer.Port)
	assert.Equal(t, config.ConnectionPlain, conf.Relay.OutgoingServer.ConnectionType)
	assert.Equal(t, config.AuthPlain, conf.Relay.OutgoingServer.AuthMethod)
	assert.Equal(t, "yaml@example.local", conf.Relay.OutgoingServer.FromEmail)
	assert.Equal(t, "secretyaml", conf.Relay.OutgoingServer.Password)
	assert.Equal(t, "useryaml@example.local", conf.Relay.OutgoingServer.Username)
	assert.True(t, conf.Relay.OutgoingServer.VerifyTLS)
}

func TestValidRelayMarshalFromENV(t *testing.T) {
	t.Setenv("RELAY_OUTGOING_SERVER_ADDRESS", "starttls://192.168.80.1:30025")
	t.Setenv("RELAY_OUTGOING_SERVER_HOST", "192.168.10.10")
	t.Setenv("RELAY_OUTGOING_SERVER_PORT", "10465")
	t.Setenv("RELAY_OUTGOING_SERVER_CONNECTION_TYPE", "plain")
	t.Setenv("RELAY_OUTGOING_SERVER_AUTH_METHOD", "plain")
	t.Setenv("RELAY_OUTGOING_SERVER_FROM_EMAIL", "env@example.local")
	t.Setenv("RELAY_OUTGOING_SERVER_PASSWORD", "secretenv")
	t.Setenv("RELAY_OUTGOING_SERVER_USERNAME", "userenv@example.local")
	t.Setenv("RELAY_OUTGOING_SERVER_VERIFY_TLS", "true")

	viperConfig := viper.New()
	conf, err := InitConfig(viperConfig)
	assert.NoError(t, err)

	assert.Equal(t, "192.168.80.1", conf.Relay.OutgoingServer.Host)
	assert.Equal(t, 30025, conf.Relay.OutgoingServer.Port)
	assert.Equal(t, config.ConnectionStartTLS, conf.Relay.OutgoingServer.ConnectionType)
	assert.Equal(t, config.AuthPlain, conf.Relay.OutgoingServer.AuthMethod)
	assert.Equal(t, "env@example.local", conf.Relay.OutgoingServer.FromEmail)
	assert.Equal(t, "secretenv", conf.Relay.OutgoingServer.Password)
	assert.Equal(t, "userenv@example.local", conf.Relay.OutgoingServer.Username)
	assert.True(t, conf.Relay.OutgoingServer.VerifyTLS)
}

func TestValidRelayHostPort(t *testing.T) {
	t.Parallel()

	viperConfig := viper.New()
	viperConfig.Set("relay.outgoing_server.address", "")
	viperConfig.Set("relay.outgoing_server.host", "192.168.10.10")
	viperConfig.Set("relay.outgoing_server.port", 10465)
	viperConfig.Set("relay.outgoing_server.connection_type", "tls")
	conf, err := InitConfig(viperConfig)
	assert.NoError(t, err)

	assert.Equal(t, "192.168.10.10", conf.Relay.OutgoingServer.Host)
	assert.Equal(t, 10465, conf.Relay.OutgoingServer.Port)
	assert.Equal(t, config.ConnectionTLS, conf.Relay.OutgoingServer.ConnectionType)
}

func TestInvalidAuthMethod(t *testing.T) {
	t.Parallel()

	viperConfig := viper.New()
	viperConfig.Set("relay.outgoing_server.auth_method", "wrong")
	_, err := InitConfig(viperConfig)

	assert.EqualError(
		t, err, "error unmarshaling config: 1 error(s) decoding:\n\n* error decoding 'Relay': invalid auth method",
	)
}

func TestInvalidAddress(t *testing.T) {
	t.Parallel()

	viperConfig := viper.New()
	viperConfig.Set("relay.outgoing_server.address", "wrong-address")
	_, err := InitConfig(viperConfig)

	assert.EqualError(
		t, err, "error unmarshaling config: 1 error(s) decoding:\n\n"+
			"* error decoding 'Relay': invalid address: missing port in address",
	)
}
