package smtp_test

import (
	"testing"
	"time"

	"github.com/ajgon/mailbowl/listener/smtp"
	"github.com/stretchr/testify/assert"
)

func TestInvalidTimeoutArgs(t *testing.T) {
	t.Parallel()

	var gotTimeout, wantTimeout *smtp.Timeout

	gotTimeout = smtp.NewTimeout("wrong", "bad", "invalid")
	wantTimeout = &smtp.Timeout{Read: time.Minute, Write: time.Minute, Data: 5 * time.Minute}

	assert.Equal(t, gotTimeout, wantTimeout)
}

func TestValidTimeoutArgs(t *testing.T) {
	t.Parallel()

	var gotTimeout, wantTimeout *smtp.Timeout

	gotTimeout = smtp.NewTimeout("30s", "10m", "2h")
	wantTimeout = &smtp.Timeout{Read: 30 * time.Second, Write: 10 * time.Minute, Data: 2 * time.Hour}

	assert.Equal(t, gotTimeout, wantTimeout)
}
