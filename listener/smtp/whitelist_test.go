package smtp_test

import (
	"testing"

	"github.com/ajgon/mailbowl/listener/smtp"
	"github.com/stretchr/testify/assert"
)

func TestWhitelist(t *testing.T) {
	t.Parallel()

	gotWhitelist := smtp.NewWhitelist([]string{
		"192.168.1.1",   // valid IPv4, invalid CIDR
		"10.97.5.2/12",  // valid CIDR
		"15.260.12.5/8", // invalid IPv4, valid CIDR
		"4.8.15.16/32",  // valid IPv4, CIDR for one IPv4
		"0.0.0.0/0",     // all IPv4s CIDR
		"2001:4860:4860:1234:5678:0000:4242:8888", // valid IPv6, invalid CIDR
		"::8:4:2:1/100",         // valid IPv6 CIDR
		"::12345:fede/80",       // invalid IPv6, valid CIDR
		"::4:8:15:16:23:42/128", // valid IPv6, CIDR for one IPv6
		"::/0",                  // all IPv6s CIDR
	})
	wantWhitelist := &smtp.Whitelist{CIDRs: []string{
		"10.97.5.2/12", "4.8.15.16/32", "0.0.0.0/0", "::8:4:2:1/100", "::4:8:15:16:23:42/128", "::/0",
	}}

	assert.Equal(t, gotWhitelist, wantWhitelist)
}
