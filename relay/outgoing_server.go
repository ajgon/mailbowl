package relay

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
)

const (
	ConnectionPlain    = "plain"
	ConnectionStartTLS = "starttls"
	ConnectionTLS      = "tls"
)

type OutgoingServer struct {
	AuthMethod     string
	ConnectionType string
	FromEmail      string
	Host           string
	Password       string
	Port           string
	Username       string
	VerifyTLS      bool
}

func NewOutgoingServer(
	address, authMetod, connectionType, fromEmail, password, username string, verifyTLS bool,
) (*OutgoingServer, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return nil, fmt.Errorf("error configuring outgoing server: %w", err)
	}

	return &OutgoingServer{
		AuthMethod:     authMetod,
		ConnectionType: connectionType,
		FromEmail:      fromEmail,
		Host:           host,
		Password:       password,
		Port:           port,
		Username:       username,
		VerifyTLS:      verifyTLS,
	}, nil
}

func (ros *OutgoingServer) Send(from string, recipients []string, message []byte) error {
	if ros.FromEmail != "" {
		from = ros.FromEmail
	}

	if ros.ConnectionType == "" {
		switch ros.Port {
		case "465", "smtps":
			ros.ConnectionType = ConnectionTLS
		case "587", "starttls":
			ros.ConnectionType = ConnectionStartTLS
		default:
			ros.ConnectionType = ConnectionPlain
		}
	}

	if ros.ConnectionType == ConnectionPlain {
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
		InsecureSkipVerify: !ros.VerifyTLS, //nolint: gosec
		ServerName:         ros.Host,
	}

	if ros.ConnectionType == ConnectionStartTLS {
		client, err = smtp.Dial(fmt.Sprintf("%s:%s", ros.Host, ros.Port))
		if err != nil {
			return nil, fmt.Errorf("outgoing starttls error: %w", err)
		}

		err = client.StartTLS(tlsconfig)
		if err != nil {
			return nil, fmt.Errorf("outgoing starttls error: %w", err)
		}

		return client, nil
	}

	conn, err = tls.Dial("tcp", fmt.Sprintf("%s:%s", ros.Host, ros.Port), tlsconfig)
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

	err := smtp.SendMail(fmt.Sprintf("%s:%s", ros.Host, ros.Port), auth, from, recipients, message)
	if err != nil {
		return fmt.Errorf("error sending email via outgoing server: %w", err)
	}

	return nil
}

func (ros *OutgoingServer) buildAuth() smtp.Auth {
	if ros.AuthMethod == "plain" {
		return smtp.PlainAuth("", ros.Username, ros.Password, ros.Host)
	}

	return nil
}
