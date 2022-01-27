package config_test

import (
	"testing"
	"time"

	"github.com/ajgon/mailbowl/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestDefaults(t *testing.T) {
	t.Parallel()

	viperConfig := viper.New()
	conf, err := InitConfig(viperConfig)
	assert.NoError(t, err)

	assert.False(t, conf.Log.Color)
	assert.Equal(t, config.Console, conf.Log.Format)
	assert.Equal(t, zapcore.WarnLevel, conf.Log.Level)
	assert.Equal(t, zapcore.ErrorLevel, conf.Log.StacktraceLevel)
	assert.Equal(t, "", conf.Relay.OutgoingServer.Host)
	assert.Equal(t, 0, conf.Relay.OutgoingServer.Port)
	assert.Equal(t, config.ConnectionTLS, conf.Relay.OutgoingServer.ConnectionType)
	assert.Equal(t, config.AuthPlain, conf.Relay.OutgoingServer.AuthMethod)
	assert.Equal(t, "", conf.Relay.OutgoingServer.FromEmail)
	assert.Equal(t, "", conf.Relay.OutgoingServer.Password)
	assert.Equal(t, "", conf.Relay.OutgoingServer.Username)
	assert.True(t, conf.Relay.OutgoingServer.VerifyTLS)
	assert.False(t, conf.SMTP.Auth.Enabled)
	assert.Equal(t, []config.SMTPAuthUser{}, conf.SMTP.Auth.Users)
	assert.Equal(t, "", conf.SMTP.Hostname)
	assert.Equal(t, 100, conf.SMTP.Limit.Connections)
	assert.Equal(t, 26214400, conf.SMTP.Limit.MessageSize)
	assert.Equal(t, 100, conf.SMTP.Limit.Recipients)
	assert.Equal(t, []config.SMTPListen{}, conf.SMTP.Listen)
	assert.Equal(t, 60*time.Second, conf.SMTP.Timeout.Read)
	assert.Equal(t, 60*time.Second, conf.SMTP.Timeout.Write)
	assert.Equal(t, 5*time.Minute, conf.SMTP.Timeout.Data)
	assert.Equal(t, "", conf.SMTP.TLS.Key)
	assert.Equal(t, "", conf.SMTP.TLS.Certificate)
	assert.Equal(t, "", conf.SMTP.TLS.KeyFile)
	assert.Equal(t, "", conf.SMTP.TLS.CertificateFile)
	assert.True(t, conf.SMTP.TLS.ForceForStartTLS)
	assert.Equal(t, []string{}, conf.SMTP.Whitelist)
}
