package smtp_test

import (
	"testing"
	"time"

	"github.com/ajgon/mailbowl/config"
	"github.com/ajgon/mailbowl/listener/smtp"
	"github.com/stretchr/testify/assert"
)

func TestValidTimeoutArgs(t *testing.T) {
	t.Parallel()

	var gotTimeout, wantTimeout *smtp.Timeout

	gotTimeout = smtp.NewTimeout(config.SMTPTimeout{Read: 30 * time.Second, Write: 10 * time.Minute, Data: 2 * time.Hour})
	wantTimeout = &smtp.Timeout{Read: 30 * time.Second, Write: 10 * time.Minute, Data: 2 * time.Hour}

	assert.Equal(t, gotTimeout, wantTimeout)
}
