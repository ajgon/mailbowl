package relay

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	"github.com/ajgon/mailbowl/config"
)

const (
	ConnectionPlain    = "plain"
	ConnectionStartTLS = "starttls"
	ConnectionTLS      = "tls"
)

type OutgoingServer struct {
	AuthMethod     config.RelayAuthMethod
	ConnectionType config.RelayConnectionType
	FromEmail      string
	Host           string
	Password       string
	Port           int
	Username       string
	VerifyTLS      bool
}

func NewOutgoingServer(conf config.RelayOutgoingServer) (*OutgoingServer, error) {
	return &OutgoingServer{
		AuthMethod:     conf.AuthMethod,
		ConnectionType: conf.ConnectionType,
		FromEmail:      conf.FromEmail,
		Host:           conf.Host,
		Password:       conf.Password,
		Port:           conf.Port,
		Username:       conf.Username,
		VerifyTLS:      conf.VerifyTLS,
	}, nil
}

func (ros *OutgoingServer) Send(from string, recipients []string, message []byte) error {
	if ros.FromEmail != "" {
		from = ros.FromEmail
	}

	if ros.ConnectionType == config.ConnectionPlain {
		return ros.sendPlain(from, recipients, message)
	}

	return ros.sendTLS(from, recipients, message)
}

func (ros *OutgoingServer) sendTLS(from string, recipients []string, message []byte) error { // nolint:cyclop
	var (
		client *smtp.Client
		err    error
	)

	auth := ros.buildAuth()

	client, err = ros.buildTLSClient()
	if err != nil {
		return fmt.Errorf("outgoing smtp error: %w", err)
	}

	// Auth
	if auth != nil {
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("outgoing smtp error: %w", err)
		}
	}

	// To && From
	if err = client.Mail(from); err != nil {
		return fmt.Errorf("outgoing smtp error: %w", err)
	}

	for _, rcpt := range recipients {
		if err = client.Rcpt(rcpt); err != nil {
			return fmt.Errorf("outgoing smtp error: %w", err)
		}
	}

	// Data
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("outgoing smtp error: %w", err)
	}

	_, err = writer.Write(message)
	if err != nil {
		return fmt.Errorf("outgoing smtp error: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("outgoing smtp error: %w", err)
	}

	err = client.Quit()
	if err != nil {
		return fmt.Errorf("outgoing smtp error: %w", err)
	}

	return nil
}

func (ros *OutgoingServer) buildTLSClient() (client *smtp.Client, err error) {
	var conn *tls.Conn

	tlsconfig := &tls.Config{
		InsecureSkipVerify: !ros.VerifyTLS, //nolint:gosec
		ServerName:         ros.Host,
	}

	if ros.ConnectionType == config.ConnectionStartTLS {
		client, err = smtp.Dial(fmt.Sprintf("%s:%d", ros.Host, ros.Port))
		if err != nil {
			return nil, fmt.Errorf("outgoing starttls error: %w", err)
		}

		err = client.StartTLS(tlsconfig)
		if err != nil {
			return nil, fmt.Errorf("outgoing starttls error: %w", err)
		}

		return client, nil
	}

	conn, err = tls.Dial("tcp", fmt.Sprintf("%s:%d", ros.Host, ros.Port), tlsconfig)
	if err != nil {
		return nil, fmt.Errorf("outgoing tls error: %w", err)
	}

	client, err = smtp.NewClient(conn, ros.Host)
	if err != nil {
		return nil, fmt.Errorf("outgoing smtp error: %w", err)
	}

	return client, nil
}

func (ros *OutgoingServer) sendPlain(from string, recipients []string, message []byte) error {
	auth := ros.buildAuth()

	err := smtp.SendMail(fmt.Sprintf("%s:%d", ros.Host, ros.Port), auth, from, recipients, message)
	if err != nil {
		return fmt.Errorf("error sending email via outgoing server: %w", err)
	}

	return nil
}

func (ros *OutgoingServer) buildAuth() smtp.Auth {
	if ros.AuthMethod == config.AuthPlain {
		return smtp.PlainAuth("", ros.Username, ros.Password, ros.Host)
	}

	if ros.AuthMethod == config.AuthCramMD5 {
		return smtp.CRAMMD5Auth(ros.Username, ros.Password)
	}

	return nil
}
