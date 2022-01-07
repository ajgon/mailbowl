package smtp

import (
	"context"
	"fmt"

	"github.com/Masterminds/log-go"
)

type SMTP struct {
	Servers []*Server
}

func NewSMTP(
	hostname string, limit *Limit, uris []string, timeout *Timeout, tls *TLS, whitelist *Whitelist,
) *SMTP {
	smtp := &SMTP{Servers: make([]*Server, 0)}
	brokenURIs := false

	for _, uri := range uris {
		smtpURI, err := NewURI(uri)
		if err != nil {
			brokenURIs = true

			log.Errorw("invalid SMTP listener URI: %s", log.Fields{"uri": uri})
		} else {
			server := NewServer(hostname, limit, smtpURI, timeout, tls, whitelist)
			smtp.Servers = append(smtp.Servers, server)
		}
	}

	if brokenURIs {
		log.Fatalf("one of SMTP listener uris is invalid, refusing to start")
	}

	return smtp
}

func (s *SMTP) GetName() string {
	return "SMTP"
}

func (s *SMTP) Serve(ctx context.Context) (err error) {
	for _, server := range s.Servers {
		if err = server.Build(); err != nil {
			log.Fatalf("error starting SMTP server (%s): %s", server.URI.String(), err.Error())
		}

		go server.Start()

		log.Infow("SMTP server started", log.Fields{"server": server.URI.String()})
	}

	<-ctx.Done()

	for _, server := range s.Servers {
		log.Debugw("stopping SMTP server", log.Fields{"server": server.URI.String()})

		if err := server.Shutdown(); err != nil {
			return fmt.Errorf("error stopping SMTP server (%s): %w", server.URI.String(), err)
		}

		log.Debugw("SMTP server shut down", log.Fields{"server": server.URI.String()})
	}

	return nil
}
