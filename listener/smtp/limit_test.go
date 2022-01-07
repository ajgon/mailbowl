package smtp_test

import (
	"testing"

	"github.com/ajgon/mailbowl/listener/smtp"
	"github.com/stretchr/testify/assert"
)

func TestInvalidLimitArgs(t *testing.T) {
	t.Parallel()

	var gotLimit, wantLimit *smtp.Limit

	gotLimit = smtp.NewLimit(0, 0, 0)
	wantLimit = &smtp.Limit{Connections: 100, MessageSize: 26214400, Recipients: 100}
	assert.Equal(t, gotLimit, wantLimit)

	gotLimit = smtp.NewLimit(0, -10, -10)
	assert.Equal(t, gotLimit, wantLimit)
}

func TestValidLimitArgs(t *testing.T) {
	t.Parallel()

	var gotLimit, wantLimit *smtp.Limit

	gotLimit = smtp.NewLimit(10, 1048576, 50)
	wantLimit = &smtp.Limit{Connections: 10, MessageSize: 1048576, Recipients: 50}

	assert.Equal(t, gotLimit, wantLimit)

	gotLimit = smtp.NewLimit(-1, 100, 1000)
	wantLimit = &smtp.Limit{Connections: -1, MessageSize: 100, Recipients: 1000}

	assert.Equal(t, gotLimit, wantLimit)
}
