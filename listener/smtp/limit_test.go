package smtp_test

import (
	"testing"

	"github.com/ajgon/mailbowl/config"
	"github.com/ajgon/mailbowl/listener/smtp"
	"github.com/stretchr/testify/assert"
)

func TestInvalidLimitArgs(t *testing.T) {
	t.Parallel()

	var gotLimit, wantLimit *smtp.Limit

	gotLimit = smtp.NewLimit(config.SMTPLimit{Connections: 0, MessageSize: 0, Recipients: 0})
	wantLimit = &smtp.Limit{Connections: 100, MessageSize: 26214400, Recipients: 100}
	assert.Equal(t, gotLimit, wantLimit)

	gotLimit = smtp.NewLimit(config.SMTPLimit{Connections: 0, MessageSize: -10, Recipients: -10})
	assert.Equal(t, gotLimit, wantLimit)
}

func TestValidLimitArgs(t *testing.T) {
	t.Parallel()

	var gotLimit, wantLimit *smtp.Limit

	gotLimit = smtp.NewLimit(config.SMTPLimit{Connections: 10, MessageSize: 1048576, Recipients: 50})
	wantLimit = &smtp.Limit{Connections: 10, MessageSize: 1048576, Recipients: 50}

	assert.Equal(t, gotLimit, wantLimit)

	gotLimit = smtp.NewLimit(config.SMTPLimit{Connections: -1, MessageSize: 100, Recipients: 1000})
	wantLimit = &smtp.Limit{Connections: -1, MessageSize: 100, Recipients: 1000}

	assert.Equal(t, gotLimit, wantLimit)
}
