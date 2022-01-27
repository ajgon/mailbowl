package smtp_test

import (
	"errors"
	"os"
	"testing"

	"github.com/ajgon/mailbowl/config"
	"github.com/ajgon/mailbowl/listener/smtp"
	"github.com/stretchr/testify/assert"
)

const (
	tlsKeyExample = `
-----BEGIN PRIVATE KEY-----
MC4CAQAwBQYDK2VwBCIEIFEAjSPdjL8yISHYSfaPozv4elFLNS6W22wtg0hmVTT9
-----END PRIVATE KEY-----
	`
	tlsCertificateExample = `
-----BEGIN CERTIFICATE-----
MIHqMIGdAhQDD5JzHsbJfBTphCK+W/f4YgxTxDAFBgMrZXAwGDEWMBQGA1UEAwwN
ZXhhbXBsZS5sb2NhbDAeFw0yMjAxMDkxMTIxMzVaFw0yMzAxMDkxMTIxMzVaMBgx
FjAUBgNVBAMMDWV4YW1wbGUubG9jYWwwKjAFBgMrZXADIQBQKVBjCeG9AkIPnb3M
JIGqrXp3fzdgWEkXVMWLMFSAyTAFBgMrZXADQQD5m7VK1sEyVz+kZXt6GoB1/rK0
cMjucM+ZnDLJX5dUjj9SmRZdqxBgsx/bRCF7f8Lieu7mykNATBLN5CxGRH4E
-----END CERTIFICATE-----
	`
)

func TestValidTLSArgs(t *testing.T) {
	t.Parallel()

	var (
		gotTLS      *smtp.TLS
		gotErr, err error
	)

	gotTLS, gotErr = smtp.NewTLS(config.SMTPTLS{ForceForStartTLS: true})

	if gotTLS != nil {
		t.Errorf("got %+v, want %+v", gotTLS, nil)
	}

	if !errors.Is(gotErr, smtp.ErrTLSNotConfigured) {
		t.Errorf("got wrong error type, want ErrTLSNotConfigured")
	}

	gotTLS, gotErr = smtp.NewTLS(
		config.SMTPTLS{Key: tlsKeyExample, Certificate: tlsCertificateExample, ForceForStartTLS: true},
	)
	assert.NoError(t, gotErr)

	if !gotTLS.ForceForStartTLS {
		t.Errorf("got %t, want %t", gotTLS.ForceForStartTLS, true)
	}

	keyFile, err := os.CreateTemp("", "keyfile")
	assert.NoError(t, err)

	defer os.Remove(keyFile.Name())

	_, err = keyFile.Write([]byte(tlsKeyExample))
	assert.NoError(t, err)

	certificateFile, err := os.CreateTemp("", "certificatefile")
	assert.NoError(t, err)

	defer os.Remove(certificateFile.Name())

	_, err = certificateFile.Write([]byte(tlsCertificateExample))
	assert.NoError(t, err)

	gotTLS, gotErr = smtp.NewTLS(
		config.SMTPTLS{KeyFile: keyFile.Name(), CertificateFile: certificateFile.Name(), ForceForStartTLS: false},
	)
	assert.NoError(t, gotErr)

	if gotTLS.ForceForStartTLS {
		t.Errorf("got %t, want %t", gotTLS.ForceForStartTLS, false)
	}
}
