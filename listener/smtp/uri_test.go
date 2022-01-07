package smtp_test

import (
	"testing"

	"github.com/ajgon/mailbowl/listener/smtp"
	"github.com/stretchr/testify/assert"
)

func TestInvalidURIArgs(t *testing.T) {
	t.Parallel()

	var (
		gotURI  *smtp.URI
		gotErr  error
		wantErr string
	)

	gotURI, gotErr = smtp.NewURI("wr%ong")
	assert.Nil(t, gotURI)

	if gotURI != nil {
		t.Errorf("got %+v, want %+v", gotURI, nil)
	}

	wantErr = "error parsing `wr%ong` uri: parse \"wr%ong\": invalid URL escape \"%on\""
	assert.Equal(t, gotErr.Error(), wantErr)

	gotURI, gotErr = smtp.NewURI("http://example.local/")
	assert.Nil(t, gotURI)

	wantErr = "invalid smtp server scheme `http`, must be one of `plain`, `tls` or `starttls`"
	assert.Equal(t, gotErr.Error(), wantErr)
}

func TestValidURIArgs(t *testing.T) {
	t.Parallel()

	var (
		gotURI, wantURI *smtp.URI
		gotErr          error
		uri             string
	)

	uri = "plain://example.local:25"
	gotURI, gotErr = smtp.NewURI(uri)
	wantURI = &smtp.URI{Scheme: "plain", Address: "example.local:25"}

	assert.Nil(t, gotErr)
	assert.Equal(t, gotURI, wantURI)
	assert.Equal(t, gotURI.String(), uri)

	uri = "tls://example.local:465"
	gotURI, gotErr = smtp.NewURI(uri)
	wantURI = &smtp.URI{Scheme: "tls", Address: "example.local:465"}

	assert.Nil(t, gotErr)
	assert.Equal(t, gotURI, wantURI)
	assert.Equal(t, gotURI.String(), uri)

	uri = "starttls://example.local:587"
	gotURI, gotErr = smtp.NewURI(uri)
	wantURI = &smtp.URI{Scheme: "starttls", Address: "example.local:587"}

	assert.Nil(t, gotErr)
	assert.Equal(t, gotURI, wantURI)
	assert.Equal(t, gotURI.String(), uri)
}
