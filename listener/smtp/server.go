package smtp

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"

	"github.com/Masterminds/log-go"
	"github.com/ajgon/mailbowl/relay"
	"github.com/chrj/smtpd"
)

const (
	ServiceNotAvailable              = 421
	AuthenticationCredentialsInvalid = 535
	TransactionFailed                = 554
)

var errMissingTLSConfig = errors.New("server configured, but TLS config is missing")

func ErrMissingTLSConfig(tlsType string) error {
	return fmt.Errorf("%s %w", tlsType, errMissingTLSConfig)
}

type Server struct {
	Auth      *Auth
	Hostname  string
	Limit     *Limit
	Timeout   *Timeout
	TLS       *TLS
	Whitelist *Whitelist

	URI      *URI
	Relay    *relay.Relay
	SMTPD    *smtpd.Server
	Listener net.Listener
}

func NewServer(
	auth *Auth, hostname string, limit *Limit, timeout *Timeout, tls *TLS, whitelist *Whitelist, uri *URI,
	relay *relay.Relay,
) *Server {
	server := &Server{
		Auth:      auth,
		Hostname:  hostname,
		Limit:     limit,
		Timeout:   timeout,
		TLS:       tls,
		Whitelist: whitelist,

		URI:   uri,
		Relay: relay,
	}

	return server
}

func (s *Server) Start() {
	if err := s.SMTPD.Serve(s.Listener); err != nil && !errors.Is(err, smtpd.ErrServerClosed) {
		log.Fatalf("SMTP server (%s) died: %s", s.URI.String(), err.Error())
	}
}

func (s *Server) Shutdown() error {
	defer s.Listener.Close()

	if err := s.SMTPD.Shutdown(true); err != nil {
		return fmt.Errorf("%w", s.SMTPD.Shutdown(true))
	}

	return nil
}

func (s *Server) Build() (err error) {
	s.SMTPD = &smtpd.Server{
		Hostname: s.Hostname,

		ReadTimeout:  s.Timeout.Read,
		WriteTimeout: s.Timeout.Write,
		DataTimeout:  s.Timeout.Data,

		MaxConnections: s.Limit.Connections,
		MaxMessageSize: s.Limit.MessageSize,
		MaxRecipients:  s.Limit.Recipients,

		ConnectionChecker: s.connectionChecker,
		Handler:           s.handler,
	}

	if s.Auth.Enabled {
		s.SMTPD.Authenticator = s.authenticator
	}

	switch s.URI.Scheme {
	case "plain":
		s.Listener, err = net.Listen("tcp", s.URI.Address)
	case "starttls":
		if s.TLS == nil {
			return ErrMissingTLSConfig("starttls")
		}

		s.SMTPD.ForceTLS = s.TLS.ForceForStartTLS
		s.SMTPD.TLSConfig = s.TLS.Config

		s.Listener, err = net.Listen("tcp", s.URI.Address)
	case "tls":
		if s.TLS == nil {
			return ErrMissingTLSConfig("tls")
		}

		s.SMTPD.TLSConfig = s.TLS.Config

		s.Listener, err = tls.Listen("tcp", s.URI.Address, s.SMTPD.TLSConfig)
	}

	if err != nil {
		return fmt.Errorf("error building listener: %w", err)
	}

	return nil
}

func (s *Server) connectionChecker(peer smtpd.Peer) error {
	remoteIP := peer.Addr.(*net.TCPAddr).IP

	log.Debugw("new SMTP connection", log.Fields{"server": s.URI.String(), "remote_ip": remoteIP})

	testIP, _, err := net.ParseCIDR(fmt.Sprintf("%s/32", remoteIP))
	if err != nil {
		return fmt.Errorf("error processing remote IP: %w", err)
	}

	for _, cidr := range s.Whitelist.CIDRs {
		_, ipnet, err := net.ParseCIDR(cidr)

		if err == nil && ipnet.Contains(testIP) {
			return nil
		}
	}

	log.Infow("IP not included in whitelist, access denied", log.Fields{"server": s.URI.String(), "remote_ip": remoteIP})

	return smtpd.Error{Code: ServiceNotAvailable, Message: "Denied"}
}

func (s *Server) authenticator(peer smtpd.Peer, username string, password string) error {
	remoteIP := peer.Addr.(*net.TCPAddr).IP

	for _, user := range s.Auth.Users {
		if user.Authenticate(username, password) {
			return nil
		}
	}

	log.Infow(
		"authorization failed, access denied",
		log.Fields{"server": s.URI.String(), "remote_ip": remoteIP, "username": username},
	)

	return smtpd.Error{Code: AuthenticationCredentialsInvalid, Message: "Authentication credentials invalid"}
}

func (s *Server) handler(peer smtpd.Peer, envelope smtpd.Envelope) error {
	if s.Relay == nil {
		return nil
	}

	remoteIP := peer.Addr.(*net.TCPAddr).IP

	log.Infow("processing email", log.Fields{
		"server": s.URI.String(), "from": envelope.Sender, "to": envelope.Recipients, "remote_ip": remoteIP,
	})

	envelope.AddReceivedLine(peer)

	err := s.Relay.Handle(envelope.Sender, envelope.Recipients, envelope.Data)
	if err != nil {
		log.Errorf("forwarding failed", log.Fields{
			"server": s.URI.String(), "from": envelope.Sender, "to": envelope.Recipients, "remote_ip": remoteIP,
			"error": err.Error(),
		})

		return smtpd.Error{Code: TransactionFailed, Message: "forwarding failed"}
	}

	log.Infow("forwarding succeeded, mail sent", log.Fields{
		"server": s.URI.String(), "from": envelope.Sender, "to": envelope.Recipients, "remote_ip": remoteIP,
	})

	return nil
}
