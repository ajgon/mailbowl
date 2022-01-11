package smtp_test

import (
	"testing"

	"github.com/ajgon/mailbowl/listener/smtp"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticate(t *testing.T) {
	t.Parallel()

	authUser := smtp.NewAuthUser("test@example.local", "$2a$10$BoHLl7lps2ZhB.B5h3Zqau.p4VAQR7BVjmWTC7nEbDAY9Kp4LjNrW")

	assert.False(t, authUser.Authenticate("test@example.local", "wrongpassword"))
	assert.False(t, authUser.Authenticate("wrong@example.local", "test"))
	assert.True(t, authUser.Authenticate("test@example.local", "test"))
}
