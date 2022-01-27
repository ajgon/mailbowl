package config_test

import (
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestValidSMTPMarshalFromObjects(t *testing.T) {
	t.Parallel()

	viperConfig := viper.New()
	viperConfig.Set("smtp.auth.enabled", true)
	viperConfig.Set(
		"smtp.auth.users", []interface{}{
			map[interface{}]interface{}{"email": "test@example.local", "password_hash": "hash"},
			map[interface{}]interface{}{"email": "other@example.local", "password_hash": "pass"},
		})
	viperConfig.Set("smtp.hostname", "object.local")
	viperConfig.Set("smtp.limit.connections", 1024)
	viperConfig.Set("smtp.limit.message_size", 2048)
	viperConfig.Set("smtp.limit.recipients", 3072)
	viperConfig.Set("smtp.listen", []string{"plain://0.0.0.0:25", "tls://192.168.0.0:465"})
	viperConfig.Set("smtp.timeout.read", "10s")
	viperConfig.Set("smtp.timeout.write", "20s")
	viperConfig.Set("smtp.timeout.data", "30s")
	viperConfig.Set("smtp.tls.key", "object-tls-key")
	viperConfig.Set("smtp.tls.key_file", "object-tls-keyfile")
	viperConfig.Set("smtp.tls.certificate", "object-tls-certificate")
	viperConfig.Set("smtp.tls.certificate_file", "object-tls-certificatefile")
	viperConfig.Set("smtp.tls.force_for_starttls", false)
	viperConfig.Set("smtp.whitelist", []string{"192.168.10.0/24", "192.168.20.0/24"})
	conf, err := InitConfig(viperConfig)
	assert.NoError(t, err)

	assert.True(t, conf.SMTP.Auth.Enabled)
	assert.Equal(t, "test@example.local", conf.SMTP.Auth.Users[0].Email)
	assert.Equal(t, "hash", conf.SMTP.Auth.Users[0].PasswordHash)
	assert.Equal(t, "other@example.local", conf.SMTP.Auth.Users[1].Email)
	assert.Equal(t, "pass", conf.SMTP.Auth.Users[1].PasswordHash)
	assert.Equal(t, 2, len(conf.SMTP.Auth.Users))
	assert.Equal(t, "object.local", conf.SMTP.Hostname)
	assert.Equal(t, 1024, conf.SMTP.Limit.Connections)
	assert.Equal(t, 2048, conf.SMTP.Limit.MessageSize)
	assert.Equal(t, 3072, conf.SMTP.Limit.Recipients)
	assert.Equal(t, "plain", conf.SMTP.Listen[0].Proto)
	assert.Equal(t, "0.0.0.0", conf.SMTP.Listen[0].Host)
	assert.Equal(t, "25", conf.SMTP.Listen[0].Port)
	assert.Equal(t, "tls", conf.SMTP.Listen[1].Proto)
	assert.Equal(t, "192.168.0.0", conf.SMTP.Listen[1].Host)
	assert.Equal(t, "465", conf.SMTP.Listen[1].Port)
	assert.Equal(t, 2, len(conf.SMTP.Listen))
	assert.Equal(t, 10*time.Second, conf.SMTP.Timeout.Read)
	assert.Equal(t, 20*time.Second, conf.SMTP.Timeout.Write)
	assert.Equal(t, 30*time.Second, conf.SMTP.Timeout.Data)
	assert.Equal(t, "object-tls-key", conf.SMTP.TLS.Key)
	assert.Equal(t, "object-tls-keyfile", conf.SMTP.TLS.KeyFile)
	assert.Equal(t, "object-tls-certificate", conf.SMTP.TLS.Certificate)
	assert.Equal(t, "object-tls-certificatefile", conf.SMTP.TLS.CertificateFile)
	assert.False(t, conf.SMTP.TLS.ForceForStartTLS)
	assert.Equal(t, []string{"192.168.10.0/24", "192.168.20.0/24"}, conf.SMTP.Whitelist)
}

func TestValidSMTPMarshalFromYAML(t *testing.T) {
	t.Parallel()

	yamlExample := `---
smtp:
  auth:
    enabled: true
    users:
      - email: test@yaml.local
        password_hash: yamlhash
      - email: other@yaml.local
        password_hash: yamlpass
  hostname: yaml.local
  limit:
    connections: 512
    message_size: 384
    recipients: 256
  listen:
    - tls://10.0.0.0:1465
    - starttls://172.12.0.0:1587
  timeout:
    read: 10m
    write: 20m
    data: 30m
  tls:
    key: yaml-tls-key
    key_file: yaml-tls-keyfile
    certificate: yaml-tls-certificate
    certificate_file: yaml-tls-certificatefile
    force_for_starttls: false
  whitelist:
    - 172.16.0.0/16
    - 172.24.0.0/16
`

	viperConfig := viper.New()

	conf, err := InitConfig(viperConfig, yamlExample)
	assert.NoError(t, err)

	assert.True(t, conf.SMTP.Auth.Enabled)
	assert.Equal(t, "test@yaml.local", conf.SMTP.Auth.Users[0].Email)
	assert.Equal(t, "yamlhash", conf.SMTP.Auth.Users[0].PasswordHash)
	assert.Equal(t, "other@yaml.local", conf.SMTP.Auth.Users[1].Email)
	assert.Equal(t, "yamlpass", conf.SMTP.Auth.Users[1].PasswordHash)
	assert.Equal(t, 2, len(conf.SMTP.Auth.Users))
	assert.Equal(t, "yaml.local", conf.SMTP.Hostname)
	assert.Equal(t, 512, conf.SMTP.Limit.Connections)
	assert.Equal(t, 384, conf.SMTP.Limit.MessageSize)
	assert.Equal(t, 256, conf.SMTP.Limit.Recipients)
	assert.Equal(t, "tls", conf.SMTP.Listen[0].Proto)
	assert.Equal(t, "10.0.0.0", conf.SMTP.Listen[0].Host)
	assert.Equal(t, "1465", conf.SMTP.Listen[0].Port)
	assert.Equal(t, "starttls", conf.SMTP.Listen[1].Proto)
	assert.Equal(t, "172.12.0.0", conf.SMTP.Listen[1].Host)
	assert.Equal(t, "1587", conf.SMTP.Listen[1].Port)
	assert.Equal(t, 2, len(conf.SMTP.Listen))
	assert.Equal(t, 10*time.Minute, conf.SMTP.Timeout.Read)
	assert.Equal(t, 20*time.Minute, conf.SMTP.Timeout.Write)
	assert.Equal(t, 30*time.Minute, conf.SMTP.Timeout.Data)
	assert.Equal(t, "yaml-tls-key", conf.SMTP.TLS.Key)
	assert.Equal(t, "yaml-tls-keyfile", conf.SMTP.TLS.KeyFile)
	assert.Equal(t, "yaml-tls-certificate", conf.SMTP.TLS.Certificate)
	assert.Equal(t, "yaml-tls-certificatefile", conf.SMTP.TLS.CertificateFile)
	assert.False(t, conf.SMTP.TLS.ForceForStartTLS)
	assert.Equal(t, []string{"172.16.0.0/16", "172.24.0.0/16"}, conf.SMTP.Whitelist)
}

func TestValidSMTPMarshalFromENV(t *testing.T) {
	t.Setenv("SMTP_AUTH_ENABLED", "true")
	t.Setenv("SMTP_AUTH_USERS", "test@env.local:envhash ot:her@env.local:envpass")
	t.Setenv("SMTP_HOSTNAME", "env.local")
	t.Setenv("SMTP_LIMIT_CONNECTIONS", "16")
	t.Setenv("SMTP_LIMIT_MESSAGE_SIZE", "12")
	t.Setenv("SMTP_LIMIT_RECIPIENTS", "8")
	t.Setenv("SMTP_LISTEN", "starttls://192.168.42.0:2587 plain://172.16.0.0:2025")
	t.Setenv("SMTP_TIMEOUT_READ", "10h")
	t.Setenv("SMTP_TIMEOUT_WRITE", "20h")
	t.Setenv("SMTP_TIMEOUT_DATA", "30h")
	t.Setenv("SMTP_TLS_KEY", "env-tls-key")
	t.Setenv("SMTP_TLS_KEY_FILE", "env-tls-keyfile")
	t.Setenv("SMTP_TLS_CERTIFICATE", "env-tls-certificate")
	t.Setenv("SMTP_TLS_CERTIFICATE_FILE", "env-tls-certificatefile")
	t.Setenv("SMTP_TLS_FORCE_FOR_STARTTLS", "0")
	t.Setenv("SMTP_WHITELIST", "10.10.10.0/8 10.20.0.0/16")

	viperConfig := viper.New()
	conf, err := InitConfig(viperConfig)
	assert.NoError(t, err)

	assert.True(t, conf.SMTP.Auth.Enabled)
	assert.Equal(t, "test@env.local", conf.SMTP.Auth.Users[0].Email)
	assert.Equal(t, "envhash", conf.SMTP.Auth.Users[0].PasswordHash)
	assert.Equal(t, "ot:her@env.local", conf.SMTP.Auth.Users[1].Email)
	assert.Equal(t, "envpass", conf.SMTP.Auth.Users[1].PasswordHash)
	assert.Equal(t, 2, len(conf.SMTP.Auth.Users))
	assert.Equal(t, "env.local", conf.SMTP.Hostname)
	assert.Equal(t, 16, conf.SMTP.Limit.Connections)
	assert.Equal(t, 12, conf.SMTP.Limit.MessageSize)
	assert.Equal(t, 8, conf.SMTP.Limit.Recipients)
	assert.Equal(t, "starttls", conf.SMTP.Listen[0].Proto)
	assert.Equal(t, "192.168.42.0", conf.SMTP.Listen[0].Host)
	assert.Equal(t, "2587", conf.SMTP.Listen[0].Port)
	assert.Equal(t, "plain", conf.SMTP.Listen[1].Proto)
	assert.Equal(t, "172.16.0.0", conf.SMTP.Listen[1].Host)
	assert.Equal(t, "2025", conf.SMTP.Listen[1].Port)
	assert.Equal(t, 2, len(conf.SMTP.Listen))
	assert.Equal(t, 10*time.Hour, conf.SMTP.Timeout.Read)
	assert.Equal(t, 20*time.Hour, conf.SMTP.Timeout.Write)
	assert.Equal(t, 30*time.Hour, conf.SMTP.Timeout.Data)
	assert.Equal(t, "env-tls-key", conf.SMTP.TLS.Key)
	assert.Equal(t, "env-tls-keyfile", conf.SMTP.TLS.KeyFile)
	assert.Equal(t, "env-tls-certificate", conf.SMTP.TLS.Certificate)
	assert.Equal(t, "env-tls-certificatefile", conf.SMTP.TLS.CertificateFile)
	assert.False(t, conf.SMTP.TLS.ForceForStartTLS)
	assert.Equal(t, []string{"10.10.10.0/8", "10.20.0.0/16"}, conf.SMTP.Whitelist)
}
