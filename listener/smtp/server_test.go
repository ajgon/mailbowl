package smtp_test

import (
	"crypto/tls"
	"fmt"
	"math/rand"
	netsmtp "net/smtp"
	"testing"
	"time"

	"github.com/ajgon/mailbowl/listener/smtp"
	"github.com/stretchr/testify/assert"
)

type RandomPorts struct {
	Ports map[int]bool
}

var randomPorts = &RandomPorts{Ports: map[int]bool{0: true}} //nolint:gochecknoglobals

func randomPort() int {
	var port int

	rand.Seed(time.Now().UnixNano())

	for randomPorts.Ports[port] {
		port = rand.Intn(64512) + 1024 //nolint:gosec
	}

	randomPorts.Ports[port] = true

	return port
}

func newTestServer(t *testing.T, proto, cidr string, includeTLSCertificate, authEnabled bool) (*smtp.Server, string) {
	t.Helper()

	var tls *smtp.TLS

	host := fmt.Sprintf("127.0.0.1:%d", randomPort())
	url := fmt.Sprintf("%s://%s", proto, host)

	tlsKeyExample := `
-----BEGIN PRIVATE KEY-----
MC4CAQAwBQYDK2VwBCIEIMsO5gzH0FPxq8AkEgyBoJEBvAxOcCnKENdzYTWbwe6Q
-----END PRIVATE KEY-----
	`
	tlsCertificateExample := `
-----BEGIN CERTIFICATE-----
MIH+MIGxoAMCAQICFHAn6jIQ9qZxySDKoL/oQXgyjL7YMAUGAytlcDAUMRIwEAYD
VQQDDAlsb2NhbGhvc3QwIBcNMjIwMTEyMjEwOTQ4WhgPMjEyMTEyMTkyMTA5NDha
MBQxEjAQBgNVBAMMCWxvY2FsaG9zdDAqMAUGAytlcAMhAMHNBNUlEKkIgCGnWMIF
m6f8MOg/ZQOOXQEgmdUyAehqoxMwETAPBgNVHREECDAGhwR/AAABMAUGAytlcANB
AGo3n53h0jGSFiTMGwBYnrV/69aPjxdB/LGr4p0/v355GVqZXyZ7idCpSuCCiYmk
DQ2hhzbuPuECiTPOUYSO5wI=
-----END CERTIFICATE-----
	`

	auth := smtp.NewAuth(
		authEnabled,
		[]*smtp.AuthUser{
			smtp.NewAuthUser("test@example.local", "$2a$10$BoHLl7lps2ZhB.B5h3Zqau.p4VAQR7BVjmWTC7nEbDAY9Kp4LjNrW"),
		},
	)
	limit := smtp.NewLimit(100, 200, 300)
	uri, err := smtp.NewURI(url)
	assert.NoError(t, err)

	timeout := smtp.NewTimeout("10s", "20s", "30s")

	if includeTLSCertificate {
		tls, err = smtp.NewTLS(tlsKeyExample, tlsCertificateExample, "", "", true)
		assert.NoError(t, err)
	} else {
		tls, err = smtp.NewTLS("", "", "", "", true)
		assert.ErrorIs(t, err, smtp.ErrTLSNotConfigured)
	}

	whitelist := smtp.NewWhitelist([]string{cidr})

	return smtp.NewServer(auth, "hostname", limit, timeout, tls, whitelist, uri, nil), host
}

func TestBuildServer(t *testing.T) {
	t.Parallel()

	server, _ := newTestServer(t, "plain", "0.0.0.0/0", false, false)
	err := server.Build()
	assert.NoError(t, err)

	defer server.Listener.Close()

	assert.Equal(t, server.Hostname, "hostname")
	assert.Equal(t, server.SMTPD.Hostname, "hostname")
	assert.Equal(t, server.Limit.Connections, 100)
	assert.Equal(t, server.SMTPD.MaxConnections, 100)
	assert.Equal(t, server.Limit.MessageSize, 200)
	assert.Equal(t, server.SMTPD.MaxMessageSize, 200)
	assert.Equal(t, server.Limit.Recipients, 300)
	assert.Equal(t, server.SMTPD.MaxRecipients, 300)
	assert.Equal(t, server.Timeout.Read, 10*time.Second)
	assert.Equal(t, server.SMTPD.WriteTimeout, 20*time.Second)
	assert.Equal(t, server.Timeout.Write, 20*time.Second)
	assert.Equal(t, server.SMTPD.DataTimeout, 30*time.Second)
	assert.Equal(t, server.Timeout.Data, 30*time.Second)
	assert.Equal(t, server.SMTPD.ReadTimeout, 10*time.Second)
}

func TestBuildPlainServer(t *testing.T) {
	t.Parallel()

	server, host := newTestServer(t, "plain", "0.0.0.0/0", false, false)
	err := server.Build()
	assert.NoError(t, err)

	defer server.Listener.Close()

	assert.Nil(t, server.TLS)
	assert.False(t, server.SMTPD.ForceTLS)
	assert.Nil(t, server.SMTPD.TLSConfig)
	assert.Equal(t, server.Listener.Addr().String(), host)
}

func TestBuildTlsServer(t *testing.T) {
	t.Parallel()

	server, host := newTestServer(t, "tls", "0.0.0.0/0", true, false)
	err := server.Build()
	assert.NoError(t, err)

	defer server.Listener.Close()

	assert.True(t, server.TLS.ForceForStartTLS)
	assert.False(t, server.SMTPD.ForceTLS)
	assert.NotNil(t, server.SMTPD.TLSConfig)
	assert.Equal(t, server.Listener.Addr().String(), host)
}

func TestBuildTlsServerWithoutCertificates(t *testing.T) {
	t.Parallel()

	server, _ := newTestServer(t, "tls", "0.0.0.0/0", false, false)
	err := server.Build()
	assert.Error(t, err)

	assert.Nil(t, server.TLS)
	assert.False(t, server.SMTPD.ForceTLS)
	assert.Nil(t, server.SMTPD.TLSConfig)
	assert.Nil(t, server.Listener)
}

func TestBuildStartTlsServer(t *testing.T) {
	t.Parallel()

	server, host := newTestServer(t, "starttls", "0.0.0.0/0", true, false)
	err := server.Build()
	assert.NoError(t, err)

	defer server.Listener.Close()

	assert.True(t, server.TLS.ForceForStartTLS)
	assert.True(t, server.SMTPD.ForceTLS)
	assert.NotNil(t, server.SMTPD.TLSConfig)
	assert.Equal(t, server.Listener.Addr().String(), host)
}

func TestBuildStartTlsServerWithoutCertificates(t *testing.T) {
	t.Parallel()

	server, _ := newTestServer(t, "starttls", "0.0.0.0/0", false, false)
	err := server.Build()
	assert.Error(t, err)

	assert.Nil(t, server.TLS)
	assert.False(t, server.SMTPD.ForceTLS)
	assert.Nil(t, server.SMTPD.TLSConfig)
	assert.Nil(t, server.Listener)
}

func TestIPWhitelistDenied(t *testing.T) {
	t.Parallel()

	server, host := newTestServer(t, "plain", "10.0.0.0/8", false, false)
	err := server.Build()
	assert.NoError(t, err)

	go server.Start()
	defer server.Shutdown() //nolint: errcheck

	err = netsmtp.SendMail(
		host, nil, "sender@example.local", []string{"receiver@example.local"}, []byte("Subject: Test"),
	)

	assert.Equal(t, err.Error(), "421 Denied")
}

func TestIPWhitelisAllowed(t *testing.T) {
	t.Parallel()

	server, host := newTestServer(t, "plain", "127.0.0.1/32", false, false)
	err := server.Build()
	assert.NoError(t, err)

	go server.Start()
	defer server.Shutdown() //nolint: errcheck

	err = netsmtp.SendMail(
		host, nil, "sender@example.local", []string{"receiver@example.local"}, []byte("Subject: Test"),
	)

	assert.NoError(t, err)
}

func TestNoAuthenticationForPlainServer(t *testing.T) {
	t.Parallel()

	t.Skip("TODO")
}

func TestAuthenticationDeniedForMissingUser(t *testing.T) {
	t.Parallel()

	server, host := newTestServer(t, "tls", "127.0.0.1/32", true, true)
	err := server.Build()
	assert.NoError(t, err)

	go server.Start()
	defer server.Shutdown() //nolint: errcheck

	tlsconfig := &tls.Config{
		InsecureSkipVerify: true, //nolint: gosec
		ServerName:         "localhost",
	}

	tlsConnection, err := tls.Dial("tcp", host, tlsconfig)
	assert.NoError(t, err)

	defer tlsConnection.Close()

	client, err := netsmtp.NewClient(tlsConnection, "127.0.0.1")
	assert.NoError(t, err)

	defer client.Quit() //nolint: errcheck

	err = client.Mail("sender@example.local")
	assert.Equal(t, err.Error(), "530 Authentication Required.")
}

func TestAuthenticationDeniedForWrongCredentials(t *testing.T) {
	t.Parallel()

	server, host := newTestServer(t, "tls", "127.0.0.1/32", true, true)
	err := server.Build()
	assert.NoError(t, err)

	go server.Start()
	defer server.Shutdown() //nolint: errcheck

	tlsconfig := &tls.Config{
		InsecureSkipVerify: true, //nolint: gosec
		ServerName:         "localhost",
	}

	tlsConnection, err := tls.Dial("tcp", host, tlsconfig)
	assert.NoError(t, err)

	defer tlsConnection.Close()

	client, err := netsmtp.NewClient(tlsConnection, "127.0.0.1")
	assert.NoError(t, err)

	defer client.Quit() //nolint: errcheck

	auth := netsmtp.PlainAuth("", "test@example.local", "wrongpassword", "127.0.0.1")
	err = client.Auth(auth)
	assert.Equal(t, err.Error(), "535 Authentication credentials invalid")
}

func TestAuthenticationValid(t *testing.T) {
	t.Parallel()

	server, host := newTestServer(t, "tls", "127.0.0.1/32", true, true)
	err := server.Build()
	assert.NoError(t, err)

	go server.Start()
	defer server.Shutdown() //nolint: errcheck

	tlsconfig := &tls.Config{
		InsecureSkipVerify: true, //nolint: gosec
		ServerName:         "localhost",
	}

	tlsConnection, err := tls.Dial("tcp", host, tlsconfig)
	assert.NoError(t, err)

	defer tlsConnection.Close()

	client, err := netsmtp.NewClient(tlsConnection, "127.0.0.1")
	assert.NoError(t, err)

	defer client.Quit() //nolint: errcheck

	auth := netsmtp.PlainAuth("", "test@example.local", "test", "127.0.0.1")
	err = client.Auth(auth)
	assert.NoError(t, err)
}
