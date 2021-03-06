package smtp

import (
	"crypto/tls"
	"errors"
	"fmt"

	"github.com/ajgon/mailbowl/config"
)

var ErrTLSNotConfigured = errors.New("TLS not configured")

type TLS struct {
	Config           *tls.Config
	ForceForStartTLS bool
}

func NewTLS(conf config.SMTPTLS) (*TLS, error) {
	var (
		cert tls.Certificate
		err  error
	)

	if conf.Key == "" && conf.Certificate == "" && conf.KeyFile == "" && conf.CertificateFile == "" {
		return nil, ErrTLSNotConfigured
	}

	tlsCipherSuites := []uint16{
		tls.TLS_AES_128_GCM_SHA256,
		tls.TLS_AES_256_GCM_SHA384,
		tls.TLS_CHACHA20_POLY1305_SHA256,
		tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_RSA_WITH_AES_128_GCM_SHA256, // does not provide PFS
		tls.TLS_RSA_WITH_AES_256_GCM_SHA384, // does not provide PFS
	}

	// first try load direct TLS (it takes precedence)
	cert, err = tls.X509KeyPair([]byte(conf.Certificate), []byte(conf.Key))
	if err != nil {
		// okay, let's try file then
		cert, err = tls.LoadX509KeyPair(conf.CertificateFile, conf.KeyFile)
		if err != nil {
			// still no luck? then fail
			return nil, fmt.Errorf("invalid TLS configuration: %w", err)
		}
	}

	return &TLS{
		Config: &tls.Config{
			PreferServerCipherSuites: true,
			MinVersion:               tls.VersionTLS12,
			CipherSuites:             tlsCipherSuites,
			Certificates:             []tls.Certificate{cert},
		},
		ForceForStartTLS: conf.ForceForStartTLS,
	}, nil
}
