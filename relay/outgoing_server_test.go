package relay_test

import (
	"context"
	"crypto/tls"
	"fmt"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/ajgon/mailbowl/listener/smtp"
	"github.com/ajgon/mailbowl/relay"
	"github.com/chrj/smtpd"
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

type SMTPTestServer struct {
	Port       int
	Sender     string
	Recipients []string
	Message    string
	TLS        *tls.Config
}

func NewSMTPTestServer() *SMTPTestServer {
	smtpTLS, _ := smtp.NewTLS(
		`
-----BEGIN PRIVATE KEY-----
MC4CAQAwBQYDK2VwBCIEIMsO5gzH0FPxq8AkEgyBoJEBvAxOcCnKENdzYTWbwe6Q
-----END PRIVATE KEY-----
		`, `
-----BEGIN CERTIFICATE-----
MIH+MIGxoAMCAQICFHAn6jIQ9qZxySDKoL/oQXgyjL7YMAUGAytlcDAUMRIwEAYD
VQQDDAlsb2NhbGhvc3QwIBcNMjIwMTEyMjEwOTQ4WhgPMjEyMTEyMTkyMTA5NDha
MBQxEjAQBgNVBAMMCWxvY2FsaG9zdDAqMAUGAytlcAMhAMHNBNUlEKkIgCGnWMIF
m6f8MOg/ZQOOXQEgmdUyAehqoxMwETAPBgNVHREECDAGhwR/AAABMAUGAytlcANB
AGo3n53h0jGSFiTMGwBYnrV/69aPjxdB/LGr4p0/v355GVqZXyZ7idCpSuCCiYmk
DQ2hhzbuPuECiTPOUYSO5wI=
-----END CERTIFICATE-----
		`,
		"", "", false,
	)

	return &SMTPTestServer{
		Port: randomPort(),
		TLS:  smtpTLS.Config,
	}
}

func (sts *SMTPTestServer) Serve(ctx context.Context, connectionType string) {
	var listener net.Listener

	smtpdServer := &smtpd.Server{
		Hostname: "127.0.0.1",

		Handler: func(_peer smtpd.Peer, envelope smtpd.Envelope) error {
			sts.Sender = envelope.Sender
			sts.Recipients = envelope.Recipients
			sts.Message = string(envelope.Data)

			return nil
		},
	}

	switch connectionType {
	case "plain":
		listener, _ = net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", sts.Port))
	case "starttls":
		listener, _ = net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", sts.Port))
		smtpdServer.TLSConfig = sts.TLS
	case "tls":
		listener, _ = tls.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", sts.Port), sts.TLS)
		smtpdServer.TLSConfig = sts.TLS
		smtpdServer.ForceTLS = false
	}

	defer listener.Close()

	go func() {
		_ = smtpdServer.Serve(listener)
	}()

	<-ctx.Done()

	_ = smtpdServer.Shutdown(true)
}

func TestSendPlainWithDefaultFrom(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())

	testSMTPServer := NewSMTPTestServer()
	go testSMTPServer.Serve(ctx, "plain")

	outgoingServer, err := relay.NewOutgoingServer(
		fmt.Sprintf("127.0.0.1:%d", testSMTPServer.Port), "", "plain", "override@example.local", "password", "user", false,
	)
	assert.NoError(t, err)
	time.Sleep(100 * time.Millisecond) // allow server to start

	err = outgoingServer.Send("from@example.local", []string{"to@example.local"}, []byte("test message"))
	assert.NoError(t, err)

	cancel()

	assert.Equal(t, "override@example.local", testSMTPServer.Sender)
	assert.Equal(t, []string{"to@example.local"}, testSMTPServer.Recipients)
	assert.Equal(t, "test message\n", testSMTPServer.Message)
}

func TestSendStartTLSWithCustomFrom(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())

	testSMTPServer := NewSMTPTestServer()
	go testSMTPServer.Serve(ctx, "starttls")

	outgoingServer, err := relay.NewOutgoingServer(
		fmt.Sprintf("127.0.0.1:%d", testSMTPServer.Port), "", "starttls", "", "password", "user", false,
	)
	assert.NoError(t, err)
	time.Sleep(100 * time.Millisecond) // allow server to start

	err = outgoingServer.Send("from@example.local", []string{"to@example.local"}, []byte("test message"))
	assert.NoError(t, err)

	cancel()

	assert.Equal(t, "from@example.local", testSMTPServer.Sender)
	assert.Equal(t, []string{"to@example.local"}, testSMTPServer.Recipients)
	assert.Equal(t, "test message\n", testSMTPServer.Message)
}

func TestSendTLSWithCustomFrom(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())

	testSMTPServer := NewSMTPTestServer()
	go testSMTPServer.Serve(ctx, "tls")

	outgoingServer, err := relay.NewOutgoingServer(
		fmt.Sprintf("127.0.0.1:%d", testSMTPServer.Port), "", "tls", "", "password", "user", false,
	)
	assert.NoError(t, err)
	time.Sleep(100 * time.Millisecond) // allow server to start

	err = outgoingServer.Send("from@example.local", []string{"to@example.local"}, []byte("test message"))
	assert.NoError(t, err)

	cancel()

	assert.Equal(t, "from@example.local", testSMTPServer.Sender)
	assert.Equal(t, []string{"to@example.local"}, testSMTPServer.Recipients)
	assert.Equal(t, "test message\n", testSMTPServer.Message)
}
